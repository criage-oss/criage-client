//go:build !embed
// +build !embed

package main

// NewEmbeddedLocalization для обычной сборки repository просто создает стандартную локализацию
// (embedded функции не доступны без тега embed)
func NewEmbeddedLocalization() *Localization {
	// Без embed тега просто используем обычную локализацию
	return NewLocalization()
}

// GetEmbeddedRepositoryLanguages возвращает пустой список без embedded сборки
func GetEmbeddedRepositoryLanguages() []string {
	return []string{}
}
