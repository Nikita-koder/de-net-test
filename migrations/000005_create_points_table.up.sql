-- Таблица баланса поинтов пользователей
CREATE TABLE IF NOT EXISTS points (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    balance INT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);