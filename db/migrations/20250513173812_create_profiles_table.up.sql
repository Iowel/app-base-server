CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE,
    avatar VARCHAR(255) DEFAULT 'static/fox-icon.png',
    about TEXT,
    friends INTEGER[],
    status TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    wallet INTEGER DEFAULT 0
);