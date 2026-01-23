package presenters

import (
	"fmt"
	"strings"

	"github.com/MaximVod/steambotgo/internal/entities"
)

const (
	maxSearchResults = 5
)

// MessageFormatter форматирует данные для отправки в Telegram
type MessageFormatter struct{}

// NewMessageFormatter создает новый форматтер сообщений
func NewMessageFormatter() *MessageFormatter {
	return &MessageFormatter{}
}

// FormatMultiRegionPrices форматирует данные о многорегиональных ценах
func (f *MessageFormatter) FormatMultiRegionPrices(data *entities.MultiRegionPriceData) string {
	if len(data.Regions) == 0 {
		return "❌ Не удалось найти цены для указанной игры."
	}

	var parts []string

	// Добавляем название игры как заголовок
	parts = append(parts, fmt.Sprintf("*%s*", data.GameName))

	isAllPricesNotAvailable := false

	gamePriceStatus := "Недоступно"

	// Проверяем на наличие того, есть ли хоть по одному из регионов цена
	for _, region := range data.Regions {
		if region.Item.Price != nil {
			isAllPricesNotAvailable = true
		}
	}

	if !isAllPricesNotAvailable {
		gamePriceStatus = "Бесплатно"
	}

	// Добавляем информацию о региональных ценах
	for _, region := range data.Regions {
		if region.Item.Price != nil {
			priceText := f.formatPriceText(region)
			parts = append(parts, fmt.Sprintf("%s - %s", region.CountryFlag, priceText))
		} else {
			parts = append(parts, fmt.Sprintf("%s - "+gamePriceStatus, region.CountryFlag))
		}
	}

	parts = append(parts, fmt.Sprintf("https://store.steampowered.com/app/%v", data.ID))

	return strings.Join(parts, "\n")
}

// FormatSteamItems форматирует список игр для отправки
func (f *MessageFormatter) FormatSteamItems(items []entities.SteamItem) string {
	if len(items) == 0 {
		return "❌ Ничего не найдено."
	}

	var parts []string
	for i, item := range items {
		if i >= maxSearchResults {
			parts = append(parts, fmt.Sprintf("\n<i>... и ещё %d результатов</i>", len(items)-maxSearchResults))
			break
		}
		parts = append(parts, item.String())
	}

	return strings.Join(parts, "\n\n")
}

// formatPriceText форматирует текст цены в зависимости от скидки и страны
func (f *MessageFormatter) formatPriceText(region *entities.RegionalPriceInfo) string {
	finalPrice := fmt.Sprintf("%.2f %s", float64(region.Item.Price.Final)/100, region.Item.Price.Currency)
	initialPrice := fmt.Sprintf("%.2f %s", float64(region.Item.Price.Initial)/100, region.Item.Price.Currency)

	hasDiscount := region.Item.Price.Initial > region.Item.Price.Final
	hasConversion := region.ConvertedRub > 0 && region.CountryCode != "RU"

	var text string
	if hasDiscount {
		text = fmt.Sprintf("Цена со скидкой - %s (вместо - %s)", finalPrice, initialPrice)
	} else {
		text = finalPrice
	}

	if hasConversion {
		text += fmt.Sprintf(" (около %.0f руб)", region.ConvertedRub)
	}

	return text
}
