package handlers

import (
	"context"
	"strings"

	"github.com/MaximVod/steambotgo/internal/interfaces"
	"github.com/MaximVod/steambotgo/internal/logger"
	"github.com/MaximVod/steambotgo/internal/presenters"
	"github.com/MaximVod/steambotgo/internal/usecases"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	commandFind = "/find"
)

// TelegramHandler обрабатывает сообщения от Telegram
type TelegramHandler struct {
	multiRegionService *usecases.MultiRegionPriceService
	searchService      *usecases.SearchGamesService
	formatter          *presenters.MessageFormatter
	logger             logger.Logger
}

// NewTelegramHandler создает новый обработчик Telegram сообщений
func NewTelegramHandler(
	steamAPI interfaces.SteamAPI,
	formatter *presenters.MessageFormatter,
	logger logger.Logger,
	countries map[string]string,
	currencyRates map[string]float64,
) *TelegramHandler {
	return &TelegramHandler{
		multiRegionService: usecases.NewMultiRegionPriceService(steamAPI, countries, currencyRates),
		searchService:      usecases.NewSearchGamesService(steamAPI),
		formatter:          formatter,
		logger:             logger,
	}
}

// Handle обрабатывает обновление от Telegram
func (h *TelegramHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Проверяем, есть ли у обновления сообщение и содержит ли оно текст
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	// Проверяем, начинается ли сообщение с команды /find
	if !strings.HasPrefix(update.Message.Text, commandFind+" ") && update.Message.Text != commandFind {
		return
	}

	// Извлекаем поисковый запрос после команды /find
	query := strings.TrimPrefix(update.Message.Text, commandFind+" ")
	query = strings.TrimSpace(query)

	// Если запрос пустой (только команда), отправляем сообщение пользователю
	if query == "" {
		h.sendMessage(ctx, b, update.Message.Chat.ID, "Пожалуйста, укажите название игры после команды /find")
		return
	}

	// Валидация запроса
	if err := h.validateQuery(query); err != nil {
		h.sendMessage(ctx, b, update.Message.Chat.ID, "❌ "+err.Error())
		return
	}

	// Пытаемся получить многорегиональные цены
	prices, err := h.multiRegionService.GetMultiRegionPrices(ctx, query)
	if err != nil {
		h.logger.Error("Ошибка получения многонациональных цен", err, "query", query)
		
		// Fallback: возвращаемся к обычному поиску
		items, err := h.searchService.FetchGames(ctx, query)
		if err != nil {
			h.logger.Error("Ошибка поиска игр", err, "query", query)
			h.sendMessage(ctx, b, update.Message.Chat.ID, "Произошла ошибка при поиске игры.")
			return
		}

		h.logger.Info("Найдено игр", "count", len(items))
		message := h.formatter.FormatSteamItems(items)
		h.sendMessage(ctx, b, update.Message.Chat.ID, message)
		return
	}

	h.logger.Info("Найдены цены для игры", "game", prices.GameName, "regions", len(prices.Regions))
	message := h.formatter.FormatMultiRegionPrices(prices)
	h.sendMessage(ctx, b, update.Message.Chat.ID, message)
}

// validateQuery проверяет валидность поискового запроса
func (h *TelegramHandler) validateQuery(query string) error {
	if len(query) < 2 {
		return &ValidationError{Message: "Поисковый запрос слишком короткий (минимум 2 символа)"}
	}
	if len(query) > 200 {
		return &ValidationError{Message: "Поисковый запрос слишком длинный (максимум 200 символов)"}
	}
	return nil
}

// sendMessage отправляет сообщение пользователю
func (h *TelegramHandler) sendMessage(ctx context.Context, b *bot.Bot, chatID int64, text string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		h.logger.Error("Ошибка отправки сообщения", err, "chatID", chatID)
	}
}

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

