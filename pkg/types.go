package pkg

import (
	"time"
)

// PackageManifest представляет манифест пакета
type PackageManifest struct {
	Name         string            `yaml:"name" json:"name"`
	Version      string            `yaml:"version" json:"version"`
	Description  string            `yaml:"description" json:"description"`
	Author       string            `yaml:"author" json:"author"`
	License      string            `yaml:"license" json:"license"`
	Homepage     string            `yaml:"homepage" json:"homepage"`
	Repository   string            `yaml:"repository" json:"repository"`
	Keywords     []string          `yaml:"keywords" json:"keywords"`
	Dependencies map[string]string `yaml:"dependencies" json:"dependencies"`
	DevDeps      map[string]string `yaml:"dev_dependencies" json:"dev_dependencies"`
	Scripts      map[string]string `yaml:"scripts" json:"scripts"`
	Files        []string          `yaml:"files" json:"files"`
	Exclude      []string          `yaml:"exclude" json:"exclude"`
	Arch         []string          `yaml:"arch" json:"arch"`
	OS           []string          `yaml:"os" json:"os"`
	MinVersion   string            `yaml:"min_version" json:"min_version"`
	Hooks        *PackageHooks     `yaml:"hooks" json:"hooks"`
	Metadata     map[string]any    `yaml:"metadata" json:"metadata"`
}

// PackageHooks представляет хуки жизненного цикла пакета
type PackageHooks struct {
	PreInstall  []string `yaml:"pre_install" json:"pre_install"`
	PostInstall []string `yaml:"post_install" json:"post_install"`
	PreRemove   []string `yaml:"pre_remove" json:"pre_remove"`
	PostRemove  []string `yaml:"post_remove" json:"post_remove"`
	PreUpdate   []string `yaml:"pre_update" json:"pre_update"`
	PostUpdate  []string `yaml:"post_update" json:"post_update"`
}

// BuildManifest представляет манифест сборки
type BuildManifest struct {
	Name          string            `yaml:"name" json:"name"`
	Version       string            `yaml:"version" json:"version"`
	BuildScript   string            `yaml:"build_script" json:"build_script"`
	BuildEnv      map[string]string `yaml:"build_env" json:"build_env"`
	OutputDir     string            `yaml:"output_dir" json:"output_dir"`
	IncludeFiles  []string          `yaml:"include_files" json:"include_files"`
	ExcludeFiles  []string          `yaml:"exclude_files" json:"exclude_files"`
	Compression   CompressionConfig `yaml:"compression" json:"compression"`
	Targets       []BuildTarget     `yaml:"targets" json:"targets"`
	Dependencies  []string          `yaml:"dependencies" json:"dependencies"`
	TestCommand   string            `yaml:"test_command" json:"test_command"`
	InstallHooks  *PackageHooks     `yaml:"install_hooks" json:"install_hooks"`
}

// BuildTarget представляет цель сборки
type BuildTarget struct {
	OS   string `yaml:"os" json:"os"`
	Arch string `yaml:"arch" json:"arch"`
}

// CompressionConfig настройки сжатия
type CompressionConfig struct {
	Format string `yaml:"format" json:"format"`
	Level  int    `yaml:"level" json:"level"`
}

// PackageInfo информация об установленном пакете
type PackageInfo struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	InstallDate  time.Time         `json:"install_date"`
	InstallPath  string            `json:"install_path"`
	Global       bool              `json:"global"`
	Dependencies map[string]string `json:"dependencies"`
	Size         int64             `json:"size"`
	Files        []string          `json:"files"`
	Scripts      map[string]string `json:"scripts"`
}

// Repository представляет репозиторий пакетов
type Repository struct {
	Name        string `yaml:"name" json:"name"`
	URL         string `yaml:"url" json:"url"`
	Type        string `yaml:"type" json:"type"`
	Priority    int    `yaml:"priority" json:"priority"`
	Enabled     bool   `yaml:"enabled" json:"enabled"`
	AuthToken   string `yaml:"auth_token" json:"auth_token"`
	Fingerprint string `yaml:"fingerprint" json:"fingerprint"`
}

// Config представляет конфигурацию criage
type Config struct {
	GlobalPath    string                 `yaml:"global_path" json:"global_path"`
	LocalPath     string                 `yaml:"local_path" json:"local_path"`
	CachePath     string                 `yaml:"cache_path" json:"cache_path"`
	TempPath      string                 `yaml:"temp_path" json:"temp_path"`
	Repositories  []Repository           `yaml:"repositories" json:"repositories"`
	Compression   CompressionConfig      `yaml:"compression" json:"compression"`
	Parallel      int                    `yaml:"parallel" json:"parallel"`
	Timeout       int                    `yaml:"timeout" json:"timeout"`
	RetryCount    int                    `yaml:"retry_count" json:"retry_count"`
	AutoUpdate    bool                   `yaml:"auto_update" json:"auto_update"`
	VerifyHashes  bool                   `yaml:"verify_hashes" json:"verify_hashes"`
	Settings      map[string]interface{} `yaml:"settings" json:"settings"`
}

// SearchResult результат поиска пакета
type SearchResult struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Repository  string    `json:"repository"`
	Downloads   int64     `json:"downloads"`
	Updated     time.Time `json:"updated"`
	Score       float64   `json:"score"`
}

// PackageEntry представляет информацию о пакете в репозитории (соответствует criage-server)
type PackageEntry struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Author        string         `json:"author"`
	License       string         `json:"license"`
	Homepage      string         `json:"homepage"`
	Repository    string         `json:"repository"`
	Keywords      []string       `json:"keywords"`
	Versions      []VersionEntry `json:"versions"`
	LatestVersion string         `json:"latest_version"`
	Downloads     int64          `json:"downloads"`
	Updated       time.Time      `json:"updated"`
}

// VersionEntry представляет конкретную версию пакета (соответствует criage-server)
type VersionEntry struct {
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Dependencies map[string]string `json:"dependencies"`
	DevDeps      map[string]string `json:"dev_dependencies"`
	Files        []FileEntry       `json:"files"`
	Size         int64             `json:"size"`
	Checksum     string            `json:"checksum"`
	Uploaded     time.Time         `json:"uploaded"`
	Downloads    int64             `json:"downloads"`
}

// FileEntry представляет файл пакета для разных платформ (соответствует criage-server)
type FileEntry struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Format   string `json:"format"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Checksum string `json:"checksum"`
}

// ApiResponse стандартный ответ API criage-server
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// RepositoryStats статистика репозитория
type RepositoryStats struct {
	TotalDownloads    int64          `json:"total_downloads"`
	PackagesByLicense map[string]int `json:"packages_by_license"`
	PackagesByAuthor  map[string]int `json:"packages_by_author"`
	PopularPackages   []string       `json:"popular_packages"`
	LastUpdated       time.Time      `json:"last_updated"`
	TotalPackages     int            `json:"total_packages"`
}

// ArchiveFormat поддерживаемые форматы архивов
type ArchiveFormat string

const (
	FormatTarZst ArchiveFormat = "tar.zst"
	FormatTarLz4 ArchiveFormat = "tar.lz4"
	FormatTarXz  ArchiveFormat = "tar.xz"
	FormatTarGz  ArchiveFormat = "tar.gz"
	FormatZip    ArchiveFormat = "zip"
)

// CompressionLevel уровни сжатия
const (
	CompressionFast   = 1
	CompressionNormal = 3
	CompressionBest   = 9
)

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		GlobalPath: "/usr/local/lib/criage",
		LocalPath:  "./criage_modules",
		CachePath:  "~/.cache/criage",
		TempPath:   "/tmp/criage",
		Repositories: []Repository{
			{
				Name:     "default",
				URL:      "https://packages.criage.io",
				Type:     "http",
				Priority: 100,
				Enabled:  true,
			},
		},
		Compression: CompressionConfig{
			Format: string(FormatTarZst),
			Level:  CompressionNormal,
		},
		Parallel:     4,
		Timeout:      60,
		RetryCount:   3,
		AutoUpdate:   false,
		VerifyHashes: true,
		Settings:     make(map[string]interface{}),
	}
}