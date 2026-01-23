package interfaces

import (
	"context"
)

// AiAPI определяет методы для работы с AI.
type AiAPI interface {
	// SearchGamesByUserQuery ищет игры по запросу юзера, когда Стим не нашел его игры.
	// Возвращает список найденных игр (может быть пустым — не ошибка!).
	// В случае сетевой/парсинг-ошибки — возвращает error.
	SearchGamesByUserQuery(ctx context.Context, query string) (string, error)
}
