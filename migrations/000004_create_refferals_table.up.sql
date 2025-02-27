-- Таблица реферальных связей
CREATE TABLE IF NOT EXISTS referrals (
    id TEXT PRIMARY KEY,
    referrer_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    referred_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (referrer_id, referred_id) -- Один пользователь может быть рефералом только одного пользователя
);

