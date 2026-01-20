package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations применяет SQL миграции к базе данных
//
// Что такое миграции?
// Миграции - это SQL скрипты, которые создают структуру БД (таблицы, индексы).
// Мы применяем их при старте приложения, чтобы убедиться, что все таблицы созданы.
func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	// Получаем путь к файлу миграции
	// filepath.Join правильно формирует путь для разных ОС (Windows/Linux/macOS)
	migrationFile := filepath.Join("migrations", "001_create_tables.sql")

	// Читаем содержимое SQL файла
	sqlContent, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл миграции %s: %w", migrationFile, err)
	}

	// Выполняем SQL команды из файла
	// Exec выполняет SQL команды, которые не возвращают данные (CREATE, INSERT, UPDATE, DELETE)
	_, err = pool.Exec(ctx, string(sqlContent))
	if err != nil {
		return fmt.Errorf("не удалось применить миграцию: %w", err)
	}

	return nil
}
