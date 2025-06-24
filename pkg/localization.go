package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// Localization управляет локализацией приложения
type Localization struct {
	currentLanguage    string
	supportedLanguages []string
	translations       map[string]map[string]string
	translationsDir    string
	mutex              sync.RWMutex
}

// Глобальный экземпляр локализации
var globalLocalization *Localization
var localizationOnce sync.Once

// Дефолтный язык (fallback)
const DefaultLanguage = "en"

// GetLocalization возвращает глобальный экземпляр локализации
func GetLocalization() *Localization {
	localizationOnce.Do(func() {
		globalLocalization = NewLocalization()
	})
	return globalLocalization
}

// NewLocalization создает новый экземпляр локализации
func NewLocalization() *Localization {
	return NewLocalizationWithDir("locale")
}

// NewLocalizationWithDir создает новый экземпляр локализации с указанной директорией
func NewLocalizationWithDir(translationsDir string) *Localization {
	l := &Localization{
		translations:    make(map[string]map[string]string),
		translationsDir: translationsDir,
	}

	// Сканируем доступные языки
	l.scanAvailableLanguages()

	// Определяем язык системы
	l.currentLanguage = l.detectSystemLanguage()

	// Инициализируем переводы
	l.initializeTranslations()

	return l
}

// scanAvailableLanguages сканирует директорию в поисках файлов переводов
func (l *Localization) scanAvailableLanguages() {
	l.supportedLanguages = []string{}

	// Регулярное выражение для поиска файлов переводов: translations_<код_языка>.json
	translationFilePattern := regexp.MustCompile(`^translations_([a-z]{2}(?:-[A-Z]{2})?)\.json$`)

	// Читаем файлы в директории
	files, err := os.ReadDir(l.translationsDir)
	if err != nil {
		// Если не можем прочитать директорию, используем дефолтный язык
		l.supportedLanguages = []string{DefaultLanguage}
		return
	}

	languageSet := make(map[string]bool)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matches := translationFilePattern.FindStringSubmatch(file.Name())
		if len(matches) == 2 {
			languageCode := matches[1]
			if !languageSet[languageCode] {
				languageSet[languageCode] = true
				l.supportedLanguages = append(l.supportedLanguages, languageCode)
			}
		}
	}

	// Если не найдено файлов переводов, добавляем дефолтный язык
	if len(l.supportedLanguages) == 0 {
		l.supportedLanguages = []string{DefaultLanguage}
	}
}

// detectSystemLanguage определяет язык системы на основе доступных языков
func (l *Localization) detectSystemLanguage() string {
	// Проверяем переменные окружения
	for _, env := range []string{"LANG", "LC_ALL", "LC_MESSAGES", "LANGUAGE"} {
		if value := os.Getenv(env); value != "" {
			// Извлекаем код языка из переменной (например, ru_RU.UTF-8 -> ru)
			langCode := strings.ToLower(strings.Split(value, "_")[0])

			// Проверяем, поддерживается ли этот язык
			for _, supportedLang := range l.supportedLanguages {
				if strings.HasPrefix(supportedLang, langCode) {
					return supportedLang
				}
			}
		}
	}

	// В Windows используем английский по умолчанию, если он доступен
	if runtime.GOOS == "windows" {
		for _, supportedLang := range l.supportedLanguages {
			if supportedLang == DefaultLanguage {
				return DefaultLanguage
			}
		}
	}

	// Возвращаем первый доступный язык или дефолтный
	if len(l.supportedLanguages) > 0 {
		return l.supportedLanguages[0]
	}

	return DefaultLanguage
}

// initializeTranslations инициализирует переводы
func (l *Localization) initializeTranslations() {
	for _, language := range l.supportedLanguages {
		l.translations[language] = make(map[string]string)

		// Пробуем загрузить переводы из файла
		filename := fmt.Sprintf("translations_%s.json", language)
		filePath := filepath.Join(l.translationsDir, filename)

		if err := l.LoadTranslationsFromFile(language, filePath); err != nil {
			// Если файл не найден или не читается, используем минимальный набор переводов
			l.translations[language] = l.getDefaultTranslations(language)
		}
	}
}

// getDefaultTranslations возвращает минимальный набор переводов (fallback)
func (l *Localization) getDefaultTranslations(language string) map[string]string {
	// Базовые переводы для разных языков
	translations := map[string]map[string]string{
		"ru": {
			"app_description":   "Высокопроизводительный пакетный менеджер",
			"cmd_install":       "Установить пакет",
			"cmd_uninstall":     "Удалить пакет",
			"cmd_search":        "Найти пакеты",
			"cmd_list":          "Показать установленные пакеты",
			"no_packages_found": "Пакеты не найдены",
		},
		"en": {
			"app_description":   "High-performance package manager",
			"cmd_install":       "Install package",
			"cmd_uninstall":     "Uninstall package",
			"cmd_search":        "Search packages",
			"cmd_list":          "List installed packages",
			"no_packages_found": "No packages found",
		},
		"de": {
			"app_description":   "Hochleistungs-Paketmanager",
			"cmd_install":       "Paket installieren",
			"cmd_uninstall":     "Paket deinstallieren",
			"cmd_search":        "Pakete suchen",
			"cmd_list":          "Installierte Pakete anzeigen",
			"no_packages_found": "Keine Pakete gefunden",
		},
		"fr": {
			"app_description":   "Gestionnaire de paquets haute performance",
			"cmd_install":       "Installer le paquet",
			"cmd_uninstall":     "Désinstaller le paquet",
			"cmd_search":        "Rechercher des paquets",
			"cmd_list":          "Afficher les paquets installés",
			"no_packages_found": "Aucun paquet trouvé",
		},
	}

	if trans, exists := translations[language]; exists {
		return trans
	}

	// Fallback на английский
	return translations["en"]
}

// SetLanguage устанавливает текущий язык
func (l *Localization) SetLanguage(language string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if _, exists := l.translations[language]; !exists {
		return fmt.Errorf("unsupported language: %s", language)
	}

	l.currentLanguage = language
	return nil
}

// GetLanguage возвращает текущий язык
func (l *Localization) GetLanguage() string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.currentLanguage
}

// Get возвращает переведенную строку
func (l *Localization) Get(key string, args ...interface{}) string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	translations, exists := l.translations[l.currentLanguage]
	if !exists {
		// Fallback to English if available, otherwise use first available language
		if fallbackTranslations, fallbackExists := l.translations[DefaultLanguage]; fallbackExists {
			translations = fallbackTranslations
		} else if len(l.supportedLanguages) > 0 {
			translations = l.translations[l.supportedLanguages[0]]
		}
	}

	if translation, exists := translations[key]; exists {
		if len(args) > 0 {
			return fmt.Sprintf(translation, args...)
		}
		return translation
	}

	// Если перевод не найден, возвращаем ключ
	return key
}

// LoadTranslationsFromFile загружает переводы из файла
func (l *Localization) LoadTranslationsFromFile(language, filePath string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return err
	}

	if l.translations[language] == nil {
		l.translations[language] = make(map[string]string)
	}

	// Объединяем переводы
	for key, value := range translations {
		l.translations[language][key] = value
	}

	return nil
}

// SaveTranslationsToFile сохраняет переводы в файл
func (l *Localization) SaveTranslationsToFile(language, filePath string) error {
	l.mutex.RLock()
	translations, exists := l.translations[language]
	l.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("language not found: %s", language)
	}

	// Создаем директорию, если она не существует
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(translations, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// GetSupportedLanguages возвращает список поддерживаемых языков
func (l *Localization) GetSupportedLanguages() []string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	languages := make([]string, 0, len(l.translations))
	for lang := range l.translations {
		languages = append(languages, lang)
	}

	return languages
}

// Глобальные функции для удобства
func T(key string, args ...interface{}) string {
	return GetLocalization().Get(key, args...)
}

func SetLanguage(language string) error {
	return GetLocalization().SetLanguage(language)
}

func GetLanguage() string {
	return GetLocalization().GetLanguage()
}
