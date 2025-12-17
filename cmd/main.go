package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/MaximVod/steambotgo/internal/adapters"
	"github.com/MaximVod/steambotgo/internal/config"
	"github.com/MaximVod/steambotgo/internal/handlers"
	"github.com/MaximVod/steambotgo/internal/logger"
	"github.com/MaximVod/steambotgo/internal/presenters"
	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Загружаем .env файл если не на Railway
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализируем логгер
	appLogger := logger.New()
	appLogger.Info("Запуск приложения")

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
	// Если не на Railway, пытаемся загрузить .env файл
	if os.Getenv("RAILWAY") == "" {
		if err := loadEnvFile(); err != nil {
			// Не критично, если .env не найден
			log.Printf("Предупреждение: не удалось загрузить .env файл: %v", err)
		}
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
