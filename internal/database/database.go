package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDB создает и возвращает пул соединений с PostgreSQL
//
// connectionString - строка подключения в формате:
// postgres://username:password@host:port/database?sslmode=disable
//
// Пример:
// postgres://postgres:postgres@localhost:5432/steambotgo?sslmode=disable
func InitDB(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	// pgxpool.New создает пул соединений с БД
	// Пул - это набор готовых соединений, которые можно переиспользовать
	// Это эффективнее, чем создавать новое соединение для каждого запроса
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать пул соединений: %w", err)
	}

	// Ping проверяет, что мы действительно можем подключиться к БД
	// Если БД недоступна, вернется ошибка
	if err := pool.Ping(ctx); err != nil {
		pool.Close() // Закрываем пул, если не удалось подключиться
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	return pool, nil
}

// Close закрывает пул соединений
// Важно вызывать эту функцию при завершении работы приложения
func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}

