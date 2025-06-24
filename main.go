package main

import (
	"fmt"
	"os"

	"criage/pkg"

	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
)

func main() {
	// Инициализируем локализацию (автоматически выбирает embedded/внешние файлы)
	l := pkg.GetLocalization()

	// Показываем режим локализации в debug режиме
	if os.Getenv("CRIAGE_DEBUG") == "1" {
		embeddedLangs := pkg.GetEmbeddedLanguages()
		if len(embeddedLangs) > 0 {
			fmt.Printf("🚀 Используются встроенные переводы: %v\n", embeddedLangs)
		} else {
			fmt.Printf("📁 Используются внешние файлы переводов: %v\n", l.GetSupportedLanguages())
		}
	}

	var rootCmd = &cobra.Command{
		Use:     "criage",
		Short:   l.Get("app_description"),
		Long:    l.Get("app_long_description"),
		Version: version,
	}

	// Команды управления пакетами
	rootCmd.AddCommand(
		newInstallCmd(),
		newUninstallCmd(),
		newUpdateCmd(),
		newSearchCmd(),
		newListCmd(),
		newInfoCmd(),
		newCreateCmd(),
		newBuildCmd(),
		newPublishCmd(),
		newConfigCmd(),
		newMetadataCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Команда установки пакетов
func newInstallCmd() *cobra.Command {
	l := pkg.GetLocalization()

	cmd := &cobra.Command{
		Use:   "install [package]",
		Short: l.Get("cmd_install"),
		Long:  l.Get("cmd_install_long"),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return installPackage(args[0])
		},
	}

	cmd.Flags().BoolP("global", "g", false, l.Get("flag_global"))
	cmd.Flags().StringP("version", "v", "", l.Get("flag_version"))
	cmd.Flags().BoolP("force", "f", false, l.Get("flag_force"))
	cmd.Flags().BoolP("dev", "d", false, l.Get("flag_dev"))
	cmd.Flags().StringP("arch", "a", "", l.Get("flag_arch"))
	cmd.Flags().StringP("os", "o", "", l.Get("flag_os"))

	return cmd
}

// Команда удаления пакетов
func newUninstallCmd() *cobra.Command {
	l := pkg.GetLocalization()

	cmd := &cobra.Command{
		Use:   "uninstall [package]",
		Short: l.Get("cmd_uninstall"),
		Long:  l.Get("cmd_uninstall_long"),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return uninstallPackage(args[0])
		},
	}

	cmd.Flags().BoolP("global", "g", false, l.Get("flag_global"))
	cmd.Flags().BoolP("purge", "p", false, l.Get("flag_purge"))

	return cmd
}

// Команда обновления пакетов
func newUpdateCmd() *cobra.Command {
	l := pkg.GetLocalization()

	cmd := &cobra.Command{
		Use:   "update [package]",
		Short: l.Get("cmd_update"),
		Long:  l.Get("cmd_update_long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return updateAllPackages()
			}
			return updatePackage(args[0])
		},
	}

	cmd.Flags().BoolP("global", "g", false, l.Get("flag_global"))
	cmd.Flags().BoolP("all", "a", false, l.Get("flag_all"))

	return cmd
}

// Команда поиска пакетов
func newSearchCmd() *cobra.Command {
	l := pkg.GetLocalization()

	return &cobra.Command{
		Use:   "search [query]",
		Short: l.Get("cmd_search"),
		Long:  l.Get("cmd_search_long"),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return searchPackages(args[0])
		},
	}
}

// Команда списка пакетов
func newListCmd() *cobra.Command {
	l := pkg.GetLocalization()

	cmd := &cobra.Command{
		Use:   "list",
		Short: l.Get("cmd_list"),
		Long:  l.Get("cmd_list_long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listPackages()
		},
	}

	cmd.Flags().BoolP("global", "g", false, l.Get("flag_global"))
	cmd.Flags().BoolP("outdated", "o", false, l.Get("flag_outdated"))

	return cmd
}

// Команда информации о пакете
func newInfoCmd() *cobra.Command {
	l := pkg.GetLocalization()

	return &cobra.Command{
		Use:   "info [package]",
		Short: l.Get("cmd_info"),
		Long:  l.Get("cmd_info_long"),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showPackageInfo(args[0])
		},
	}
}

// Команда создания пакета
func newCreateCmd() *cobra.Command {
	l := pkg.GetLocalization()

	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: l.Get("cmd_create"),
		Long:  l.Get("cmd_create_long"),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return createPackage(args[0])
		},
	}

	cmd.Flags().StringP("template", "t", "basic", l.Get("flag_template"))
	cmd.Flags().StringP("author", "a", "", l.Get("flag_author"))
	cmd.Flags().StringP("description", "d", "", l.Get("flag_description"))

	return cmd
}

// Команда сборки пакета
func newBuildCmd() *cobra.Command {
	l := pkg.GetLocalization()

	cmd := &cobra.Command{
		Use:   "build",
		Short: l.Get("cmd_build"),
		Long:  l.Get("cmd_build_long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")
			format, _ := cmd.Flags().GetString("format")
			compression, _ := cmd.Flags().GetInt("compression")

			return packageManager.BuildPackage(output, format, compression)
		},
	}

	cmd.Flags().StringP("output", "o", "", l.Get("flag_output"))
	cmd.Flags().StringP("format", "f", "tar.zst", l.Get("flag_format"))
	cmd.Flags().IntP("compression", "c", 3, l.Get("flag_compression"))

	return cmd
}

// Команда публикации пакета
func newPublishCmd() *cobra.Command {
	l := pkg.GetLocalization()

	cmd := &cobra.Command{
		Use:   "publish",
		Short: l.Get("cmd_publish"),
		Long:  l.Get("cmd_publish_long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return publishPackage()
		},
	}

	cmd.Flags().StringP("registry", "r", "", l.Get("flag_registry"))
	cmd.Flags().StringP("token", "t", "", l.Get("flag_token"))

	return cmd
}

// Команда конфигурации
func newConfigCmd() *cobra.Command {
	l := pkg.GetLocalization()

	cmd := &cobra.Command{
		Use:   "config",
		Short: l.Get("cmd_config"),
		Long:  l.Get("cmd_config_long"),
	}

	// Подкоманда set
	setCmd := &cobra.Command{
		Use:   "set [key] [value]",
		Short: l.Get("config_set"),
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setConfig(args[0], args[1])
		},
	}

	// Подкоманда get
	getCmd := &cobra.Command{
		Use:   "get [key]",
		Short: l.Get("config_get"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getConfig(args[0])
		},
	}

	// Подкоманда list
	listCmd := &cobra.Command{
		Use:   "list",
		Short: l.Get("config_list"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return listConfig()
		},
	}

	cmd.AddCommand(setCmd, getCmd, listCmd)
	return cmd
}

// Команда метаданных
func newMetadataCmd() *cobra.Command {
	l := pkg.GetLocalization()

	return &cobra.Command{
		Use:   "metadata [archive]",
		Short: l.Get("cmd_metadata"),
		Long:  l.Get("cmd_metadata_long"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showArchiveMetadata(args[0])
		},
	}
}
