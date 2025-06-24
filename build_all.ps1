# Универсальный PowerShell скрипт сборки всего проекта Criage
# Собирает клиент и repository сервер с встроенной локализацией

$ErrorActionPreference = "Stop"

Write-Host "🏗️  Полная сборка проекта Criage с помощью собственного пакетного менеджера" -ForegroundColor Green
Write-Host "======================================================================" -ForegroundColor Green

# Функция для красивого вывода времени
function Write-TimeStamp {
    param([string]$Message)
    $timestamp = Get-Date -Format "HH:mm:ss"
    Write-Host "⏰ $timestamp`: $Message" -ForegroundColor Yellow
}

# Функция для проверки успешности операции
function Test-Success {
    param([string]$Operation)
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ $Operation - успешно!" -ForegroundColor Green
    } else {
        Write-Host "❌ $Operation - ошибка!" -ForegroundColor Red
        exit 1
    }
}

Write-TimeStamp "Начало сборки"

# 1. Проверяем и собираем основной пакетный менеджер
Write-Host "📦 Этап 1: Подготовка пакетного менеджера..." -ForegroundColor Cyan
if (-not (Test-Path "./criage.exe")) {
    Write-TimeStamp "Сборка пакетного менеджера"
    go build -o criage.exe .
    Test-Success "Сборка пакетного менеджера"
} else {
    Write-Host "✅ Пакетный менеджер уже существует" -ForegroundColor Green
}

# 2. Сборка клиента
Write-Host ""
Write-Host "🖥️  Этап 2: Сборка клиента с встроенной локализацией..." -ForegroundColor Cyan
Write-TimeStamp "Начало сборки клиента"

# Показываем информацию о клиенте
Write-Host "📋 Информация о пакете клиента:" -ForegroundColor White
Write-Host "  • Имя: criage" -ForegroundColor Gray
Write-Host "  • Версия: 1.0.0" -ForegroundColor Gray
Write-Host "  • Тип: Пакетный менеджер с embedded локализацией" -ForegroundColor Gray

.\criage.exe build -o criage-client-embedded.tar.zst -f tar.zst -c 6
Test-Success "Сборка клиента"

if (Test-Path "criage-client-embedded.tar.zst") {
    $clientSize = (Get-Item "criage-client-embedded.tar.zst").Length
    $clientSizeFormatted = [math]::Round($clientSize / 1MB, 2)
    Write-Host "📏 Размер архива клиента: $clientSizeFormatted MB" -ForegroundColor White
}

# 3. Сборка repository сервера
Write-Host ""
Write-Host "🌐 Этап 3: Сборка repository сервера с встроенной локализацией..." -ForegroundColor Cyan
Write-TimeStamp "Начало сборки repository"

# Переходим в директорию repository
Push-Location repository

# Показываем информацию о repository
Write-Host "📋 Информация о пакете repository:" -ForegroundColor White
Write-Host "  • Имя: criage-repository" -ForegroundColor Gray
Write-Host "  • Версия: 1.0.0" -ForegroundColor Gray
Write-Host "  • Тип: Repository сервер с embedded локализацией" -ForegroundColor Gray

..\criage.exe build -o criage-repository-embedded.tar.zst -f tar.zst -c 6
Test-Success "Сборка repository сервера"

if (Test-Path "criage-repository-embedded.tar.zst") {
    $repoSize = (Get-Item "criage-repository-embedded.tar.zst").Length
    $repoSizeFormatted = [math]::Round($repoSize / 1MB, 2)
    Write-Host "📏 Размер архива repository: $repoSizeFormatted MB" -ForegroundColor White
}

# Возвращаемся в корневую директорию
Pop-Location

# 4. Итоговая информация
Write-Host ""
Write-Host "🎉 СБОРКА ЗАВЕРШЕНА УСПЕШНО!" -ForegroundColor Green
Write-Host "======================================================================" -ForegroundColor Green
Write-TimeStamp "Конец сборки"

Write-Host ""
Write-Host "📦 Созданные архивы:" -ForegroundColor White
if (Test-Path "criage-client-embedded.tar.zst") {
    $clientSize = (Get-Item "criage-client-embedded.tar.zst").Length
    $clientSizeFormatted = [math]::Round($clientSize / 1MB, 2)
    Write-Host "  • criage-client-embedded.tar.zst ($clientSizeFormatted MB)" -ForegroundColor Gray
}

if (Test-Path "repository/criage-repository-embedded.tar.zst") {
    $repoSize = (Get-Item "repository/criage-repository-embedded.tar.zst").Length
    $repoSizeFormatted = [math]::Round($repoSize / 1MB, 2)
    Write-Host "  • repository/criage-repository-embedded.tar.zst ($repoSizeFormatted MB)" -ForegroundColor Gray
}

Write-Host ""
Write-Host "🚀 Готово к развертыванию:" -ForegroundColor Green
Write-Host "  • Клиент: полностью автономный пакетный менеджер" -ForegroundColor Gray
Write-Host "  • Сервер: репозиторий с веб-интерфейсом" -ForegroundColor Gray
Write-Host "  • Локализация: встроена в оба архива (ru, en, de, es)" -ForegroundColor Gray
Write-Host ""
Write-Host "💡 Для тестирования используйте команды:" -ForegroundColor Yellow
Write-Host "  .\criage.exe metadata criage-client-embedded.tar.zst" -ForegroundColor White
Write-Host "  .\criage.exe metadata repository\criage-repository-embedded.tar.zst" -ForegroundColor White 