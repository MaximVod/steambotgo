package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/MaximVod/steambotgo/internal/adapters"
	"github.com/MaximVod/steambotgo/internal/config"
	"github.com/MaximVod/steambotgo/internal/database"
	"github.com/MaximVod/steambotgo/internal/handlers"
	"github.com/MaximVod/steambotgo/internal/logger"
	"github.com/MaximVod/steambotgo/internal/presenters"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализируем логгер
	appLogger := logger.New()
	appLogger.Info("Запуск приложения")

	// Инициализируем подключение к базе данных
	appLogger.Info("Подключение к базе данных", "url", cfg.Database.URL)
	dbPool, err := database.InitDB(ctx, cfg.Database.URL)
	if err != nil {
		appLogger.Error("Ошибка подключения к БД", err)
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer database.Close(dbPool) // Закрываем соединение при завершении приложения
	appLogger.Info("Успешно подключено к базе данных")

	// Применяем миграции (создаем таблицы, если их еще нет)
	appLogger.Info("Применение миграций")
	if err := database.RunMigrations(ctx, dbPool); err != nil {
		appLogger.Error("Ошибка применения миграций", err)
		log.Fatalf("Не удалось применить миграции: %v", err)
	}
	appLogger.Info("Миграции применены успешно")

	// Инициализируем компоненты
	steamAPI := adapters.NewSteamGamesAPI(cfg.Steam.BaseURL, cfg.Steam.Timeout)
	formatter := presenters.NewMessageFormatter()
	telegramHandler := handlers.NewTelegramHandler(
		steamAPI,
		formatter,
		appLogger,
		cfg.App.SupportedCountries,
		cfg.App.CurrencyRates,
	)

	// Инициализируем бота
	opts := []bot.Option{
		bot.WithDefaultHandler(telegramHandler.Handle),
	}

	b, err := bot.New(cfg.Telegram.BotToken, opts...)
	if err != nil {
		appLogger.Error("Ошибка создания бота", err)
		log.Fatalf("Не удалось создать бота: %v", err)
	}

	appLogger.Info("Бот запущен и готов к работе")
	b.Start(ctx)
}

// loadConfig загружает конфигурацию с учетом окружения
func loadConfig() (*config.Config, error) {
	// Пытаемся загрузить .env файл
	if err := loadEnvFile(); err != nil {
		// Не критично, если .env не найден
		log.Printf("Предупреждение: не удалось загрузить .env файл: %v", err)
	}

	return config.Load()
}

// loadEnvFile пытается загрузить .env файл из разных мест
func loadEnvFile() error {
	// Пытаемся загрузить из директории исполняемого файла
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		envPath := filepath.Join(exeDir, ".env")
		if err := godotenv.Load(envPath); err == nil {
			return nil
		}
	}

	// Пытаемся загрузить из текущей директории
	return godotenv.Load()
}
