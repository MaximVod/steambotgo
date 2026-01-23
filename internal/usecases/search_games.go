package usecases

import (
	"context"
	"fmt"

	"github.com/MaximVod/steambotgo/internal/entities"
	"github.com/MaximVod/steambotgo/internal/interfaces"
	"github.com/MaximVod/steambotgo/internal/logger"
)

type SearchGamesService struct {
	steamAPI interfaces.SteamAPI
	aiAPI    interfaces.AiAPI
}

func NewSearchGamesService(api interfaces.SteamAPI, aiAPI interfaces.AiAPI) *SearchGamesService {
	return &SearchGamesService{
		steamAPI: api,
		aiAPI:    aiAPI,
	}
}

// FetchGames ищет игры по запросу и возвращает список найденных игр.
func (s *SearchGamesService) FetchGames(ctx context.Context, query string) ([]entities.SteamItem, error) {
	items, err := s.steamAPI.SearchGamesByName(ctx, query)
	appLogger := logger.New()
	appLogger.Info("Result Exist " + fmt.Sprintf("SearchGamesService ", query))
	if err != nil {
		return nil, fmt.Errorf("не удалось найти игры: %w", err)
	}
	return items, nil
}

func (s *SearchGamesService) AiSearchGames(ctx context.Context, query string) error {
	_, err := s.aiAPI.SearchGamesByUserQuery(ctx, query)
	if err != nil {
		return fmt.Errorf("не удалось найти игры: %w", err)
	}
	return nil
}
