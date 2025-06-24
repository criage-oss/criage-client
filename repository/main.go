package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Инициализируем локализацию repository (автоматически выбирает embedded/внешние файлы)
	l := GetLocalization()

	// Показываем режим локализации в debug режиме
	if os.Getenv("CRIAGE_DEBUG") == "1" {
		embeddedLangs := GetEmbeddedRepositoryLanguages()
		if len(embeddedLangs) > 0 {
			log.Printf("🚀 Repository: используются встроенные переводы: %v", embeddedLangs)
		} else {
			log.Printf("📁 Repository: используются внешние файлы переводов: %v", l.GetSupportedLanguages())
		}
	}

	// Параметры командной строки
	configPath := flag.String("config", "config.json", "Path to configuration file")
	flag.Parse()

	// Загружаем конфигурацию
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаем необходимые директории
	if err := createDirectories(config); err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	// Инициализируем менеджер индекса
	indexManager, err := NewIndexManager(config)
	if err != nil {
		log.Fatalf("Failed to create index manager: %v", err)
	}

	// Создаем API сервер
	apiServer := NewApiServer(indexManager, config)

	// Запускаем сервер
	log.Fatal(apiServer.Start())
}

// loadConfig загружает конфигурацию из файла
func loadConfig(configPath string) (*Config, error) {
	// Создаем конфигурацию по умолчанию
	config := &Config{
		Port:        8080,
		StoragePath: "./packages",
		IndexPath:   "./index.json",
		UploadToken: "your-secret-token",
		MaxFileSize: 100 * 1024 * 1024, // 100MB
		AllowedFormats: []string{
			"criage", "tar.zst", "tar.lz4", "tar.xz", "tar.gz", "zip",
		},
		EnableCORS: true,
		LogLevel:   "info",
	}

	// Если файл конфигурации существует, загружаем его
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, config); err != nil {
			return nil, err
		}
	} else {
		// Создаем файл конфигурации по умолчанию
		log.Printf("Creating default config at %s", configPath)
		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return nil, err
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// createDirectories создает необходимые директории
func createDirectories(config *Config) error {
	// Создаем директорию для пакетов
	if err := os.MkdirAll(config.StoragePath, 0755); err != nil {
		return err
	}

	// Создаем директорию для веб интерфейса
	webDir := "./web"
	if err := os.MkdirAll(webDir, 0755); err != nil {
		return err
	}

	// Создаем простой index.html если его нет
	indexHTML := filepath.Join(webDir, "index.html")
	if _, err := os.Stat(indexHTML); os.IsNotExist(err) {
		html := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Criage Package Repository</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { color: #333; }
        .info { background: #f5f5f5; padding: 20px; border-radius: 5px; }
        .api-endpoint { margin: 10px 0; }
        .method { font-weight: bold; color: #007bff; }
        code { background: #e9ecef; padding: 2px 4px; border-radius: 3px; }
    </style>
</head>
<body>
    <h1>🚀 Criage Package Repository</h1>
    
    <div class="info">
        <h2>API Endpoints</h2>
        
        <div class="api-endpoint">
            <span class="method">GET</span> <code>/api/v1/</code> - Информация о репозитории
        </div>
        
        <div class="api-endpoint">
            <span class="method">GET</span> <code>/api/v1/packages</code> - Список всех пакетов
        </div>
        
        <div class="api-endpoint">
            <span class="method">GET</span> <code>/api/v1/packages/{name}</code> - Информация о пакете
        </div>
        
        <div class="api-endpoint">
            <span class="method">GET</span> <code>/api/v1/search?q={query}</code> - Поиск пакетов
        </div>
        
        <div class="api-endpoint">
            <span class="method">GET</span> <code>/api/v1/download/{name}/{version}/{file}</code> - Скачать пакет
        </div>
        
        <div class="api-endpoint">
            <span class="method">POST</span> <code>/api/v1/upload</code> - Загрузить пакет (требует токен)
        </div>
    </div>

    <h2>📊 Статистика</h2>
    <div id="stats">Загрузка...</div>

    <script>
        // Загружаем статистику репозитория
        fetch('/api/v1/')
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    document.getElementById('stats').innerHTML = 
                        '<p><strong>Всего пакетов:</strong> ' + data.data.total_packages + '</p>' +
                        '<p><strong>Последнее обновление:</strong> ' + new Date(data.data.last_updated).toLocaleString() + '</p>' +
                        '<p><strong>Поддерживаемые форматы:</strong> ' + data.data.formats.join(', ') + '</p>';
                }
            })
            .catch(error => {
                document.getElementById('stats').innerHTML = '<p>Ошибка загрузки статистики</p>';
            });
    </script>
</body>
</html>`
		if err := os.WriteFile(indexHTML, []byte(html), 0644); err != nil {
			return err
		}
	}

	return nil
}
