//go:build embed
// +build embed

package pkg

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"regexp"
)

// Встроенные файлы переводов - активируются только при сборке с тегом embed
// Для активации используйте: go build -tags embed
//
//go:embed locale/*.json
var embeddedLocaleFS embed.FS

// NewEmbeddedLocalization создает новый экземпляр локализации со встроенными переводами
func NewEmbeddedLocalization() *Localization {
	l := &Localization{
		translations: make(map[string]map[string]string),
		useEmbedded:  true,
	}

	// Сканируем встроенные языки
	l.scanEmbeddedLanguages()

	// Определяем язык системы
	l.currentLanguage = l.detectSystemLanguage()

	// Инициализируем встроенные переводы
	l.initializeEmbeddedTranslations()

	return l
}

// scanEmbeddedLanguages сканирует встроенные файлы переводов
func (l *Localization) scanEmbeddedLanguages() {
	l.supportedLanguages = []string{}

	// Регулярное выражение для поиска файлов переводов
	translationFilePattern := regexp.MustCompile(`^locale/translations_([a-z]{2}(?:-[A-Z]{2})?)\.json$`)

	languageSet := make(map[string]bool)

	// Читаем встроенные файлы
	fs.WalkDir(embeddedLocaleFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Игнорируем ошибки
		}

		if d.IsDir() {
			return nil
		}

		matches := translationFilePattern.FindStringSubmatch(path)
		if len(matches) == 2 {
			languageCode := matches[1]
			if !languageSet[languageCode] {
				languageSet[languageCode] = true
				l.supportedLanguages = append(l.supportedLanguages, languageCode)
			}
		}

		return nil
	})

	// Если не найдено встроенных файлов, добавляем дефолтный язык
	if len(l.supportedLanguages) == 0 {
		l.supportedLanguages = []string{DefaultLanguage}
	}
}

// initializeEmbeddedTranslations инициализирует встроенные переводы
func (l *Localization) initializeEmbeddedTranslations() {
	for _, language := range l.supportedLanguages {
		l.translations[language] = make(map[string]string)

		// Пробуем загрузить переводы из встроенного файла
		filename := fmt.Sprintf("locale/translations_%s.json", language)

		if err := l.loadEmbeddedTranslations(language, filename); err != nil {
			// Если встроенный файл не найден, используем дефолтные переводы
			l.translations[language] = l.getDefaultTranslations(language)
		}
	}
}

// loadEmbeddedTranslations загружает переводы из встроенного файла
func (l *Localization) loadEmbeddedTranslations(language, filename string) error {
	data, err := embeddedLocaleFS.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("embedded file not found: %s", filename)
	}

	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return fmt.Errorf("failed to parse embedded translations for %s: %v", language, err)
	}

	l.translations[language] = translations
	return nil
}

// GetEmbeddedLanguages возвращает список встроенных языков
func GetEmbeddedLanguages() []string {
	// Сканируем встроенные языки без создания полной локализации
	translationFilePattern := regexp.MustCompile(`^locale/translations_([a-z]{2}(?:-[A-Z]{2})?)\.json$`)
	languageSet := make(map[string]bool)
	languages := []string{}

	fs.WalkDir(embeddedLocaleFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		matches := translationFilePattern.FindStringSubmatch(path)
		if len(matches) == 2 {
			languageCode := matches[1]
			if !languageSet[languageCode] {
				languageSet[languageCode] = true
				languages = append(languages, languageCode)
			}
		}

		return nil
	})

	return languages
}
