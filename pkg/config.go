package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	commontypes "github.com/criage-oss/criage-common/types"
	"gopkg.in/yaml.v3"
)

const (
	ConfigFileName    = "config.yaml"
	LocalConfigName   = "criage.yaml"
	DefaultConfigDir  = ".config/criage"
	DefaultCacheDir   = ".cache/criage"
	DefaultGlobalPath = "/usr/local/lib/criage"
	DefaultLocalPath  = "./criage_modules"
)

// ConfigManager управляет конфигурацией criage
type ConfigManager struct {
	configPath string
	config     *Config
}

// NewConfigManager создает новый менеджер конфигурации
func NewConfigManager() (*ConfigManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, DefaultConfigDir)
	configPath := filepath.Join(configDir, ConfigFileName)

	cm := &ConfigManager{
		configPath: configPath,
	}

	// Загружаем конфигурацию
	config, err := cm.loadConfig()
	if err != nil {
		// Если конфигурация не найдена, создаем новую
		config = DefaultConfig()

		// Настраиваем пути относительно домашней директории
		config.CachePath = filepath.Join(homeDir, DefaultCacheDir)
		config.TempPath = filepath.Join(os.TempDir(), "criage")

		// Создаем директории
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}

		if err := cm.saveConfig(config); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}
	}

	cm.config = config
	return cm, nil
}

// GetConfig возвращает текущую конфигурацию
func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// SetValue устанавливает значение конфигурации
func (cm *ConfigManager) SetValue(key, value string) error {
	switch key {
	case "global_path":
		cm.config.GlobalPath = value
	case "local_path":
		cm.config.LocalPath = value
	case "cache_path":
		cm.config.CachePath = value
	case "temp_path":
		cm.config.TempPath = value
	case "compression.format":
		cm.config.Compression.Format = value
	case "compression.level":
		level := 3 // По умолчанию
		if value == "fast" || value == "1" {
			level = CompressionFast
		} else if value == "normal" || value == "3" {
			level = CompressionNormal
		} else if value == "best" || value == "9" {
			level = CompressionBest
		}
		cm.config.Compression.Level = level
	case "parallel":
		var parallel int
		if _, err := fmt.Sscanf(value, "%d", &parallel); err != nil {
			return fmt.Errorf("invalid parallel value: %s", value)
		}
		if parallel < 1 || parallel > 64 {
			return fmt.Errorf("parallel must be between 1 and 64")
		}
		cm.config.Parallel = parallel
	case "timeout":
		var timeout int
		if _, err := fmt.Sscanf(value, "%d", &timeout); err != nil {
			return fmt.Errorf("invalid timeout value: %s", value)
		}
		cm.config.Timeout = timeout
	case "retry_count":
		var retryCount int
		if _, err := fmt.Sscanf(value, "%d", &retryCount); err != nil {
			return fmt.Errorf("invalid retry_count value: %s", value)
		}
		cm.config.RetryCount = retryCount
	case "auto_update":
		cm.config.AutoUpdate = strings.ToLower(value) == "true"
	case "verify_hashes":
		cm.config.VerifyHashes = strings.ToLower(value) == "true"
	default:
		// Произвольные настройки
		if cm.config.Settings == nil {
			cm.config.Settings = make(map[string]interface{})
		}
		cm.config.Settings[key] = value
	}

	return cm.saveConfig(cm.config)
}

// GetValue получает значение конфигурации
func (cm *ConfigManager) GetValue(key string) (string, error) {
	switch key {
	case "global_path":
		return cm.config.GlobalPath, nil
	case "local_path":
		return cm.config.LocalPath, nil
	case "cache_path":
		return cm.config.CachePath, nil
	case "temp_path":
		return cm.config.TempPath, nil
	case "compression.format":
		return cm.config.Compression.Format, nil
	case "compression.level":
		return fmt.Sprintf("%d", cm.config.Compression.Level), nil
	case "parallel":
		return fmt.Sprintf("%d", cm.config.Parallel), nil
	case "timeout":
		return fmt.Sprintf("%d", cm.config.Timeout), nil
	case "retry_count":
		return fmt.Sprintf("%d", cm.config.RetryCount), nil
	case "auto_update":
		return fmt.Sprintf("%t", cm.config.AutoUpdate), nil
	case "verify_hashes":
		return fmt.Sprintf("%t", cm.config.VerifyHashes), nil
	default:
		if cm.config.Settings != nil {
			if value, exists := cm.config.Settings[key]; exists {
				return fmt.Sprintf("%v", value), nil
			}
		}
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}

// ListValues возвращает все значения конфигурации
func (cm *ConfigManager) ListValues() map[string]string {
	values := map[string]string{
		"global_path":        cm.config.GlobalPath,
		"local_path":         cm.config.LocalPath,
		"cache_path":         cm.config.CachePath,
		"temp_path":          cm.config.TempPath,
		"compression.format": cm.config.Compression.Format,
		"compression.level":  fmt.Sprintf("%d", cm.config.Compression.Level),
		"parallel":           fmt.Sprintf("%d", cm.config.Parallel),
		"timeout":            fmt.Sprintf("%d", cm.config.Timeout),
		"retry_count":        fmt.Sprintf("%d", cm.config.RetryCount),
		"auto_update":        fmt.Sprintf("%t", cm.config.AutoUpdate),
		"verify_hashes":      fmt.Sprintf("%t", cm.config.VerifyHashes),
	}

	// Добавляем произвольные настройки
	if cm.config.Settings != nil {
		for key, value := range cm.config.Settings {
			values[key] = fmt.Sprintf("%v", value)
		}
	}

	return values
}

// AddRepository добавляет новый репозиторий
func (cm *ConfigManager) AddRepository(name, url, repoType string, priority int) error {
	// Проверяем, не существует ли уже репозиторий с таким именем
	for i, repo := range cm.config.Repositories {
		if repo.Name == name {
			// Обновляем существующий
			cm.config.Repositories[i] = Repository{
				Name:     name,
				URL:      url,
				Priority: priority,
				Enabled:  true,
			}
			return cm.saveConfig(cm.config)
		}
	}

	// Добавляем новый репозиторий
	cm.config.Repositories = append(cm.config.Repositories, Repository{
		Name:     name,
		URL:      url,
		Priority: priority,
		Enabled:  true,
	})

	return cm.saveConfig(cm.config)
}

// RemoveRepository удаляет репозиторий
func (cm *ConfigManager) RemoveRepository(name string) error {
	for i, repo := range cm.config.Repositories {
		if repo.Name == name {
			cm.config.Repositories = append(cm.config.Repositories[:i], cm.config.Repositories[i+1:]...)
			return cm.saveConfig(cm.config)
		}
	}

	return fmt.Errorf("repository not found: %s", name)
}

// GetRepositories возвращает список репозиториев
func (cm *ConfigManager) GetRepositories() []Repository {
	return cm.config.Repositories
}

// GetRepositoriesCommon возвращает репозитории в формате common/types.Repository
func (cm *ConfigManager) GetRepositoriesCommon() []commontypes.Repository {
	return cm.config.Repositories
}

// SetRepositoriesCommon устанавливает репозитории из формата common/types.Repository
func (cm *ConfigManager) SetRepositoriesCommon(repos []commontypes.Repository) error {
	cm.config.Repositories = repos
	return cm.saveConfig(cm.config)
}

// LoadLocalConfig загружает локальную конфигурацию проекта
func (cm *ConfigManager) LoadLocalConfig(projectPath string) (*PackageManifest, error) {
	configPath := filepath.Join(projectPath, LocalConfigName)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("local config not found")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read local config: %w", err)
	}

	var manifest PackageManifest
	if err := yaml.Unmarshal(data, &manifest); err == nil {
		return &manifest, nil
	}

	// Попробуем распарсить как общую структуру и сконвертировать
	var common commontypes.PackageManifest
	if err := yaml.Unmarshal(data, &common); err == nil {
		return &common, nil
	}
	if err := json.Unmarshal(data, &common); err == nil {
		return &common, nil
	}

	return nil, fmt.Errorf("failed to parse local config: unsupported format")
}

// SaveLocalConfig сохраняет локальную конфигурацию проекта
func (cm *ConfigManager) SaveLocalConfig(projectPath string, manifest *PackageManifest) error {
	configPath := filepath.Join(projectPath, LocalConfigName)

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// LoadBuildConfig загружает конфигурацию сборки
func (cm *ConfigManager) LoadBuildConfig(projectPath string) (*BuildManifest, error) {
	configPath := filepath.Join(projectPath, "build.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("build config not found")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read build config: %w", err)
	}

	var manifest BuildManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse build config: %w", err)
	}

	return &manifest, nil
}

// SaveBuildConfig сохраняет конфигурацию сборки
func (cm *ConfigManager) SaveBuildConfig(projectPath string, manifest *BuildManifest) error {
	configPath := filepath.Join(projectPath, "build.json")

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal build config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write build config: %w", err)
	}

	return nil
}

// EnsureDirectories создает необходимые директории
func (cm *ConfigManager) EnsureDirectories() error {
	dirs := []string{
		cm.config.CachePath,
		cm.config.TempPath,
		cm.config.LocalPath,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// loadConfig загружает конфигурацию из файла
func (cm *ConfigManager) loadConfig() (*Config, error) {
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found")
	}

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// saveConfig сохраняет конфигурацию в файл
func (cm *ConfigManager) saveConfig(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetCachePath возвращает путь к кешу для пакета
func (cm *ConfigManager) GetCachePath(packageName, version string) string {
	return filepath.Join(cm.config.CachePath, packageName, version)
}

// GetTempPath возвращает временный путь для операций
func (cm *ConfigManager) GetTempPath(suffix string) string {
	return filepath.Join(cm.config.TempPath, suffix)
}

// GetInstallPath возвращает путь установки пакета
func (cm *ConfigManager) GetInstallPath(packageName string, global bool) string {
	if global {
		return filepath.Join(cm.config.GlobalPath, packageName)
	}
	return filepath.Join(cm.config.LocalPath, packageName)
}
