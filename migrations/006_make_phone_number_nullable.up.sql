-- Изменяем колонку phone_number, делая её nullable
ALTER TABLE users ALTER COLUMN phone_number DROP NOT NULL;
