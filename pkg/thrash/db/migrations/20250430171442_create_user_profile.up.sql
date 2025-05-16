CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,             -- Идентификатор профиля
    user_id INT REFERENCES users(id) ON DELETE CASCADE,  -- Ссылка на пользователя (внешний ключ)
    avatar VARCHAR(255),               -- URL изображения профиля
    about TEXT,                        -- Информация о пользователе (описание)
    friends INT[],                    -- Массив ID друзей (например, список друзей может быть представлен как массив строк или целых чисел)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Время создания профиля
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP  -- Время последнего обновления
);
