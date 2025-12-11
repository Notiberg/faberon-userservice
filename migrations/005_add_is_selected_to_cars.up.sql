-- Добавляем поле is_selected в таблицу cars
ALTER TABLE cars ADD COLUMN is_selected BOOLEAN NOT NULL DEFAULT false;

-- Создаем частичный уникальный индекс: только одна машина у пользователя может быть is_selected = true
CREATE UNIQUE INDEX idx_cars_user_selected ON cars(user_id) WHERE is_selected = true;

-- Для существующих пользователей с одним автомобилем устанавливаем is_selected = true
UPDATE cars c1
SET is_selected = true
WHERE EXISTS (
    SELECT 1
    FROM (
        SELECT user_id
        FROM cars
        GROUP BY user_id
        HAVING COUNT(*) = 1
    ) single_car_users
    WHERE single_car_users.user_id = c1.user_id
);
