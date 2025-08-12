package pkg

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// TestApiResponseStructure проверяет соответствие структуры ApiResponse схеме API
func TestApiResponseStructure(t *testing.T) {
	// Проверяем, что структура ApiResponse содержит все необходимые поля
	response := ApiResponse{
		Success: true,
		Message: "Test message",
		Data:    map[string]interface{}{"test": "data"},
		Error:   "Test error",
	}

	// Сериализуем и десериализуем для проверки корректности JSON тегов
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal ApiResponse: %v", err)
	}

	var unmarshaled ApiResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ApiResponse: %v", err)
	}

	if unmarshaled.Success != response.Success {
		t.Errorf("Success field mismatch: expected %v, got %v", response.Success, unmarshaled.Success)
	}

	if unmarshaled.Message != response.Message {
		t.Errorf("Message field mismatch: expected %s, got %s", response.Message, unmarshaled.Message)
	}

	if unmarshaled.Error != response.Error {
		t.Errorf("Error field mismatch: expected %s, got %s", response.Error, unmarshaled.Error)
	}
}

// TestPackageEntryStructure проверяет соответствие структуры PackageEntry схеме API
func TestPackageEntryStructure(t *testing.T) {
	now := time.Now()
	packageEntry := PackageEntry{
		Name:          "test-package",
		Description:   "Test package description",
		Author:        "Test Author",
		License:       "MIT",
		Homepage:      "https://example.com",
		Repository:    "https://github.com/example/test",
		Keywords:      []string{"test", "example"},
		Versions:      []VersionEntry{},
		LatestVersion: "1.0.0",
		Downloads:     100,
		Updated:       now,
	}

	// Проверяем сериализацию/десериализацию
	data, err := json.Marshal(packageEntry)
	if err != nil {
		t.Fatalf("Failed to marshal PackageEntry: %v", err)
	}

	var unmarshaled PackageEntry
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal PackageEntry: %v", err)
	}

	if unmarshaled.Name != packageEntry.Name {
		t.Errorf("Name field mismatch: expected %s, got %s", packageEntry.Name, unmarshaled.Name)
	}

	if len(unmarshaled.Keywords) != len(packageEntry.Keywords) {
		t.Errorf("Keywords length mismatch: expected %d, got %d", len(packageEntry.Keywords), len(unmarshaled.Keywords))
	}
}

// TestVersionEntryStructure проверяет соответствие структуры VersionEntry схеме API
func TestVersionEntryStructure(t *testing.T) {
	now := time.Now()
	versionEntry := VersionEntry{
		Version:      "1.0.0",
		Description:  "Initial version",
		Dependencies: map[string]string{"dep1": "^1.0.0"},
		DevDeps:      map[string]string{"devdep1": "^2.0.0"},
		Files:        []FileEntry{},
		Size:         1024,
		Checksum:     "sha256:abcd1234",
		Uploaded:     now,
		Downloads:    50,
	}

	// Проверяем сериализацию/десериализацию
	data, err := json.Marshal(versionEntry)
	if err != nil {
		t.Fatalf("Failed to marshal VersionEntry: %v", err)
	}

	var unmarshaled VersionEntry
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal VersionEntry: %v", err)
	}

	if unmarshaled.Version != versionEntry.Version {
		t.Errorf("Version field mismatch: expected %s, got %s", versionEntry.Version, unmarshaled.Version)
	}

	if unmarshaled.Size != versionEntry.Size {
		t.Errorf("Size field mismatch: expected %d, got %d", versionEntry.Size, unmarshaled.Size)
	}
}

// TestFileEntryStructure проверяет соответствие структуры FileEntry схеме API
func TestFileEntryStructure(t *testing.T) {
	fileEntry := FileEntry{
		OS:       "linux",
		Arch:     "amd64",
		Format:   "tar.zst",
		Filename: "test-package-1.0.0-linux-amd64.tar.zst",
		Size:     2048,
		Checksum: "sha256:efgh5678",
	}

	// Проверяем сериализацию/десериализацию
	data, err := json.Marshal(fileEntry)
	if err != nil {
		t.Fatalf("Failed to marshal FileEntry: %v", err)
	}

	var unmarshaled FileEntry
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal FileEntry: %v", err)
	}

	if unmarshaled.OS != fileEntry.OS {
		t.Errorf("OS field mismatch: expected %s, got %s", fileEntry.OS, unmarshaled.OS)
	}

	if unmarshaled.Arch != fileEntry.Arch {
		t.Errorf("Arch field mismatch: expected %s, got %s", fileEntry.Arch, unmarshaled.Arch)
	}
}

// TestSearchResultStructure проверяет соответствие структуры SearchResult схеме API
func TestSearchResultStructure(t *testing.T) {
	now := time.Now()
	searchResult := SearchResult{
		Name:        "search-test",
		Version:     "1.0.0",
		Description: "Search test package",
		Author:      "Search Author",
		Downloads:   75,
		Updated:     now,
		Score:       0.95,
	}

	// Проверяем сериализацию/десериализацию
	data, err := json.Marshal(searchResult)
	if err != nil {
		t.Fatalf("Failed to marshal SearchResult: %v", err)
	}

	var unmarshaled SearchResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal SearchResult: %v", err)
	}

	if unmarshaled.Name != searchResult.Name {
		t.Errorf("Name field mismatch: expected %s, got %s", searchResult.Name, unmarshaled.Name)
	}

	if unmarshaled.Score != searchResult.Score {
		t.Errorf("Score field mismatch: expected %f, got %f", searchResult.Score, unmarshaled.Score)
	}
}

// TestPackageListResponseStructure проверяет новую структуру для списка пакетов
func TestPackageListResponseStructure(t *testing.T) {
	packageList := PackageListResponse{
		Packages: []*PackageEntry{
			{
				Name:          "test1",
				LatestVersion: "1.0.0",
			},
			{
				Name:          "test2",
				LatestVersion: "2.0.0",
			},
		},
		Total:      100,
		Page:       1,
		Limit:      20,
		TotalPages: 5,
	}

	// Проверяем сериализацию/десериализацию
	data, err := json.Marshal(packageList)
	if err != nil {
		t.Fatalf("Failed to marshal PackageListResponse: %v", err)
	}

	var unmarshaled PackageListResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal PackageListResponse: %v", err)
	}

	if len(unmarshaled.Packages) != len(packageList.Packages) {
		t.Errorf("Packages length mismatch: expected %d, got %d", len(packageList.Packages), len(unmarshaled.Packages))
	}

	if unmarshaled.Total != packageList.Total {
		t.Errorf("Total field mismatch: expected %d, got %d", packageList.Total, unmarshaled.Total)
	}
}

// TestApiSchemaCompatibility проверяет совместимость с API схемой из документации
func TestApiSchemaCompatibility(t *testing.T) {
	// Проверяем, что все основные типы данных имеют правильные JSON теги
	testCases := []struct {
		name       string
		structType reflect.Type
		fieldTests map[string]string // поле -> ожидаемый JSON тег
	}{
		{
			name:       "ApiResponse",
			structType: reflect.TypeOf(ApiResponse{}),
			fieldTests: map[string]string{
				"Success": "success",
				"Message": "message,omitempty",
				"Data":    "data,omitempty",
				"Error":   "error,omitempty",
			},
		},
		{
			name:       "PackageEntry",
			structType: reflect.TypeOf(PackageEntry{}),
			fieldTests: map[string]string{
				"Name":          "name",
				"Description":   "description",
				"Author":        "author",
				"License":       "license",
				"LatestVersion": "latestVersion",
				"Downloads":     "downloads",
			},
		},
		{
			name:       "VersionEntry",
			structType: reflect.TypeOf(VersionEntry{}),
			fieldTests: map[string]string{
				"Version":      "version",
				"Dependencies": "dependencies,omitempty",
				"DevDeps":      "devDependencies,omitempty",
				"Size":         "size",
				"Checksum":     "checksum",
				"Downloads":    "downloads",
			},
		},
		{
			name:       "FileEntry",
			structType: reflect.TypeOf(FileEntry{}),
			fieldTests: map[string]string{
				"OS":       "os",
				"Arch":     "arch",
				"Format":   "format",
				"Filename": "filename",
				"Size":     "size",
				"Checksum": "checksum",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for fieldName, expectedTag := range tc.fieldTests {
				field, found := tc.structType.FieldByName(fieldName)
				if !found {
					t.Errorf("Field %s not found in %s", fieldName, tc.name)
					continue
				}

				jsonTag := field.Tag.Get("json")
				if jsonTag != expectedTag {
					t.Errorf("Field %s in %s has wrong JSON tag: expected %s, got %s", fieldName, tc.name, expectedTag, jsonTag)
				}
			}
		})
	}
}

// TestRateLimiterFunctionality проверяет работу rate limiter
func TestRateLimiterFunctionality(t *testing.T) {
	// Создаем rate limiter с высокой частотой для быстрого тестирования
	rl := NewRateLimiter(100) // 100 запросов в секунду
	defer rl.Close()

	// Проверяем, что rate limiter не блокирует нормальные запросы
	start := time.Now()
	for i := 0; i < 5; i++ {
		rl.Wait()
	}
	elapsed := time.Since(start)

	// Должно занимать меньше секунды для 5 запросов при лимите 100/сек
	if elapsed > time.Second {
		t.Errorf("Rate limiter is too slow: took %v for 5 requests", elapsed)
	}

	// Проверяем, что rate limiter действительно ограничивает частоту
	rl2 := NewRateLimiter(2) // 2 запроса в секунду
	defer rl2.Close()

	start = time.Now()
	for i := 0; i < 3; i++ {
		rl2.Wait()
	}
	elapsed = time.Since(start)

	// Должно занимать как минимум 1 секунду для 3 запросов при лимите 2/сек
	if elapsed < time.Second {
		t.Errorf("Rate limiter is not working: took only %v for 3 requests with 2/sec limit", elapsed)
	}
}

// BenchmarkRateLimiter бенчмарк для rate limiter
func BenchmarkRateLimiter(b *testing.B) {
	rl := NewRateLimiter(1000) // 1000 запросов в секунду
	defer rl.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.Wait()
	}
}

// TestNewApiEndpoints проверяет новые эндпоинты API
func TestNewApiEndpoints(t *testing.T) {
	pm, err := NewPackageManager()
	if err != nil {
		t.Skipf("Failed to create PackageManager: %v", err)
	}

	// Проверяем, что методы существуют (компиляция пройдет только если методы определены)
	// Вызываем методы с пустыми параметрами для проверки их наличия
	_, err = pm.ListRepositoryPackages("", 1, 10)
	if err == nil {
		t.Log("ListRepositoryPackages method is available")
	}

	_, err = pm.GetPackageVersion("", "", "")
	if err == nil {
		t.Log("GetPackageVersion method is available")
	}

	// Проверяем, что rate limiter инициализирован
	if pm.rateLimiter == nil {
		t.Error("Rate limiter is not initialized")
	}
}
