package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/MaximVod/steambotgo/internal/adapters"
	"github.com/MaximVod/steambotgo/internal/entities"
	"github.com/MaximVod/steambotgo/internal/usecases"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

// Отправьте любое текстовое сообщение боту после его запуска

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("Не удалось загрузить .env файл: %v", err)
	}

	// Получаем токен из переменной окружения
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Необходимо установить переменную окружения TELEGRAM_BOT_TOKEN")
	}

	// Инициализация

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Проверяем, есть ли у обновления сообщение и содержит ли оно текст
	if update.Message == nil || update.Message.Text == "" {
		// Игнорируем сообщения без текста
		return
	}

	// Проверяем, начинается ли сообщение с команды /find
	if !strings.HasPrefix(update.Message.Text, "/find ") && update.Message.Text != "/find" {
		// Игнорируем сообщения, которые не начинаются с команды /find
		return
	}

	// Извлекаем поисковый запрос после команды /find
	query := strings.TrimPrefix(update.Message.Text, "/find ")
	query = strings.TrimSpace(query) // Удаляем лишние пробелы

	// Если запрос пустой (только команда), отправляем сообщение пользователю
	if query == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Пожалуйста, укажите название игры после команды /find",
		})
		return
	}

	steamAPI := adapters.NewSteamGamesAPI()

	// Используем новый сервис многонациональных цен
	multiRegionService := usecases.NewMultiRegionPriceService(steamAPI)
	prices, err := multiRegionService.GetMultiRegionPrices(ctx, query)
	if err != nil {
		log.Printf("Ошибка получения многонациональных цен: %v", err)
		// Возвращаемся к старому поиску, если многонациональный поиск не удался
		searchService := usecases.NewSearchGamesService(steamAPI)
		items, err := searchService.FetchGames(ctx, update.Message.Text)
		if err != nil {
			log.Printf("Ошибка: %v", err)
			// Отправляем сообщение об ошибке пользователю
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Произошла ошибка при поиске игры.",
			})
			return
		}
		log.Printf("Found %d games", len(items))
		for _, item := range items[:3] { // первые 3
			price := "—"
			if item.Price != nil {
				price = fmt.Sprintf("$%.2f", float64(item.Price.Final)/100)
			}
			log.Printf("🎮 %s | %s", item.Name, price)
		}

		log.Println(items)
		message := FormatSteamItems(items)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   message,
		})
		return
	}

	log.Printf("Found prices for %s in %d countries", prices.GameName, len(prices.Regions))

	message := FormatMultiRegionPrices(prices)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   message,
	})
}

func FormatMultiRegionPrices(data *entities.MultiRegionPriceData) string {
	if len(data.Regions) == 0 {
		return "❌ Не удалось найти цены для указанной игры."
	}

	var parts []string

	// Добавляем название игры как заголовок
	parts = append(parts, fmt.Sprintf("*%s*", data.GameName))

	// Добавляем информацию о региональных ценах
	for _, region := range data.Regions {
		if region.Item.Price != nil {
			price := fmt.Sprintf("%.2f %s", float64(region.Item.Price.Final)/100, region.Item.Price.Currency)
			if region.ConvertedRub > 0 && region.CountryCode != "RU" {
				parts = append(parts, fmt.Sprintf("%s - %s (около %.0f руб)",
					region.CountryFlag, price, region.ConvertedRub))
			} else {
				parts = append(parts, fmt.Sprintf("%s - %s", region.CountryFlag, price))
			}
		} else {
			parts = append(parts, fmt.Sprintf("%s - бесплатно", region.CountryFlag))
		}
	}

	return strings.Join(parts, "\n")
}

func FormatSteamItems(items []entities.SteamItem) string {
	if len(items) == 0 {
		return "❌ Ничего не найдено."
	}

	var parts []string
	for i, item := range items {
		if i >= 5 { // не спамим — максимум 5 игр
			parts = append(parts, fmt.Sprintf("\n<i>... и ещё %d результатов</i>", len(items)-5))
			break
		}
		parts = append(parts, item.String()) // или встроить логику сюда
	}

	return strings.Join(parts, "\n\n")
}
