package main

import (
	"fmt"
	"os"

	"criage/pkg"
)

func main() {
	fmt.Println("=== Тест системы локализации Criage ===")

	// Получаем экземпляр локализации
	l := pkg.GetLocalization()

	fmt.Printf("Текущий язык: %s\n", l.GetLanguage())
	fmt.Printf("Поддерживаемые языки: %v\n", l.GetSupportedLanguages())

	fmt.Println("\n=== Примеры переводов ===")

	// Основные команды
	fmt.Printf("Описание приложения: %s\n", pkg.T("app_description"))
	fmt.Printf("Команда установки: %s\n", pkg.T("cmd_install"))
	fmt.Printf("Команда поиска: %s\n", pkg.T("cmd_search"))

	// Сообщения с параметрами
	fmt.Printf("Установка пакета: %s\n", pkg.T("installing_package", "example-package"))
	fmt.Printf("Найдено пакетов: %s\n", pkg.T("packages_found", 5))

	// Информация о пакете
	fmt.Printf("Название пакета: %s\n", pkg.T("package_name"))
	fmt.Printf("Версия пакета: %s\n", pkg.T("package_version"))

	fmt.Println("\n=== Переключение языка ===")

	// Переключаемся на английский
	if err := pkg.SetLanguage("en"); err != nil {
		fmt.Printf("Ошибка переключения языка: %v\n", err)
	} else {
		fmt.Printf("Язык переключен на: %s\n", pkg.GetLanguage())

		// Те же сообщения на английском
		fmt.Printf("App description: %s\n", pkg.T("app_description"))
		fmt.Printf("Install command: %s\n", pkg.T("cmd_install"))
		fmt.Printf("Search command: %s\n", pkg.T("cmd_search"))
		fmt.Printf("Installing package: %s\n", pkg.T("installing_package", "example-package"))
		fmt.Printf("Packages found: %s\n", pkg.T("packages_found", 5))
	}

	// Переключаемся обратно на русский
	if err := pkg.SetLanguage("ru"); err != nil {
		fmt.Printf("Ошибка переключения языка: %v\n", err)
	} else {
		fmt.Printf("\nЯзык переключен обратно на: %s\n", pkg.GetLanguage())
	}

	fmt.Println("\n=== Автоопределение языка ===")
	fmt.Printf("LANG: %s\n", os.Getenv("LANG"))
	fmt.Printf("LC_ALL: %s\n", os.Getenv("LC_ALL"))
	fmt.Printf("LANGUAGE: %s\n", os.Getenv("LANGUAGE"))

	fmt.Println("\n=== Сохранение переводов ===")

	// Демонстрация сохранения переводов в файл
	if err := l.SaveTranslationsToFile("ru", "./translations_ru.json"); err != nil {
		fmt.Printf("Ошибка сохранения русских переводов: %v\n", err)
	} else {
		fmt.Println("Русские переводы сохранены в translations_ru.json")
	}

	if err := l.SaveTranslationsToFile("en", "./translations_en.json"); err != nil {
		fmt.Printf("Ошибка сохранения английских переводов: %v\n", err)
	} else {
		fmt.Println("Английские переводы сохранены в translations_en.json")
	}

	fmt.Println("\n=== Тест завершен ===")
}
