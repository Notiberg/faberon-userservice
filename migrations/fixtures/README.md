# UserService - Test Fixtures

Тестовые данные для интеграции с BookingService.

## Загрузка фикстур

```bash
# Через docker
docker exec -i userservice-db psql -U postgres -d smc_userservice < migrations/fixtures/001_test_users.sql

# Или через Makefile (если есть команда)
make fixtures
```

## Содержимое

### 001_test_users.sql

**11 пользователей, 7 автомобилей с выбранным (`is_selected = true`)**

#### Обычные пользователи (7 человек)

| tg_user_id | Имя | Car ID | Марка | Модель | Госномер | Класс |
|------------|-----|--------|-------|--------|----------|-------|
| 123456789 | Иван Петров | 1001 | BMW | X5 | А123БВ799 | L |
| 987654321 | Мария Сидорова | 2001 | Mercedes | E-Class | В999КС777 | E |
| 111222333 | Алексей Иванов | 3001 | Audi | A4 | С555АА199 | D |
| 444555666 | Екатерина Смирнова | 4001 | Tesla | Model 3 | Т123КХ777 | D |
| 555666777 | Дмитрий Волков | 5001 | Volkswagen | Polo | О777ОО799 | B |
| 666777888 | Сергей Николаев | 6001 | Porsche | Cayenne | Н123МР777 | J |
| 777888999 | Ольга Кузнецова | 7001 | Lexus | RX350 | К888КК199 | E |

#### Менеджеры компаний (3 человека)

| tg_user_id | Имя | Компания |
|------------|-----|----------|
| 777777777 | Менеджер Автомойка Премиум | Компания 1 |
| 888888888 | Менеджер СТО Профи | Компания 2 |
| 999999000 | Менеджер Детейлинг Центр | Компания 3 |

#### Пользователь без автомобиля (1 человек)

| tg_user_id | Имя | Автомобили |
|------------|-----|------------|
| 999999999 | Пользователь Без Авто | Нет |

**Используется для тестирования негативного сценария TC-2.7** (попытка создать бронирование без выбранного авто).

## Классы автомобилей

Используются для интеграции с PriceService:

| Класс | Описание | Примеры |
|-------|----------|---------|
| A | Мини | Smart, Fiat 500 |
| B | Малый | Polo, Rio |
| C | Средний | Golf, Focus |
| D | Средний+ | Camry, A4, Model 3 |
| E | Бизнес | E-Class, 5-Series, RX350 |
| F | Люкс | S-Class, 7-Series |
| J | SUV | Cayenne, X5 |
| M | MPV | Odyssey, Carnival |
| S | Спорт | 911, R8 |

## Проверка загруженных данных

```sql
-- Подключиться к БД
psql -U postgres -d smc_userservice

-- Проверить пользователей
SELECT tg_user_id, name, phone_number FROM users ORDER BY tg_user_id;

-- Проверить автомобили
SELECT c.id, c.user_id, u.name as user_name,
       c.brand, c.model, c.license_plate, c.size, c.is_selected
FROM cars c
JOIN users u ON c.user_id = u.tg_user_id
ORDER BY c.id;

-- Проверить выбранные автомобили
SELECT u.tg_user_id, u.name,
       c.brand, c.model, c.license_plate
FROM users u
JOIN cars c ON u.tg_user_id = c.user_id
WHERE c.is_selected = true
ORDER BY u.tg_user_id;

-- Проверить пользователя без автомобиля
SELECT u.tg_user_id, u.name,
       COUNT(c.id) as cars_count
FROM users u
LEFT JOIN cars c ON u.tg_user_id = c.user_id
GROUP BY u.tg_user_id, u.name
HAVING COUNT(c.id) = 0;
```

## Интеграция с BookingService

Эти данные полностью совместимы с фикстурами BookingService.

### Соответствие данных

Все `user_id` и `car_id` из BookingService фикстур присутствуют в этом файле:

- Пользователь **123456789** создаёт несколько бронирований (BMW X5)
- Пользователь **987654321** создаёт бронирования (Mercedes E-Class)
- И так далее...

### Важные моменты

- **is_selected = true**: только у одного автомобиля на пользователя
- **Уникальный индекс**: гарантирует, что у пользователя только один выбранный автомобиль
- **Пользователь 999999999**: специально создан БЕЗ автомобилей для негативных тестов

## Сброс данных

```bash
# Удалить все данные
psql -U postgres -d smc_userservice << EOF
TRUNCATE users CASCADE;
EOF

# Применить фикстуры заново
docker exec -i userservice-db psql -U postgres -d smc_userservice < migrations/fixtures/001_test_users.sql
```
