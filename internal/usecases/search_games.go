package usecases

import (
	"context"
	"fmt"

	"github.com/MaximVod/steambotgo/internal/entities"
	"github.com/MaximVod/steambotgo/internal/interfaces"
)

type SearchGamesService struct {
	api interfaces.SteamAPI
}

func NewSearchGamesService(api interfaces.SteamAPI) *SearchGamesService {
	return &SearchGamesService{api: api}
}

// FetchGames ищет игры по запросу и возвращает список найденных игр.
func (s *SearchGamesService) FetchGames(ctx context.Context, query string) ([]entities.SteamItem, error) {
	items, err := s.api.SearchGamesByName(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти игры: %w", err)
	}
	return items, nil
}

