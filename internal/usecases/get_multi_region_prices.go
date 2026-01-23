package usecases

import (
	"context"
	"fmt"

	"github.com/MaximVod/steambotgo/internal/entities"
	"github.com/MaximVod/steambotgo/internal/interfaces"
)

type MultiRegionPriceService struct {
	api                interfaces.SteamAPI
	aiApi              interfaces.AiAPI
	supportedCountries map[string]string // country code -> flag emoji
	currencyRates      map[string]float64
}

func NewMultiRegionPriceService(api interfaces.SteamAPI, aiApi interfaces.AiAPI, countries map[string]string, rates map[string]float64) *MultiRegionPriceService {
	return &MultiRegionPriceService{
		api:                api,
		aiApi:              aiApi,
		supportedCountries: countries,
		currencyRates:      rates,
	}
}

// GetMultiRegionPrices извлекает цены на игры из нескольких стран
func (s *MultiRegionPriceService) GetMultiRegionPrices(ctx context.Context, query string) (*entities.MultiRegionPriceData, error) {
	data := &entities.MultiRegionPriceData{}

	var correctedQuery string

	// Сначала находим игру с помощью стандартного поиска (американский магазин)
	game, err := s.api.SearchGameByQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти игру: %w", err)
	}

	// Если игра не найдена, пытаемся использовать AI для исправления запроса
	if game == nil {
		correctedQuery, err = s.aiApi.SearchGamesByUserQuery(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("не удалось найти игру c помощью AI: %w", err)
		}

		// Пытаемся найти игру с исправленным названием
		game, err = s.api.SearchGameByQuery(ctx, correctedQuery)
		if err != nil {
			return nil, fmt.Errorf("не удалось найти игру после исправления AI: %w", err)
		}

		// Если и после AI ничего не найдено, возвращаем пустой результат
		if game == nil {
			return &entities.MultiRegionPriceData{
				GameName: correctedQuery, // Используем исправленное название, даже если не нашли
				Regions:  []*entities.RegionalPriceInfo{},
			}, nil
		}
	}

	// Устанавливаем данные игры
	data.GameName = game.Name
	data.ID = game.ID

	// Получаем цены для каждой страны
	// Пробуем разные варианты поиска для максимальной совместимости
	for countryCode, flag := range s.supportedCountries {
		var item *entities.SteamItem

		// Сначала пробуем по game.Name (точное название из Steam)
		item, err = s.api.GetGamePricesByCountryCode(ctx, game.Name, countryCode, game.ID)
		if err != nil {
			// Пропускаем эту страну, если произошла ошибка
			continue
		}

		// Если не нашли по game.Name и есть исправленный AI запрос, пробуем его
		if item == nil && correctedQuery != "" && correctedQuery != game.Name {
			item, err = s.api.GetGamePricesByCountryCode(ctx, correctedQuery, countryCode, game.ID)
			if err != nil {
				// Пропускаем эту страну, если произошла ошибка
				continue
			}
		}

		// Если все еще не нашли, пробуем оригинальный запрос
		if item == nil && query != game.Name && query != correctedQuery {
			item, err = s.api.GetGamePricesByCountryCode(ctx, query, countryCode, game.ID)
			if err != nil {
				// Пропускаем эту страну, если произошла ошибка
				continue
			}
		}

		// Если нашли игру (с нужным ID), добавляем в результат
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

	// Используем курсы из конфигурации
	rate, exists := s.currencyRates[currency]
	if !exists {
		// Для неизвестных валют используем курс USD по умолчанию
		rate = s.currencyRates["USD"]
		if rate == 0 {
			rate = 90 // Fallback значение
		}
	}

	return price * rate
}
