-- Миграция 001: Создание таблиц для системы отслеживания цен
-- Этот файл создает структуру базы данных для хранения:
-- 1. Отслеживаемых игр (какие игры отслеживают пользователи)
-- 2. Истории цен (когда и какая была цена)

-- Таблица отслеживаемых игр
-- Хранит информацию о том, какие игры отслеживает каждый пользователь
CREATE TABLE IF NOT EXISTS tracked_games (
    -- id - уникальный идентификатор записи (автоматически увеличивается)
    id SERIAL PRIMARY KEY,
    
    -- game_id - ID игры в Steam (например, 1091500 для Cyberpunk 2077)
    game_id BIGINT NOT NULL,
    
    -- game_name - название игры (для удобства, чтобы не делать запрос к Steam API каждый раз)
    game_name VARCHAR(255) NOT NULL,
    
    -- user_chat_id - ID чата пользователя в Telegram (чтобы знать, кому отправлять уведомления)
    user_chat_id BIGINT NOT NULL,
    
    -- created_at - когда начали отслеживать игру
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    -- last_checked - когда последний раз проверяли цену этой игры
    last_checked TIMESTAMP,
    
    -- UNIQUE - один пользователь не может отслеживать одну игру дважды
    -- Это предотвращает дубликаты
    UNIQUE(game_id, user_chat_id)
);

-- Таблица истории цен
-- Хранит все снимки цен на игры (для отслеживания изменений)
CREATE TABLE IF NOT EXISTS price_snapshots (
    -- id - уникальный идентификатор записи
    id SERIAL PRIMARY KEY,
    
    -- game_id - ID игры в Steam (связь с tracked_games)
    game_id BIGINT NOT NULL,
    
    -- price - цена в центах (например, 999 = $9.99)
    -- Используем INTEGER, потому что цены в Steam API приходят в центах
    price INTEGER NOT NULL,
    
    -- currency - валюта (USD, EUR, RUB и т.д.)
    currency VARCHAR(3) NOT NULL,
    
    -- discount - процент скидки (0 если скидки нет, 50 если скидка 50%)
    discount INTEGER DEFAULT 0,
    
    -- checked_at - когда был сделан этот снимок цены
    checked_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы для производительности
-- Индексы ускоряют поиск данных в таблице

-- Индекс по user_chat_id - для быстрого поиска всех игр пользователя
CREATE INDEX IF NOT EXISTS idx_tracked_games_user ON tracked_games(user_chat_id);

-- Индекс по last_checked - для быстрого поиска игр, которые нужно проверить
CREATE INDEX IF NOT EXISTS idx_tracked_games_last_checked ON tracked_games(last_checked);

-- Индекс по game_id и checked_at - для быстрого получения истории цен игры
-- DESC означает сортировку по убыванию (новые цены первыми)
CREATE INDEX IF NOT EXISTS idx_price_snapshots_game ON price_snapshots(game_id, checked_at DESC);

