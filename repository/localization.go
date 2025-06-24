package main

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

// Localization управляет локализацией для repository сервера
type Localization struct {
	currentLanguage    string
	supportedLanguages []string
	translations       map[string]map[string]string
	translationsDir    string
	useEmbedded        bool
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
		// Автоматически выбираем режим локализации для repository
		embeddedLangs := GetEmbeddedRepositoryLanguages()
		if len(embeddedLangs) > 0 {
			// Если есть встроенные языки, используем их
			globalLocalization = NewEmbeddedLocalization()
		} else {
			// Иначе используем внешние файлы
			globalLocalization = NewLocalization()
		}
	})
	return globalLocalization
}

// SetGlobalLocalization устанавливает глобальный экземпляр локализации для repository
func SetGlobalLocalization(localization *Localization) {
	globalLocalization = localization
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

// getDefaultTranslations возвращает минимальный набор переводов для repository сервера
func (l *Localization) getDefaultTranslations(language string) map[string]string {
	// Базовые переводы для repository сервера
	translations := map[string]map[string]string{
		"ru": {
			"server_started":    "Сервер запущен на порту %d",
			"server_stopped":    "Сервер остановлен",
			"package_uploaded":  "Пакет загружен",
			"package_not_found": "Пакет не найден",
			"invalid_request":   "Неверный запрос",
			"internal_error":    "Внутренняя ошибка сервера",
		},
		"en": {
			"server_started":    "Server started on port %d",
			"server_stopped":    "Server stopped",
			"package_uploaded":  "Package uploaded",
			"package_not_found": "Package not found",
			"invalid_request":   "Invalid request",
			"internal_error":    "Internal server error",
		},
		"de": {
			"server_started":    "Server auf Port %d gestartet",
			"server_stopped":    "Server gestoppt",
			"package_uploaded":  "Paket hochgeladen",
			"package_not_found": "Paket nicht gefunden",
			"invalid_request":   "Ungültige Anfrage",
			"internal_error":    "Interner Serverfehler",
		},
		"fr": {
			"server_started":    "Serveur démarré sur le port %d",
			"server_stopped":    "Serveur arrêté",
			"package_uploaded":  "Paquet téléchargé",
			"package_not_found": "Paquet introuvable",
			"invalid_request":   "Requête invalide",
			"internal_error":    "Erreur interne du serveur",
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

// IsEmbedded возвращает true, если используются встроенные переводы для repository
func (l *Localization) IsEmbedded() bool {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.useEmbedded
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
