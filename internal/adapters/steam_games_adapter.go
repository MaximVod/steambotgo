package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/MaximVod/steambotgo/internal/entities"
	"github.com/MaximVod/steambotgo/internal/interfaces"
)

type SteamGamesAPI struct {
	baseURL string
	client  *http.Client
}

func NewSteamGamesAPI(baseURL string, timeout time.Duration) *SteamGamesAPI {
	return &SteamGamesAPI{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// SearchGamesByName реализует interfaces.SteamAPI.
func (f *SteamGamesAPI) SearchGamesByName(ctx context.Context, query string) ([]entities.SteamItem, error) {
	// Формируем URL с экранированием query
	endpoint := fmt.Sprintf(
		"%s/api/storesearch/?term=%s&l=english&cc=US",
		f.baseURL,
		url.QueryEscape(query), // ← защищает от " ", "&", "%"
	)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, errors.New("поиск отменен")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.New("поиск превысил время ожидания")
		}
		return nil, fmt.Errorf("HTTP запрос не удался: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус %d", resp.StatusCode)
	}

	// Парсим ответ
	var result entities.SteamResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("не удалось декодировать JSON: %w", err)
	}

	return result.Items, nil
}

// SearchGameByQuery реализует interfaces.SteamAPI.
func (f *SteamGamesAPI) SearchGameByQuery(ctx context.Context, query string) (*entities.SteamItem, error) {
	items, err := f.SearchGamesByName(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	return &items[0], nil
}

// GetGamePricesByCountryCode реализует interfaces.SteamAPI.
func (f *SteamGamesAPI) GetGamePricesByCountryCode(ctx context.Context, query string, countryCode string) (*entities.SteamItem, error) {
	// Формируем URL с экранированием query и указанием страны
	endpoint := fmt.Sprintf(
		"%s/api/storesearch/?term=%s&l=english&cc=%s",
		f.baseURL,
		url.QueryEscape(query),
		countryCode,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать запрос: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, errors.New("поиск отменен")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.New("поиск превысил время ожидания")
		}
		return nil, fmt.Errorf("HTTP запрос не удался: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус %d", resp.StatusCode)
	}

	// Парсим ответ
	var result entities.SteamResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("не удалось декодировать JSON: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	return &result.Items[0], nil
}

// Компиляторная проверка реализации интерфейса.
var _ interfaces.SteamAPI = (*SteamGamesAPI)(nil)
