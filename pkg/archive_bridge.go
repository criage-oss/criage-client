package pkg

import (
	commonarchive "github.com/criage-oss/criage-common/archive"
	commonconfig "github.com/criage-oss/criage-common/config"
	commontypes "github.com/criage-oss/criage-common/types"
)

// NewCommonArchiveManager создает общий архивный менеджер на базе конфигурации клиента
func NewCommonArchiveManager(cfg *Config, version string) (*commonarchive.Manager, error) {
	commonCfg := toCommonConfig(cfg)
	if commonCfg == nil {
		commonCfg = commonconfig.DefaultConfig()
	}
	return commonarchive.NewManager(commonCfg, version)
}

// toCommonConfig минимальный маппинг локального Config в общий config.Config
func toCommonConfig(cfg *Config) *commonconfig.Config {
	if cfg == nil {
		return commonconfig.DefaultConfig()
	}
	repos := make([]commontypes.Repository, 0, len(cfg.Repositories))
	for _, r := range cfg.Repositories {
		repos = append(repos, commontypes.Repository{
			Name:      r.Name,
			URL:       r.URL,
			Priority:  r.Priority,
			Enabled:   r.Enabled,
			AuthToken: r.AuthToken,
		})
	}
	return &commonconfig.Config{
		InstallPath:      cfg.GlobalPath,
		CachePath:        cfg.CachePath,
		TempPath:         cfg.TempPath,
		Timeout:          cfg.Timeout,
		MaxConnections:   10,
		UserAgent:        "CriageClient",
		Repositories:     repos,
		CompressionLevel: cfg.Compression.Level,
		PreferredFormat:  cfg.Compression.Format,
		Parallel:         cfg.Parallel > 1,
		MaxParallel:      cfg.Parallel,
		Language:         "en",
		Debug:            false,
	}
}
