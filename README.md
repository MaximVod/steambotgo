# SteamBotGo

A Go application for interacting with Steam store APIs.

## Features
- Search for games on Steam
- Get pricing information for games across different regions
- Clean architecture with separation of concerns

## Architecture
This project follows a clean architecture pattern with:
- **cmd**: Entry point of the application
- **internal/adapters**: External service adapters (Steam API)
- **internal/entities**: Domain entities
- **internal/interfaces**: Port interfaces
- **internal/usecases**: Business logic

## Setup

### Установка зависимостей
```bash
go mod tidy
```

### Конфигурация

#### Основной бот
Создайте файл `.env` в корне проекта:
```env
TELEGRAM_BOT_TOKEN=your_production_bot_token
OPENAI_API_KEY=your_openai_api_key
```

#### Тестовый бот
Создайте файл `.env.test` в корне проекта:
```env
TELEGRAM_BOT_TOKEN_TEST=your_test_bot_token
OPENAI_API_KEY=your_openai_api_key
```

**Важно:** Файл `.env.test` должен быть в `.gitignore` (уже добавлен). 
Приложение автоматически определит тестовый режим, если файл `.env.test` существует локально 
или если установлена переменная окружения `USE_TEST_BOT=true`.

## Запуск

### Локальный запуск

#### Основной бот
```bash
make run
# или
go run ./cmd/main.go
```

#### Тестовый бот
```bash
make run-test
```

При запуске тестового бота приложение автоматически использует токен из переменной `TELEGRAM_BOT_TOKEN_TEST`, если она установлена.

### Запуск через Docker

#### Основной бот
```bash
docker-compose up --build
```

#### Тестовый бот
```bash
make docker-test
# или
docker-compose -f docker-compose.test.yml up --build
```

### Остановка

```bash
# Остановить основной бот
make stop

# Остановить тестовый бот
make stop-test
```

## Доступные команды Make

- `make run` - Запустить основной бот локально
- `make run-test` - Запустить тестовый бот локально
- `make build` - Собрать приложение
- `make docker-test` - Запустить тестовый бот через Docker
- `make stop-test` - Остановить тестовый бот в Docker
- `make stop` - Остановить основной бот в Docker
- `make help` - Показать справку

## License
MIT
