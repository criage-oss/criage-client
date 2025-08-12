package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"criage/pkg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run test_metadata.go <путь к архиву>")
		os.Exit(1)
	}

	archivePath := os.Args[1]

	fmt.Printf("Проверка метаданных в архиве: %s\n", archivePath)

	// Создаем архивный менеджер
	config := pkg.DefaultConfig()
	archiveManager, err := pkg.NewCommonArchiveManager(config, "1.0.0")
	if err != nil {
		log.Fatalf("Ошибка создания архивного менеджера: %v", err)
	}
	defer archiveManager.Close()

	// Определяем формат архива
	format := archiveManager.DetectFormat(archivePath)
	fmt.Printf("Обнаружен формат: %s\n", format)

	// Извлекаем метаданные
	metadata, err := archiveManager.ExtractMetadataFromArchive(archivePath, format)
	if err != nil {
		log.Fatalf("Ошибка извлечения метаданных: %v", err)
	}

	// Выводим метаданные
	fmt.Printf("\n=== Метаданные архива %s ===\n", archivePath)
	fmt.Printf("Тип сжатия: %s\n", metadata.CompressionType)
	fmt.Printf("Создан: %s\n", metadata.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Создано с помощью: %s\n", metadata.CreatedBy)

	if metadata.PackageManifest != nil {
		fmt.Printf("\n=== Манифест пакета ===\n")
		fmt.Printf("Название: %s\n", metadata.PackageManifest.Name)
		fmt.Printf("Версия: %s\n", metadata.PackageManifest.Version)
		fmt.Printf("Описание: %s\n", metadata.PackageManifest.Description)
		fmt.Printf("Автор: %s\n", metadata.PackageManifest.Author)
		fmt.Printf("Лицензия: %s\n", metadata.PackageManifest.License)
		if len(metadata.PackageManifest.Dependencies) > 0 {
			fmt.Printf("Зависимости:\n")
			for name, version := range metadata.PackageManifest.Dependencies {
				fmt.Printf("  - %s: %s\n", name, version)
			}
		}
	}

	if metadata.BuildManifest != nil {
		fmt.Printf("\n=== Манифест сборки ===\n")
		fmt.Printf("Скрипт сборки: %s\n", metadata.BuildManifest.BuildScript)
		fmt.Printf("Выходная директория: %s\n", metadata.BuildManifest.OutputDir)
		fmt.Printf("Формат сжатия: %s (уровень %d)\n",
			metadata.BuildManifest.Compression.Format,
			metadata.BuildManifest.Compression.Level)
		if len(metadata.BuildManifest.Targets) > 0 {
			fmt.Printf("Целевые платформы:\n")
			for _, target := range metadata.BuildManifest.Targets {
				fmt.Printf("  - %s/%s\n", target.OS, target.Arch)
			}
		}
	}

	fmt.Println("\n✅ Метаданные успешно извлечены!")
}
