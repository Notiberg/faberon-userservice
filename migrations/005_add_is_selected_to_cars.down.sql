-- Удаляем индекс
DROP INDEX IF EXISTS idx_cars_user_selected;

-- Удаляем колонку is_selected
ALTER TABLE cars DROP COLUMN IF EXISTS is_selected;
