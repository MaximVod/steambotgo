package config

import (
	"fmt"
	"os"
	"time"
)

// Config содержит всю конфигурацию приложения
type Config struct {
	Telegram TelegramConfig
	Steam    SteamConfig
	App      AppConfig
}

// TelegramConfig содержит настройки Telegram бота
type TelegramConfig struct {
	BotToken string
}

// SteamConfig содержит настройки для работы с Steam API
type SteamConfig struct {
	BaseURL    string
	Timeout    time.Duration
	MaxRetries int
}

// AppConfig содержит общие настройки приложения
type AppConfig struct {
	MaxSearchResults    int
	MaxRegionResults    int
	SupportedCountries  map[string]string // country code -> flag emoji
	CurrencyRates       map[string]float64 // currency code -> rate to RUB
	IsRailway           bool
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	cfg := &Config{
		Telegram: TelegramConfig{
			BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		},
		Steam: SteamConfig{
			BaseURL:    getEnvOrDefault("STEAM_BASE_URL", "https://store.steampowered.com"),
			Timeout:    10 * time.Second,
			MaxRetries: 3,
		},
		App: AppConfig{
			MaxSearchResults: 5,
			MaxRegionResults: 10,
			SupportedCountries: map[string]string{
				"RU": "🇷🇺", // Россия
				"KZ": "🇰🇿", // Казахстан
				"TR": "🇹🇷", // Турция
				"PL": "🇵🇱", // Польша
			},
			CurrencyRates: map[string]float64{
				"RUB": 1.0,   // Уже в рублях
				"USD": 90.0,  // 1 USD ≈ 90 RUB
				"EUR": 99.0,  // 1 EUR ≈ 99 RUB
				"KZT": 0.2,   // 1 KZT ≈ 0.2 RUB
				"TRY": 2.2,   // 1 TRY ≈ 2.2 RUB
				"PLN": 23.0,  // 1 PLN ≈ 23 RUB
				"GBP": 110.0, // 1 GBP ≈ 110 RUB
				"CNY": 13.0,  // 1 CNY ≈ 13 RUB
			},
			IsRailway: os.Getenv("RAILWAY") != "",
		},
	}

	if cfg.Telegram.BotToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN не установлен")
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

