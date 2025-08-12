package pkg

import (
	commontypes "github.com/criage-oss/criage-common/types"
)

type PackageManifest = commontypes.PackageManifest

type PackageHooks = commontypes.PackageHooks

type BuildManifest = commontypes.BuildManifest

type BuildTarget = commontypes.BuildTarget

type CompressionConfig = commontypes.CompressionConfig

type PackageMetadata = commontypes.PackageMetadata

type PackageInfo = commontypes.PackageInfo

type Repository = commontypes.Repository

// Config представляет конфигурацию criage
type Config struct {
	GlobalPath   string                 `yaml:"global_path" json:"global_path"`
	LocalPath    string                 `yaml:"local_path" json:"local_path"`
	CachePath    string                 `yaml:"cache_path" json:"cache_path"`
	TempPath     string                 `yaml:"temp_path" json:"temp_path"`
	Repositories []Repository           `yaml:"repositories" json:"repositories"`
	Compression  CompressionConfig      `yaml:"compression" json:"compression"`
	Parallel     int                    `yaml:"parallel" json:"parallel"`
	Timeout      int                    `yaml:"timeout" json:"timeout"`
	RetryCount   int                    `yaml:"retry_count" json:"retry_count"`
	AutoUpdate   bool                   `yaml:"auto_update" json:"auto_update"`
	VerifyHashes bool                   `yaml:"verify_hashes" json:"verify_hashes"`
	Settings     map[string]interface{} `yaml:"settings" json:"settings"`
}

type SearchResult = commontypes.SearchResult

type PackageEntry = commontypes.PackageEntry

type VersionEntry = commontypes.VersionEntry

type FileEntry = commontypes.FileEntry

type ApiResponse = commontypes.ApiResponse

type RepositoryStats = commontypes.Statistics

type ArchiveFormat = commontypes.ArchiveFormat

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
				Priority: 100,
				Enabled:  true,
			},
		},
		Compression: CompressionConfig{
			Format: string(commontypes.FormatTarZst),
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
