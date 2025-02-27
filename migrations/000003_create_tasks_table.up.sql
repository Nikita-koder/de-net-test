-- Таблица выполненных заданий пользователями
CREATE TABLE IF NOT EXISTS user_tasks (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, task_id) -- Один пользователь может выполнить задание только один раз
);

