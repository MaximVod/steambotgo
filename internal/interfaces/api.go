package interfaces

import (
	"context"
	"github.com/MaximVod/steambotgo/internal/entities"
)

// SteamAPI определяет методы для работы с API Steam.
type SteamAPI interface {
	// SearchGamesByName ищет игры по названию.
	// Возвращает список найденных игр (может быть пустым — не ошибка!).
	// В случае сетевой/парсинг-ошибки — возвращает error.
	SearchGamesByName(ctx context.Context, query string) ([]entities.SteamItem, error)

	// SearchGameByQuery ищет первую игру по названию.
	// Возвращает первую найденную игру или nil, если ничего не найдено.
	SearchGameByQuery(ctx context.Context, query string) (*entities.SteamItem, error)

	// GetGamePricesByCountryCode ищет цены на игру в разных странах.
	// Возвращает информацию об игре с ценами в указанной стране.
	// Если gameID указан (не 0), ищет игру с этим ID в результатах поиска.
	GetGamePricesByCountryCode(ctx context.Context, query string, countryCode string, gameID int) (*entities.SteamItem, error)
}
