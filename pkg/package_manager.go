package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"
)

// RateLimiter простой rate limiter для HTTP запросов
type RateLimiter struct {
	ticker   *time.Ticker
	requests chan struct{}
}

// NewRateLimiter создает новый rate limiter с заданной частотой запросов в секунду
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	if requestsPerSecond <= 0 {
		requestsPerSecond = 10 // по умолчанию 10 запросов в секунду
	}

	interval := time.Second / time.Duration(requestsPerSecond)
	ticker := time.NewTicker(interval)
	requests := make(chan struct{}, requestsPerSecond)

	// Заполняем буфер
	for i := 0; i < requestsPerSecond; i++ {
		requests <- struct{}{}
	}

	rl := &RateLimiter{
		ticker:   ticker,
		requests: requests,
	}

	// Запускаем горутину для пополнения буфера
	go func() {
		for range ticker.C {
			select {
			case requests <- struct{}{}:
			default:
				// Буфер полон, пропускаем
			}
		}
	}()

	return rl
}

// Wait ждет разрешения на выполнение запроса
func (rl *RateLimiter) Wait() {
	<-rl.requests
}

// Close останавливает rate limiter
func (rl *RateLimiter) Close() {
	rl.ticker.Stop()
	close(rl.requests)
}

// PackageManager основной менеджер пакетов
type PackageManager struct {
	configManager     *ConfigManager
	archiveManager    *ArchiveManager
	installedPackages map[string]*PackageInfo
	packagesMutex     sync.RWMutex
	httpClient        *http.Client
	rateLimiter       *RateLimiter
}

// NewPackageManager создает новый пакетный менеджер
func NewPackageManager() (*PackageManager, error) {
	configManager, err := NewConfigManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create config manager: %w", err)
	}

	// Получаем версию из переменной окружения или используем "1.0.0"
	version := os.Getenv("CRIAGE_VERSION")
	if version == "" {
		version = "1.0.0"
	}

	archiveManager, err := NewArchiveManager(configManager.GetConfig(), version)
	if err != nil {
		return nil, fmt.Errorf("failed to create archive manager: %w", err)
	}

	// Настраиваем HTTP клиент
	httpClient := &http.Client{
		Timeout: time.Duration(configManager.GetConfig().Timeout) * time.Second,
	}

	pm := &PackageManager{
		configManager:     configManager,
		archiveManager:    archiveManager,
		installedPackages: make(map[string]*PackageInfo),
		httpClient:        httpClient,
		rateLimiter:       NewRateLimiter(5), // 5 запросов в секунду
	}

	// Создаем необходимые директории
	if err := configManager.EnsureDirectories(); err != nil {
		return nil, fmt.Errorf("failed to ensure directories: %w", err)
	}

	// Загружаем информацию об установленных пакетах
	if err := pm.loadInstalledPackages(); err != nil {
		return nil, fmt.Errorf("failed to load installed packages: %w", err)
	}

	return pm, nil
}

// InstallPackage устанавливает пакет
func (pm *PackageManager) InstallPackage(packageName, version string, global, force, dev bool, arch, osName string) error {
	fmt.Printf(T("installing_package", packageName))

	// Проверяем, не установлен ли уже пакет
	if !force {
		if info, exists := pm.getInstalledPackage(packageName); exists {
			if version == "" || info.Version == version {
				fmt.Printf(T("package_already_installed", packageName, info.Version))
				return nil
			}
		}
	}

	// Определяем архитектуру и ОС
	if arch == "" {
		arch = runtime.GOARCH
	}
	if osName == "" {
		osName = runtime.GOOS
	}

	// Поиск пакета в репозиториях
	packageInfo, downloadURL, err := pm.findPackage(packageName, version, arch, osName)
	if err != nil {
		return fmt.Errorf(T("error_failed_to_find"), err)
	}

	// Скачиваем пакет
	archivePath, err := pm.downloadPackage(downloadURL, packageName, packageInfo.Version)
	if err != nil {
		return fmt.Errorf(T("error_failed_to_download"), err)
	}
	defer os.Remove(archivePath)

	// Извлекаем архив
	tempDir := pm.configManager.GetTempPath(fmt.Sprintf("install_%s_%d", packageName, time.Now().Unix()))
	defer os.RemoveAll(tempDir)

	format := pm.archiveManager.DetectFormat(archivePath)
	if err := pm.archiveManager.ExtractArchive(archivePath, tempDir, format); err != nil {
		return fmt.Errorf(T("error_failed_to_extract"), err)
	}

	// Загружаем манифест пакета
	manifest, err := pm.loadManifestFromDir(tempDir)
	if err != nil {
		return fmt.Errorf(T("error_failed_to_load"), err)
	}

	// Проверяем зависимости
	if err := pm.checkDependencies(manifest, dev); err != nil {
		return fmt.Errorf(T("error_dependency_check"), err)
	}

	// Выполняем пре-установочные хуки
	if err := pm.executeHooks(manifest.Hooks, manifest.Hooks.PreInstall, tempDir); err != nil {
		return fmt.Errorf(T("error_pre_install_hooks"), err)
	}

	// Определяем путь установки
	installPath := pm.configManager.GetInstallPath(packageName, global)

	// Удаляем старую версию, если она есть
	if force {
		if err := os.RemoveAll(installPath); err != nil {
			return fmt.Errorf(T("error_failed_to_remove"), err)
		}
	}

	// Создаем директорию установки
	if err := os.MkdirAll(installPath, 0755); err != nil {
		return fmt.Errorf(T("error_failed_to_create"), err)
	}

	// Копируем файлы
	if err := pm.copyFiles(tempDir, installPath, manifest.Files); err != nil {
		return fmt.Errorf(T("error_failed_to_copy"), err)
	}

	// Создаем информацию о пакете
	packageInfo = &PackageInfo{
		Name:         manifest.Name,
		Version:      manifest.Version,
		Description:  manifest.Description,
		Author:       manifest.Author,
		InstallDate:  time.Now(),
		InstallPath:  installPath,
		Global:       global,
		Dependencies: manifest.Dependencies,
		Size:         pm.calculateDirSize(installPath),
		Files:        manifest.Files,
		Scripts:      manifest.Scripts,
	}

	// Сохраняем информацию о пакете
	if err := pm.savePackageInfo(packageInfo); err != nil {
		return fmt.Errorf(T("error_failed_to_save"), err)
	}

	// Обновляем кеш установленных пакетов
	pm.packagesMutex.Lock()
	pm.installedPackages[packageName] = packageInfo
	pm.packagesMutex.Unlock()

	// Выполняем пост-установочные хуки
	if err := pm.executeHooks(manifest.Hooks, manifest.Hooks.PostInstall, installPath); err != nil {
		fmt.Printf(T("error_post_install_hooks", err))
	}

	fmt.Printf(T("package_installed", packageName, packageInfo.Version))
	return nil
}

// UninstallPackage удаляет пакет
func (pm *PackageManager) UninstallPackage(packageName string, global, purge bool) error {
	fmt.Printf(T("uninstalling_package", packageName))

	// Проверяем, установлен ли пакет
	packageInfo, exists := pm.getInstalledPackage(packageName)
	if !exists {
		return fmt.Errorf(T("package_not_installed", packageName))
	}

	// Загружаем манифест
	manifest, err := pm.loadManifestFromDir(packageInfo.InstallPath)
	if err != nil {
		fmt.Printf(T("warning_failed_to_load", err))
	}

	// Выполняем пре-удаление хуки
	if manifest != nil && manifest.Hooks != nil {
		if err := pm.executeHooks(manifest.Hooks, manifest.Hooks.PreRemove, packageInfo.InstallPath); err != nil {
			fmt.Printf(T("warning_pre_remove_hooks", err))
		}
	}

	// Удаляем файлы пакета
	if err := os.RemoveAll(packageInfo.InstallPath); err != nil {
		return fmt.Errorf(T("error_failed_to_remove"), err)
	}

	// Удаляем информацию о пакете
	if err := pm.removePackageInfo(packageName); err != nil {
		return fmt.Errorf(T("error_failed_to_remove"), err)
	}

	// Обновляем кеш
	pm.packagesMutex.Lock()
	delete(pm.installedPackages, packageName)
	pm.packagesMutex.Unlock()

	// Выполняем пост-удаление хуки
	if manifest != nil && manifest.Hooks != nil {
		if err := pm.executeHooks(manifest.Hooks, manifest.Hooks.PostRemove, ""); err != nil {
			fmt.Printf(T("warning_post_remove_hooks", err))
		}
	}

	fmt.Printf(T("package_uninstalled", packageName))
	return nil
}

// UpdatePackage обновляет пакет
func (pm *PackageManager) UpdatePackage(packageName string) error {
	fmt.Printf("Обновление пакета %s...\n", packageName)

	// Проверяем, установлен ли пакет
	packageInfo, exists := pm.getInstalledPackage(packageName)
	if !exists {
		return fmt.Errorf("package not installed: %s", packageName)
	}

	// Ищем последнюю версию
	latestInfo, _, err := pm.findPackage(packageName, "", runtime.GOARCH, runtime.GOOS)
	if err != nil {
		return fmt.Errorf("failed to find latest version: %w", err)
	}

	// Проверяем, нужно ли обновление
	if latestInfo.Version == packageInfo.Version {
		fmt.Printf("Пакет %s уже имеет последнюю версию (%s)\n", packageName, packageInfo.Version)
		return nil
	}

	// Выполняем обновление через переустановку
	return pm.InstallPackage(packageName, latestInfo.Version, packageInfo.Global, true, false, "", "")
}

// SearchPackages ищет пакеты в репозиториях
func (pm *PackageManager) SearchPackages(query string) ([]SearchResult, error) {
	var results []SearchResult

	repositories := pm.configManager.GetRepositories()

	for _, repo := range repositories {
		if !repo.Enabled {
			continue
		}

		repoResults, err := pm.searchInRepository(repo, query)
		if err != nil {
			fmt.Printf("Предупреждение: failed to search in repository %s: %v\n", repo.Name, err)
			continue
		}

		results = append(results, repoResults...)
	}

	// Сортируем по релевантности
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

// ListPackages возвращает список установленных пакетов
func (pm *PackageManager) ListPackages(global, outdated bool) ([]*PackageInfo, error) {
	pm.packagesMutex.RLock()
	defer pm.packagesMutex.RUnlock()

	var packages []*PackageInfo

	for _, pkg := range pm.installedPackages {
		if global && !pkg.Global {
			continue
		}
		if !global && pkg.Global {
			continue
		}

		if outdated {
			// Проверяем, есть ли более новая версия
			latestInfo, _, err := pm.findPackage(pkg.Name, "", runtime.GOARCH, runtime.GOOS)
			if err != nil || latestInfo.Version == pkg.Version {
				continue
			}
		}

		packages = append(packages, pkg)
	}

	// Сортируем по имени
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})

	return packages, nil
}

// GetPackageInfo возвращает информацию о пакете
func (pm *PackageManager) GetPackageInfo(packageName string) (*PackageInfo, error) {
	info, exists := pm.getInstalledPackage(packageName)
	if !exists {
		return nil, fmt.Errorf("package not installed: %s", packageName)
	}

	return info, nil
}

// CreatePackage создает новый пакет
func (pm *PackageManager) CreatePackage(name, template, author, description string) error {
	if err := os.MkdirAll(name, 0755); err != nil {
		return fmt.Errorf("failed to create package directory: %w", err)
	}

	manifest := &PackageManifest{
		Name:         name,
		Version:      "1.0.0",
		Description:  description,
		Author:       author,
		License:      "MIT",
		Keywords:     []string{},
		Dependencies: make(map[string]string),
		DevDeps:      make(map[string]string),
		Scripts:      make(map[string]string),
		Files:        []string{"*"},
		Exclude:      []string{".git", "node_modules", "*.log"},
		Arch:         []string{"amd64", "arm64"},
		OS:           []string{"linux", "darwin", "windows"},
		MinVersion:   "1.0.0",
		Metadata:     make(map[string]any),
	}

	// Создаем основные файлы
	if err := pm.configManager.SaveLocalConfig(name, manifest); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	// Создаем README
	readmeContent := fmt.Sprintf("# %s\n\n%s\n\n## Установка\n\n```bash\ncriage install %s\n```\n", name, description, name)
	if err := os.WriteFile(filepath.Join(name, "README.md"), []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to create README: %w", err)
	}

	// Создаем основные директории
	dirs := []string{"src", "bin", "docs"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(name, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	fmt.Printf("Пакет %s создан успешно\n", name)
	return nil
}

// BuildPackage собирает пакет с встроенными метаданными
func (pm *PackageManager) BuildPackage(outputPath, format string, compressionLevel int) error {
	fmt.Println("Сборка пакета...")

	// Загружаем локальную конфигурацию
	manifest, err := pm.configManager.LoadLocalConfig(".")
	if err != nil {
		return fmt.Errorf("failed to load local config: %w", err)
	}

	// Проверяем сборочную конфигурацию
	buildManifest, err := pm.configManager.LoadBuildConfig(".")
	if err != nil {
		// Создаем базовую конфигурацию сборки
		buildManifest = &BuildManifest{
			Name:         manifest.Name,
			Version:      manifest.Version,
			BuildScript:  "make",
			OutputDir:    "./build",
			IncludeFiles: manifest.Files,
			ExcludeFiles: manifest.Exclude,
			Compression: CompressionConfig{
				Format: format,
				Level:  compressionLevel,
			},
			Targets: []BuildTarget{
				{OS: runtime.GOOS, Arch: runtime.GOARCH},
			},
		}
	}

	// Выполняем скрипт сборки
	if buildManifest.BuildScript != "" {
		fmt.Printf("Выполнение скрипта сборки: %s\n", buildManifest.BuildScript)
		if err := pm.executeBuildScript(buildManifest); err != nil {
			return fmt.Errorf("build script failed: %w", err)
		}
	}

	// Определяем выходной файл
	if outputPath == "" {
		outputPath = fmt.Sprintf("%s-%s.criage", manifest.Name, manifest.Version)
	}

	// Создаем структуру метаданных для встраивания в архив
	metadata := &PackageMetadata{
		PackageManifest: manifest,
		BuildManifest:   buildManifest,
		CompressionType: format,
		CreatedBy:       "criage",
	}

	// Создаем архив с встроенными метаданными
	archiveFormat := ArchiveFormat(format)
	if err := pm.archiveManager.CreateArchiveWithMetadata(".", outputPath, archiveFormat, buildManifest.IncludeFiles, buildManifest.ExcludeFiles, metadata); err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}

	fmt.Printf("Пакет собран с встроенными метаданными: %s\n", outputPath)
	return nil
}

// PublishPackage публикует пакет в репозитории
func (pm *PackageManager) PublishPackage(registryURL, token string) error {
	fmt.Println("Публикация пакета...")

	// Загружаем локальную конфигурацию
	manifest, err := pm.configManager.LoadLocalConfig(".")
	if err != nil {
		return fmt.Errorf("failed to load local config: %w", err)
	}

	// Собираем пакет
	archivePath := fmt.Sprintf("%s-%s.tar.zst", manifest.Name, manifest.Version)
	if err := pm.BuildPackage(archivePath, "tar.zst", CompressionNormal); err != nil {
		return fmt.Errorf("failed to build package: %w", err)
	}
	defer os.Remove(archivePath)

	// Публикуем пакет
	if err := pm.uploadPackage(registryURL, token, archivePath, manifest); err != nil {
		return fmt.Errorf("failed to upload package: %w", err)
	}

	fmt.Printf("Пакет %s версии %s успешно опубликован\n", manifest.Name, manifest.Version)
	return nil
}

// getInstalledPackage возвращает информацию об установленном пакете
func (pm *PackageManager) getInstalledPackage(packageName string) (*PackageInfo, bool) {
	pm.packagesMutex.RLock()
	defer pm.packagesMutex.RUnlock()

	info, exists := pm.installedPackages[packageName]
	return info, exists
}

// findPackage ищет пакет в репозиториях
func (pm *PackageManager) findPackage(packageName, version, arch, osName string) (*PackageInfo, string, error) {
	repositories := pm.configManager.GetRepositories()

	// Сортируем репозитории по приоритету
	sort.Slice(repositories, func(i, j int) bool {
		return repositories[i].Priority > repositories[j].Priority
	})

	for _, repo := range repositories {
		if !repo.Enabled {
			continue
		}

		packageInfo, downloadURL, err := pm.findInRepository(repo, packageName, version, arch, osName)
		if err == nil {
			return packageInfo, downloadURL, nil
		}
	}

	return nil, "", fmt.Errorf("package not found: %s", packageName)
}

// findInRepository ищет пакет в конкретном репозитории
func (pm *PackageManager) findInRepository(repo Repository, packageName, version, arch, osName string) (*PackageInfo, string, error) {
	// Используем API v1 criage-server
	url := fmt.Sprintf("%s/api/v1/packages/%s", repo.URL, packageName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	if repo.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+repo.AuthToken)
	}

	resp, err := pm.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("package not found in repository")
	}

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, "", err
	}

	if !apiResp.Success {
		return nil, "", fmt.Errorf("API error: %s", apiResp.Error)
	}

	// Преобразуем данные в PackageEntry
	packageEntryData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, "", fmt.Errorf("unexpected API response format")
	}

	var packageEntry PackageEntry
	packageEntryBytes, err := json.Marshal(packageEntryData)
	if err != nil {
		return nil, "", err
	}

	if err := json.Unmarshal(packageEntryBytes, &packageEntry); err != nil {
		return nil, "", err
	}

	// Выбираем версию
	var selectedVersion *VersionEntry
	if version == "" {
		// Берем последнюю версию
		if len(packageEntry.Versions) > 0 {
			selectedVersion = &packageEntry.Versions[len(packageEntry.Versions)-1]
		}
	} else {
		// Ищем указанную версию
		for _, v := range packageEntry.Versions {
			if v.Version == version {
				selectedVersion = &v
				break
			}
		}
	}

	if selectedVersion == nil {
		return nil, "", fmt.Errorf("version %s not found", version)
	}

	// Ищем подходящий файл
	var selectedFile *FileEntry
	for _, file := range selectedVersion.Files {
		if file.OS == osName && file.Arch == arch {
			selectedFile = &file
			break
		}
	}

	if selectedFile == nil {
		return nil, "", fmt.Errorf("file for %s/%s not found", osName, arch)
	}

	// Создаем PackageInfo из PackageEntry
	packageInfo := &PackageInfo{
		Name:         packageEntry.Name,
		Version:      selectedVersion.Version,
		Description:  packageEntry.Description,
		Author:       packageEntry.Author,
		Dependencies: selectedVersion.Dependencies,
		Size:         selectedFile.Size,
	}

	// Строим URL для скачивания используя API v1
	downloadURL := fmt.Sprintf("%s/api/v1/download/%s/%s/%s",
		repo.URL, packageEntry.Name, selectedVersion.Version, selectedFile.Filename)

	return packageInfo, downloadURL, nil
}

// downloadPackage скачивает пакет
func (pm *PackageManager) downloadPackage(url, packageName, version string) (string, error) {
	cachePath := pm.configManager.GetCachePath(packageName, version)
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	archivePath := filepath.Join(cachePath, "package.tar.zst")

	// Проверяем, есть ли уже файл в кеше
	if _, err := os.Stat(archivePath); err == nil {
		fmt.Printf("Используется кешированная версия пакета\n")
		return archivePath, nil
	}

	fmt.Printf("Скачивание пакета из %s\n", url)

	resp, err := pm.httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download package: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download package: HTTP %d", resp.StatusCode)
	}

	outFile, err := os.Create(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to create cache file: %w", err)
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return "", fmt.Errorf("failed to save package: %w", err)
	}

	return archivePath, nil
}

// loadManifestFromDir загружает манифест из директории
func (pm *PackageManager) loadManifestFromDir(dir string) (*PackageManifest, error) {
	manifestPath := filepath.Join(dir, "criage.yaml")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("manifest not found")
	}

	return pm.configManager.LoadLocalConfig(dir)
}

// checkDependencies проверяет и устанавливает зависимости
func (pm *PackageManager) checkDependencies(manifest *PackageManifest, dev bool) error {
	dependencies := manifest.Dependencies
	if dev {
		for name, version := range manifest.DevDeps {
			dependencies[name] = version
		}
	}

	for depName, depVersion := range dependencies {
		if _, exists := pm.getInstalledPackage(depName); !exists {
			fmt.Printf("Установка зависимости: %s@%s\n", depName, depVersion)
			if err := pm.InstallPackage(depName, depVersion, false, false, false, "", ""); err != nil {
				return fmt.Errorf("failed to install dependency %s: %w", depName, err)
			}
		}
	}

	return nil
}

// executeHooks выполняет хуки жизненного цикла
func (pm *PackageManager) executeHooks(hooks *PackageHooks, commands []string, workDir string) error {
	if hooks == nil || len(commands) == 0 {
		return nil
	}

	for _, command := range commands {
		cmd := exec.Command("sh", "-c", command)
		if workDir != "" {
			cmd.Dir = workDir
		}

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("hook command failed: %s: %w", command, err)
		}
	}

	return nil
}

// Дополнительные методы для полноты реализации будут добавлены в следующих частях...

// Close освобождает ресурсы
func (pm *PackageManager) Close() error {
	return pm.archiveManager.Close()
}
