package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// loadInstalledPackages загружает информацию об установленных пакетах
func (pm *PackageManager) loadInstalledPackages() error {
	// Загружаем из локальных директорий
	localPath := pm.configManager.GetConfig().LocalPath
	globalPath := pm.configManager.GetConfig().GlobalPath

	// Загружаем локальные пакеты
	if err := pm.loadPackagesFromDir(localPath, false); err != nil {
		fmt.Printf("Предупреждение: failed to load local packages: %v\n", err)
	}

	// Загружаем глобальные пакеты
	if err := pm.loadPackagesFromDir(globalPath, true); err != nil {
		fmt.Printf("Предупреждение: failed to load global packages: %v\n", err)
	}

	return nil
}

// loadPackagesFromDir загружает пакеты из директории
func (pm *PackageManager) loadPackagesFromDir(dir string, global bool) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil // Директория не существует, это нормально
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		packageDir := filepath.Join(dir, entry.Name())
		infoPath := filepath.Join(packageDir, ".criage", "package.json")

		if _, err := os.Stat(infoPath); os.IsNotExist(err) {
			continue
		}

		var info PackageInfo
		data, err := os.ReadFile(infoPath)
		if err != nil {
			continue
		}

		if err := json.Unmarshal(data, &info); err != nil {
			continue
		}

		pm.packagesMutex.Lock()
		pm.installedPackages[info.Name] = &info
		pm.packagesMutex.Unlock()
	}

	return nil
}

// savePackageInfo сохраняет информацию о пакете
func (pm *PackageManager) savePackageInfo(info *PackageInfo) error {
	infoDir := filepath.Join(info.InstallPath, ".criage")
	if err := os.MkdirAll(infoDir, 0755); err != nil {
		return err
	}

	infoPath := filepath.Join(infoDir, "package.json")
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(infoPath, data, 0644)
}

// removePackageInfo удаляет информацию о пакете
func (pm *PackageManager) removePackageInfo(packageName string) error {
	info, exists := pm.getInstalledPackage(packageName)
	if !exists {
		return nil
	}

	infoPath := filepath.Join(info.InstallPath, ".criage", "package.json")
	return os.Remove(infoPath)
}

// copyFiles копирует файлы из исходной директории в целевую
func (pm *PackageManager) copyFiles(srcDir, dstDir string, files []string) error {
	for _, pattern := range files {
		matches, err := filepath.Glob(filepath.Join(srcDir, pattern))
		if err != nil {
			return err
		}

		for _, src := range matches {
			rel, err := filepath.Rel(srcDir, src)
			if err != nil {
				return err
			}

			dst := filepath.Join(dstDir, rel)

			if err := pm.copyFile(src, dst); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile копирует отдельный файл
func (pm *PackageManager) copyFile(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return os.MkdirAll(dst, srcInfo.Mode())
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

// calculateDirSize вычисляет размер директории
func (pm *PackageManager) calculateDirSize(dir string) int64 {
	var size int64

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size
}

// searchInRepository выполняет поиск в репозитории
func (pm *PackageManager) searchInRepository(repo Repository, query string) ([]SearchResult, error) {
	// Используем API v1 criage-server
	url := fmt.Sprintf("%s/api/v1/search?q=%s", repo.URL, query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if repo.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+repo.AuthToken)
	}

	resp, err := pm.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed: HTTP %d", resp.StatusCode)
	}

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("API error: %s", apiResp.Error)
	}

	// Парсим данные поиска
	searchData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected API response format")
	}

	resultsData, ok := searchData["results"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected search results format")
	}

	var results []SearchResult
	for _, resultItem := range resultsData {
		resultBytes, err := json.Marshal(resultItem)
		if err != nil {
			continue
		}

		var result SearchResult
		if err := json.Unmarshal(resultBytes, &result); err != nil {
			continue
		}

		// Устанавливаем репозиторий для каждого результата
		result.Repository = repo.Name
		results = append(results, result)
	}

	return results, nil
}

// executeBuildScript выполняет скрипт сборки
func (pm *PackageManager) executeBuildScript(manifest *BuildManifest) error {
	cmd := exec.Command("sh", "-c", manifest.BuildScript)

	// Устанавливаем переменные окружения
	env := os.Environ()
	for key, value := range manifest.BuildEnv {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// uploadPackage загружает пакет в репозиторий
func (pm *PackageManager) uploadPackage(registryURL, token, archivePath string, manifest *PackageManifest) error {
	// Используем API v1 criage-server
	url := fmt.Sprintf("%s/api/v1/upload", registryURL)

	// Открываем файл для загрузки
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Создаем multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Добавляем файл в form
	part, err := writer.CreateFormFile("package", filepath.Base(archivePath))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	writer.Close()

	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := pm.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid authorization token")
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upload failed: HTTP %d", resp.StatusCode)
	}

	// Читаем ответ
	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return fmt.Errorf("upload failed: %s", apiResp.Error)
	}

	return nil
}

// GetConfigManager возвращает менеджер конфигурации
func (pm *PackageManager) GetConfigManager() *ConfigManager {
	return pm.configManager
}

// GetRepositoryInfo получает информацию о репозитории
func (pm *PackageManager) GetRepositoryInfo(repositoryURL string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/", repositoryURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := pm.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %d", resp.StatusCode)
	}

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("operation failed: %s", apiResp.Error)
	}

	infoData, ok := apiResp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return infoData, nil
}

// GetRepositoryStats получает статистику репозитория
func (pm *PackageManager) GetRepositoryStats(repositoryURL string) (*RepositoryStats, error) {
	url := fmt.Sprintf("%s/api/v1/stats", repositoryURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := pm.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server error: %d", resp.StatusCode)
	}

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("operation failed: %s", apiResp.Error)
	}

	statsBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stats: %w", err)
	}

	var stats RepositoryStats
	if err := json.Unmarshal(statsBytes, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stats: %w", err)
	}

	return &stats, nil
}

// RefreshRepositoryIndex обновляет индекс пакетов в репозитории
func (pm *PackageManager) RefreshRepositoryIndex(repositoryURL, authToken string) error {
	url := fmt.Sprintf("%s/api/v1/refresh", repositoryURL)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if authToken != "" {
		req.Header.Set("Authorization", "Bearer "+authToken)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := pm.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("invalid authorization token")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server error: %d", resp.StatusCode)
	}

	var apiResp ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return fmt.Errorf("operation failed: %s", apiResp.Error)
	}

	return nil
}
