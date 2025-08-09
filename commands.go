package main

import (
	"fmt"
	"os"

	"criage/pkg"

	"github.com/spf13/cobra"
)

var packageManager *pkg.PackageManager

func init() {
	var err error
	packageManager, err = pkg.NewPackageManager()
	if err != nil {
		fmt.Print(pkg.T("error_init_package_manager", err))
		os.Exit(1)
	}
}

// installPackage устанавливает пакет
func installPackage(packageName string) error {
	return packageManager.InstallPackage(packageName, "", false, false, false, "", "")
}

// uninstallPackage удаляет пакет
func uninstallPackage(packageName string) error {
	return packageManager.UninstallPackage(packageName, false, false)
}

// updatePackage обновляет пакет
func updatePackage(packageName string) error {
	return packageManager.UpdatePackage(packageName)
}

// updateAllPackages обновляет все пакеты
func updateAllPackages() error {
	packages, err := packageManager.ListPackages(false, true)
	if err != nil {
		return err
	}

	for _, packageInfo := range packages {
		if err := packageManager.UpdatePackage(packageInfo.Name); err != nil {
			fmt.Print(pkg.T("failed_to_update", packageInfo.Name, err))
		}
	}
	return nil
}

// searchPackages выполняет поиск пакетов
func searchPackages(query string) error {
	results, err := packageManager.SearchPackages(query)
	if err != nil {
		return err
	}

	fmt.Print(pkg.T("packages_found", len(results)))
	for _, result := range results {
		fmt.Printf("- %s (%s): %s\n", result.Name, result.Version, result.Description)
	}
	return nil
}

// listPackages показывает список установленных пакетов
func listPackages() error {
	packages, err := packageManager.ListPackages(false, false)
	if err != nil {
		return err
	}

	fmt.Print(pkg.T("packages_installed", len(packages)))
	for _, pkg := range packages {
		fmt.Printf("- %s (%s)\n", pkg.Name, pkg.Version)
	}
	return nil
}

// showPackageInfo показывает информацию о пакете
func showPackageInfo(packageName string) error {
	info, err := packageManager.GetPackageInfo(packageName)
	if err != nil {
		return err
	}

	fmt.Printf("%s: %s\n", pkg.T("package_name"), info.Name)
	fmt.Printf("%s: %s\n", pkg.T("package_version"), info.Version)
	fmt.Printf("%s: %s\n", pkg.T("package_description"), info.Description)
	fmt.Printf("%s: %s\n", pkg.T("package_author"), info.Author)
	fmt.Printf("%s: %s\n", pkg.T("package_installed_date"), info.InstallDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("%s: %d байт\n", pkg.T("package_size"), info.Size)
	return nil
}

// createPackage создает новый пакет
func createPackage(name string) error {
	return packageManager.CreatePackage(name, "basic", "", "")
}

// publishPackage публикует пакет
func publishPackage() error {
	return packageManager.PublishPackage("", "")
}

// showArchiveMetadata показывает метаданные архива
func showArchiveMetadata(archivePath string) error {
	archiveManager, err := pkg.NewArchiveManager(pkg.DefaultConfig(), version)
	if err != nil {
		return fmt.Errorf("failed to create archive manager: %w", err)
	}
	defer archiveManager.Close()

	format := archiveManager.DetectFormat(archivePath)
	metadata, err := archiveManager.ExtractMetadataFromArchive(archivePath, format)
	if err != nil {
		return fmt.Errorf("failed to extract metadata: %w", err)
	}

	fmt.Print(pkg.T("archive_metadata_title", archivePath))
	fmt.Printf("%s: %s\n", pkg.T("compression_type"), metadata.CompressionType)
	fmt.Printf("%s: %s\n", pkg.T("created_at"), metadata.CreatedAt)
	fmt.Printf("%s: %s\n", pkg.T("created_by"), metadata.CreatedBy)

	if metadata.PackageManifest != nil {
		fmt.Printf("\n%s\n", pkg.T("package_manifest_title"))
		fmt.Printf("%s: %s\n", pkg.T("package_name"), metadata.PackageManifest.Name)
		fmt.Printf("%s: %s\n", pkg.T("package_version"), metadata.PackageManifest.Version)
		fmt.Printf("%s: %s\n", pkg.T("package_description"), metadata.PackageManifest.Description)
		fmt.Printf("%s: %s\n", pkg.T("package_author"), metadata.PackageManifest.Author)
		fmt.Printf("%s: %s\n", pkg.T("package_license"), metadata.PackageManifest.License)
		if len(metadata.PackageManifest.Dependencies) > 0 {
			fmt.Printf("%s:\n", pkg.T("package_dependencies"))
			for name, version := range metadata.PackageManifest.Dependencies {
				fmt.Printf("  - %s: %s\n", name, version)
			}
		}
	}

	if metadata.BuildManifest != nil {
		fmt.Printf("\n%s\n", pkg.T("build_manifest_title"))
		fmt.Printf("%s: %s\n", pkg.T("build_script"), metadata.BuildManifest.BuildScript)
		fmt.Printf("%s: %s\n", pkg.T("output_dir"), metadata.BuildManifest.OutputDir)
		fmt.Printf("%s: %s (уровень %d)\n", pkg.T("compression_format"),
			metadata.BuildManifest.Compression.Format,
			metadata.BuildManifest.Compression.Level)
		if len(metadata.BuildManifest.Targets) > 0 {
			fmt.Printf("%s:\n", pkg.T("target_platforms"))
			for _, target := range metadata.BuildManifest.Targets {
				fmt.Printf("  - %s/%s\n", target.OS, target.Arch)
			}
		}
	}

	return nil
}

// setConfig устанавливает значение конфигурации
func setConfig(key, value string) error {
	fmt.Print(pkg.T("config_set", key, value))
	// Здесь будет реализация установки конфигурации
	return nil
}

// getConfig получает значение конфигурации
func getConfig(key string) error {
	fmt.Print(pkg.T("config_get", key))
	// Здесь будет реализация получения конфигурации
	return nil
}

// listConfig показывает все настройки
func listConfig() error {
	fmt.Println(pkg.T("config_list"))
	// Здесь будет реализация показа всех настроек
	return nil
}

// getCurrentCommand возвращает текущую команду для доступа к флагам
func getCurrentCommand() *cobra.Command {
	// Это вспомогательная функция для получения текущей команды
	// В реальной реализации нужно передавать команду через контекст
	return &cobra.Command{}
}

// formatSize форматирует размер в человеко-читаемый формат
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
