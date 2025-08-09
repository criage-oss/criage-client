<div align="center">
  <img src="logo.png" alt="Criage Logo" width="200">
  
# Criage - Высокопроизводительный пакетный менеджер
  
  Criage - это современный пакетный менеджер, написанный на Go, обеспечивающий быструю установку, обновление и управление пакетами с поддержкой различных форматов сжатия.
  
  [🇬🇧 English Version](README.md) | 🇷🇺 Русская версия
</div>

## Возможности

### Основные функции

- 🚀 **Высокая производительность** - использование быстрых алгоритмов сжатия (Zstandard, LZ4)
- 📦 **Единое расширение пакетов** - все пакеты используют расширение `.criage` с встроенными метаданными о типе сжатия
- 🔧 **Управление зависимостями** - автоматическое разрешение и установка зависимостей
- 🌐 **Множественные репозитории** - поддержка нескольких источников пакетов
- 🎯 **Кроссплатформенность** - поддержка Linux, macOS, Windows
- ⚡ **Параллельные операции** - многопоточная обработка для ускорения
- 🌍 **Многоязычная поддержка** - динамическая система локализации - см. [Руководство по локализации](LOCALIZATION.md)

### Управление пакетами

- Установка и удаление пакетов
- Обновление до последних версий  
- Поиск пакетов в репозиториях
- Просмотр информации о пакетах
- Глобальная и локальная установка

### Разработка пакетов

- Создание новых пакетов из шаблонов
- Сборка пакетов с настраиваемыми скриптами
- Публикация в репозитории
- Хуки жизненного цикла (pre/post install/remove)
- Манифесты сборки

## Установка

### Из исходников

```bash
git clone https://github.com/criage-oss/criage-client.git
cd criage-client
go build -o criage
sudo mv criage /usr/local/bin/
```

### Проверка установки

```bash
criage --version
```

## Использование

### Основные команды

#### Установка пакетов

```bash
# Установить пакет
criage install package-name

# Установить определенную версию
criage install package-name --version 1.2.3

# Установить из конкретного репозитория
criage install package-name --repo myrepo

# Глобальная установка
criage install package-name --global

# Установка с dev зависимостями
criage install package-name --dev

# Установить локальный файл .criage
criage install ./my-package-1.0.0.criage
```

#### Удаление пакетов

```bash
# Удалить пакет
criage uninstall package-name

# Полное удаление с конфигурацией
criage uninstall package-name --purge
```

#### Обновление пакетов

```bash
# Обновить конкретный пакет
criage update package-name

# Обновить все пакеты
criage update --all
```

#### Поиск и информация

```bash
# Найти пакеты во всех репозиториях
criage search keyword

# Найти пакеты в конкретном репозитории
criage search keyword --repo myrepo

# Показать все доступные пакеты
criage search "*" --all-repos

# Показать установленные пакеты
criage list

# Показать только устаревшие пакеты
criage list --outdated

# Подробная информация о пакете
criage info package-name

# Информация о пакете из конкретного репозитория
criage info package-name --repo myrepo
```

### Разработка пакетов

#### Создание нового пакета

```bash
# Создать пакет из базового шаблона
criage create my-package --author "Your Name" --description "Package description"
```

#### Сборка пакета

```bash
# Собрать с настройками по умолчанию (создаст файл .criage)
criage build

# Указать тип сжатия и уровень сжатия
criage build --format tar.zst --compression 6 --output my-package-1.0.0.criage
```

#### Публикация пакета

```bash
# Опубликовать в репозитории
criage publish --registry https://packages.example.com --token YOUR_TOKEN
```

### Управление репозиториями

#### Добавление репозиториев

```bash
# Добавить новый репозиторий
criage repo add myrepo https://packages.example.com

# Добавить репозиторий с токеном авторизации
criage repo add private-repo https://private.example.com --token YOUR_TOKEN

# Добавить репозиторий с приоритетом
criage repo add priority-repo https://priority.example.com --priority 10
```

#### Управление репозиториями

```bash
# Показать список репозиториев
criage repo list

# Показать подробную информацию о репозитории
criage repo info myrepo

# Удалить репозиторий
criage repo remove myrepo

# Обновить индексы всех репозиториев
criage repo update

# Проверить доступность репозиториев
criage repo check
```

#### Приоритет репозиториев

```bash
# Установить приоритет репозитория (чем выше число, тем выше приоритет)
criage repo priority myrepo 15

# При поиске пакетов используется следующий порядок:
# 1. Репозитории с высшим приоритетом
# 2. Официальный репозиторий (приоритет 10)
# 3. Пользовательские репозитории (приоритет 5 по умолчанию)
```

### Конфигурация

#### Просмотр настроек

```bash
# Показать все настройки
criage config list

# Получить значение конкретной настройки
criage config get cache_path
```

#### Изменение настроек

```bash
# Изменить путь кеша
criage config set cache_path /custom/cache/path

# Изменить уровень сжатия по умолчанию
criage config set compression.level 6

# Изменить количество параллельных потоков
criage config set parallel 8

# Установить репозиторий по умолчанию
criage config set default_registry https://packages.criage.ru

# Настроить тайм-аут для сетевых операций
criage config set network.timeout 30s
```

## Структура проекта

```
criage/
├── main.go              # Основная точка входа
├── commands.go          # Реализация CLI команд
├── go.mod               # Go модуль
├── go.sum               # Зависимости
└── pkg/                 # Основные пакеты
    ├── types.go         # Структуры данных
    ├── archive.go       # Работа с архивами
    ├── config.go        # Управление конфигурацией
    ├── package_manager.go        # Основная логика пакетного менеджера
    └── package_manager_helpers.go # Вспомогательные функции
```

## Форматы файлов

### Манифест пакета (criage.yaml)

```yaml
name: my-package
version: 1.0.0
description: Example package
author: Your Name
license: MIT
homepage: https://github.com/user/my-package
repository: https://github.com/user/my-package

keywords:
  - utility
  - tool

dependencies:
  some-lib: ^1.0.0

dev_dependencies:
  test-framework: ^2.0.0

scripts:
  build: make build
  test: make test
  install: make install

files:
  - "bin/*"
  - "lib/*"
  - "README.md"

exclude:
  - "*.log"
  - ".git"
  - "node_modules"

arch:
  - amd64
  - arm64

os:
  - linux  
  - darwin
  - windows

hooks:
  pre_install:
    - echo "Installing package..."
  post_install:
    - echo "Package installed successfully"
```

### Конфигурация сборки (build.json)

```json
{
  "name": "my-package",
  "version": "1.0.0",
  "build_script": "make build",
  "build_env": {
    "CGO_ENABLED": "0",
    "GOOS": "linux"
  },
  "output_dir": "./dist",
  "include_files": ["bin/*", "lib/*"],
  "exclude_files": ["*.log", "test/*"],
  "compression": {
    "format": "tar.zst",
    "level": 3
  },
  "targets": [
    {"os": "linux", "arch": "amd64"},
    {"os": "linux", "arch": "arm64"},
    {"os": "darwin", "arch": "amd64"},
    {"os": "windows", "arch": "amd64"}
  ]
}
```

### Конфигурация репозиториев

```json
{
  "repositories": [
    {
      "name": "official",
      "url": "https://packages.criage.ru",
      "priority": 10,
      "enabled": true,
      "type": "official"
    },
    {
      "name": "company-internal",
      "url": "https://packages.company.com",
      "priority": 15,
      "enabled": true,
      "type": "private",
      "auth": {
        "token": "your-company-token"
      }
    },
    {
      "name": "community",
      "url": "https://community.criage.org",
      "priority": 5,
      "enabled": true,
      "type": "community"
    }
  ],
  "cache": {
    "ttl": "1h",
    "max_size": "1GB",
    "path": "~/.criage/cache"
  },
  "network": {
    "timeout": "30s",
    "retries": 3,
    "parallel_downloads": 4
  }
}
```

## Встраивание метаданных в архивы

Criage поддерживает встраивание метаданных пакетов (`criage.yaml` и `build.json`) непосредственно в архивы. Это позволяет получать информацию о пакете без необходимости его распаковки.

### Поддерживаемые форматы

#### TAR архивы (tar.zst, tar.lz4, tar.xz, tar.gz)

- Использует **PAX Extended Headers** - стандартный механизм для хранения дополнительных метаданных
- Совместимо с большинством современных архиваторов
- Метаданные хранятся в полях `criage.metadata`, `criage.package_manifest`, `criage.build_manifest`

#### ZIP архивы

- Использует **ZIP Comment** для основных метаданных
- Дополнительно создает файл `.criage_metadata.json` внутри архива
- Полная обратная совместимость

### Встраиваемые данные

- **Манифест пакета** (`criage.yaml`) - название, версия, зависимости, автор
- **Манифест сборки** (`build.json`) - настройки сборки, целевые платформы
- **Тип сжатия** - формат и уровень сжатия
- **Метаданные создания** - дата, версия criage
- **Контрольные суммы** - для проверки целостности

### Примеры использования

#### Создание архива с метаданными

```bash
# Собрать пакет с автоматическим встраиванием метаданных
criage build --format tar.zst --compression 6

# Результат: test-package-1.0.0.criage с встроенными метаданными
```

#### Просмотр метаданных архива

```bash
# Показать все метаданные архива
criage metadata test-package-1.0.0.criage

# Пример вывода:
# === Метаданные архива test-package-1.0.0.criage ===
# Тип сжатия: tar.zst
# Создан: 2024-01-15T10:30:45Z
# Создано с помощью: criage/1.0.0
# 
# === Манифест пакета ===
# Название: test-package
# Версия: 1.0.0
# Описание: Тестовый пакет
# Автор: Developer Name
# Лицензия: MIT
# Зависимости:
#   - some-lib: ^1.0.0
# 
# === Манифест сборки ===
# Скрипт сборки: echo Building...
# Выходная директория: ./build
# Формат сжатия: tar.zst (уровень 6)
# Целевые платформы:
#   - linux/amd64
#   - linux/arm64
```

### Преимущества встраивания метаданных

1. **Самодостаточность** - архив содержит всю необходимую информацию
2. **Быстрый доступ** - не нужно распаковывать для получения информации
3. **Стандартность** - использует стандартные механизмы архивных форматов
4. **Совместимость** - работает с любыми архиваторами, поддерживающими PAX
5. **Безопасность** - встроенные контрольные суммы для проверки целостности

### Технические детали

#### Структура метаданных

```json
{
  "package_manifest": {
    "name": "test-package",
    "version": "1.0.0",
    "dependencies": {...}
  },
  "build_manifest": {
    "build_script": "echo Building...",
    "compression": {...}
  },
  "compression_type": "tar.zst",
  "created_at": "2024-01-15T10:30:45Z",
  "created_by": "criage/1.0.0",
  "checksum": "sha256:..."
}
```

#### Расположение в архиве

- **TAR**: PAX Extended Headers в начале архива
- **ZIP**: Комментарий архива + файл `.criage_metadata.json`

## Производительность

Criage оптимизирован для максимальной производительности:

- **Zstandard сжатие** - до 3x быстрее чем gzip при лучшем сжатии
- **LZ4 сжатие** - экстремально быстрое сжатие/распаковка
- **Параллельная обработка** - использование всех доступных CPU ядер
- **Умное кеширование** - избежание повторных загрузок
- **Эффективное разрешение зависимостей** - минимизация сетевых запросов

## Сравнение форматов сжатия

| Формат | Скорость сжатия | Скорость распаковки | Размер | Использование |
|--------|----------------|---------------------|--------|---------------|
| tar.zst | Средняя | Очень быстрая | Отличное | По умолчанию |
| tar.lz4 | Очень быстрая | Очень быстрая | Среднее | Быстрые операции |
| tar.xz | Медленная | Средняя | Отличное | Минимальный размер |
| tar.gz | Средняя | Средняя | Хорошее | Совместимость |
| zip | Средняя | Быстрая | Хорошее | Windows совместимость |

## Разработка

### Требования

- Go 1.24.4 или выше
- Git

### Сборка из исходников

```bash
git clone https://github.com/criage-oss/criage-client.git
cd criage-client
go mod tidy
go build -o criage
```

### Запуск тестов

```bash
go test ./...
```

### Форматирование кода

```bash
go fmt ./...
```

## Лицензия

MIT License - см. файл [LICENSE](LICENSE) для подробностей.

## Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для новой функции (`git checkout -b feature/amazing-feature`)
3. Зафиксируйте изменения (`git commit -m 'Add amazing feature'`)
4. Отправьте в ветку (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## Запуск собственного репозитория

Criage поддерживает создание собственных репозиториев для частного использования или организаций.

### Быстрый старт репозитория

```bash
# Клонировать проект
git clone https://github.com/criage-oss/criage-client.git
cd criage-client/repository

# Собрать сервер репозитория
go build -o criage-repository

# Запустить с конфигурацией по умолчанию
./criage-repository
```

### Конфигурация сервера

Отредактируйте `config.json`:

```json
{
  "port": 8081,
  "storage_path": "./packages",
  "upload_token": "your-secure-token",
  "allowed_formats": ["criage", "tar.zst", "tar.lz4"],
  "enable_cors": true
}
```

### Загрузка пакетов в репозиторий

```bash
# Загрузить пакет через API
curl -X POST http://localhost:8081/api/v1/upload \
  -H "Authorization: Bearer your-secure-token" \
  -F "file=@my-package-1.0.0.criage"

# Или скопировать файл в папку packages/
cp my-package-1.0.0.criage ./packages/

# Обновить индекс
curl -X POST http://localhost:8081/api/v1/refresh \
  -H "Authorization: Bearer your-secure-token"
```

### Использование собственного репозитория

```bash
# Добавить репозиторий
criage repo add mycompany http://localhost:8081

# Установить пакеты из своего репозитория
criage install my-package --repo mycompany
```

## Поддержка

- 📧 Email: <support@criage.ru>
- 🐛 Баги: <https://github.com/criage-oss/criage-client/issues>
- 📖 Документация: <https://docs.criage.ru>
