//go:build !embed
// +build !embed

package pkg

// NewEmbeddedLocalization для обычной сборки просто создает стандартную локализацию
// (embedded функции не доступны без тега embed)
func NewEmbeddedLocalization() *Localization {
	// Без embed тега просто используем обычную локализацию
	return NewLocalization()
}

// GetEmbeddedLanguages возвращает пустой список без embedded сборки
func GetEmbeddedLanguages() []string {
	return []string{}
}
