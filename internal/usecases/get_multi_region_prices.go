package usecases

import (
	"context"
	"fmt"

	"github.com/MaximVod/steambotgo/internal/entities"
	"github.com/MaximVod/steambotgo/internal/interfaces"
)

type MultiRegionPriceService struct {
	api              interfaces.SteamAPI
	supportedCountries map[string]string // country code -> flag emoji
	currencyRates    map[string]float64   // currency code -> rate to RUB
}

func NewMultiRegionPriceService(api interfaces.SteamAPI, countries map[string]string, rates map[string]float64) *MultiRegionPriceService {
	return &MultiRegionPriceService{
		api:                api,
		supportedCountries: countries,
		currencyRates:      rates,
	}
}

// GetMultiRegionPrices извлекает цены на игры из нескольких стран
func (s *MultiRegionPriceService) GetMultiRegionPrices(ctx context.Context, query string) (*entities.MultiRegionPriceData, error) {

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
	for countryCode, flag := range s.supportedCountries {
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
