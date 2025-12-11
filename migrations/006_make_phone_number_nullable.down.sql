-- Откатываем изменение: делаем phone_number обязательным
-- ВНИМАНИЕ: эта миграция может упасть, если в таблице есть записи с NULL в phone_number
ALTER TABLE users ALTER COLUMN phone_number SET NOT NULL;
