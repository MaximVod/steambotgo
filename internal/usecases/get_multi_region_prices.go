package usecases

import (
	"context"
	"fmt"

	"github.com/MaximVod/steambotgo/internal/entities"
	"github.com/MaximVod/steambotgo/internal/interfaces"
)

type MultiRegionPriceService struct {
	api interfaces.SteamAPI
}

func NewMultiRegionPriceService(api interfaces.SteamAPI) *MultiRegionPriceService {
	return &MultiRegionPriceService{api: api}
}

// GetMultiRegionPrices извлекает цены на игры из нескольких стран
func (s *MultiRegionPriceService) GetMultiRegionPrices(ctx context.Context, query string) (*entities.MultiRegionPriceData, error) {
	// Определяем страны, для которых мы хотим получить цены
	countries := map[string]string{
		"RU": "🇷🇺", // Россия
		"KZ": "🇰🇿", // Казахстан
		"TR": "🇹🇷", // Турция
		"PL": "🇵🇱", // Польша
	}

	data := &entities.MultiRegionPriceData{}

	// Сначала находим игру с помощью стандартного поиска (американский магазин), чтобы получить название игры
	game, err := s.api.SearchGameByQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти игру: %w", err)
	}

	if game == nil {
		// Если игра не найдена по поисковому запросу, возвращаем пустой результат
		return &entities.MultiRegionPriceData{
			GameName: query,
			Regions:  []*entities.RegionalPriceInfo{},
		}, nil
	}

	data.GameName = game.Name
	data.ID = game.ID

	// Получаем цены для каждой страны
	for countryCode, flag := range countries {
		item, err := s.api.GetGamePricesByCountryCode(ctx, query, countryCode)
		if err != nil {
			// Пропускаем эту страну, если произошла ошибка
			continue
		}

		if item != nil {
			// Рассчитываем значение в рублях
			var convertedRub float64
			if item.Price != nil {
				convertedRub = s.convertPriceToRubles(float64(item.Price.Final)/100, item.Price.Currency)
			}

			regionalPrice := &entities.RegionalPriceInfo{
				CountryCode:  countryCode,
				CountryFlag:  flag,
				Item:         item,
				ConvertedRub: convertedRub,
			}

			data.Regions = append(data.Regions, regionalPrice)
		}
	}

	return data, nil
}

// convertPriceToRubles обеспечивает приблизительную конвертацию в рубли на основе валюты
func (s *MultiRegionPriceService) convertPriceToRubles(price float64, currency string) float64 {
	// Примечание: API поиска Steam возвращает данные о ценах, которые могут не полностью отражать
	// региональные различия, так как ограничены используемым нами конечным пунктом.
	// Для получения точных региональных цен нам нужно использовать API обзора цен Steam для каждого конкретного ID приложения.

	// цена уже указана в местной валюте указанной страны
	// currency - это фактический 3-буквенный код валюты, возвращаемый API Steam (например, "RUB", "KZT", "TRY", "PLN", etc.)

	// Примерные курсы обмена для конвертации местных цен в рубли (по состоянию на декабрь 2025)
	switch currency {
	case "RUB":
		// Уже в рублях
		return price
	case "USD":
		// Конвертируем доллары США в рубли (приблизительно)
		return price * 90 // 1 USD ≈ 90 RUB
	case "EUR":
		// Конвертируем евро в рубли (приблизительно)
		return price * 99 // 1 EUR ≈ 99 RUB
	case "KZT":
		// Конвертируем казахстанские тенге в рубли (приблизительно)
		return price * 0.2 // 1 KZT ≈ 0.2 RUB (приблизительно)
	case "TRY":
		// Конвертируем турецкую лиру в рубли (приблизительно)
		return price * 2.2 // 1 TRY ≈ 2.2 RUB (приблизительно)
	case "PLN":
		// Конвертируем польские злотые в рубли (приблизительно)
		return price * 23 // 1 PLN ≈ 23 RUB (приблизительно)
	case "GBP":
		// Конвертируем британские фунты в рубли (приблизительно)
		return price * 110 // 1 GBP ≈ 110 RUB (приблизительно)
	case "CNY":
		// Конвертируем китайский юань в рубли (приблизительно)
		return price * 13 // 1 CNY ≈ 13 RUB (приблизительно)
	default:
		// Для неизвестных валют возвращаем цену как есть, но, вероятно, требуется ручная конвертация
		return price * 90 // Приблизительная оценка по умолчанию
	}
}
