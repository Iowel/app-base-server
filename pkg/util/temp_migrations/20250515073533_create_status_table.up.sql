CREATE TABLE statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

-- Вставка статусов по умолчанию
INSERT INTO statuses (name) VALUES 
('Серебренный'),
('Золотой'),
('Бриллиантовый');
