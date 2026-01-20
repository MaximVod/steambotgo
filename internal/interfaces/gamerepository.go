package interfaces

import (
	"context"

	"github.com/MaximVod/steambotgo/internal/entities"
)

// GameRepository для записи игр в базу данных.
type GameRepository interface {
	// SaveTrackedGame сохраняет игру в базу данных.
	SaveTrackedGame(ctx context.Context, game *entities.TrackedGame) error
}