-- Таблица заданий
CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    points_reward INT NOT NULL CHECK (points_reward > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

