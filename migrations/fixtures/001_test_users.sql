-- ==========================================
-- UserService - Тестовые фикстуры
-- ==========================================
-- Создаёт тестовых пользователей и автомобили для интеграции с BookingService
-- Соответствует данным из BookingService фикстур

-- ==========================================
-- Пользователь 1: 123456789 (BMW X5)
-- ==========================================

INSERT INTO users (tg_user_id, name, phone_number, tg_link)
VALUES (
    123456789,
    'Иван Петров',
    '+79161234567',
    '@ivan_petrov'
)
ON CONFLICT (tg_user_id) DO UPDATE SET
    name = EXCLUDED.name,
    phone_number = EXCLUDED.phone_number,
    tg_link = EXCLUDED.tg_link;

-- Автомобиль пользователя 123456789 (выбранный)
INSERT INTO cars (id, user_id, brand, model, license_plate, color, size, is_selected)
VALUES (
    1001,
    123456789,
    'BMW',
    'X5',
    'А123БВ799',
    'Черный',
    'L',
    true
)
ON CONFLICT (id) DO UPDATE SET
    brand = EXCLUDED.brand,
    model = EXCLUDED.model,
    license_plate = EXCLUDED.license_plate,
    color = EXCLUDED.color,
    size = EXCLUDED.size,
    is_selected = EXCLUDED.is_selected;

-- ==========================================
-- Пользователь 2: 987654321 (Mercedes E-Class)
-- ==========================================

INSERT INTO users (tg_user_id, name, phone_number, tg_link)
VALUES (
    987654321,
    'Мария Сидорова',
    '+79169876543',
    '@maria_sidorova'
)
ON CONFLICT (tg_user_id) DO UPDATE SET
    name = EXCLUDED.name,
    phone_number = EXCLUDED.phone_number,
    tg_link = EXCLUDED.tg_link;

-- Автомобиль пользователя 987654321 (выбранный)
INSERT INTO cars (id, user_id, brand, model, license_plate, color, size, is_selected)
VALUES (
    2001,
    987654321,
    'Mercedes',
    'E-Class',
    'В999КС777',
    'Серебристый',
    'E',
    true
)
ON CONFLICT (id) DO UPDATE SET
    brand = EXCLUDED.brand,
    model = EXCLUDED.model,
    license_plate = EXCLUDED.license_plate,
    color = EXCLUDED.color,
    size = EXCLUDED.size,
    is_selected = EXCLUDED.is_selected;

-- ==========================================
-- Пользователь 3: 111222333 (Audi A4)
-- ==========================================

INSERT INTO users (tg_user_id, name, phone_number, tg_link)
VALUES (
    111222333,
    'Алексей Иванов',
    '+79161112223',
    '@alex_ivanov'
)
ON CONFLICT (tg_user_id) DO UPDATE SET
    name = EXCLUDED.name,
    phone_number = EXCLUDED.phone_number,
    tg_link = EXCLUDED.tg_link;

-- Автомобиль пользователя 111222333 (выбранный)
INSERT INTO cars (id, user_id, brand, model, license_plate, color, size, is_selected)
VALUES (
    3001,
    111222333,
    'Audi',
    'A4',
    'С555АА199',
    'Белый',
    'D',
    true
)
ON CONFLICT (id) DO UPDATE SET
    brand = EXCLUDED.brand,
    model = EXCLUDED.model,
    license_plate = EXCLUDED.license_plate,
    color = EXCLUDED.color,
    size = EXCLUDED.size,
    is_selected = EXCLUDED.is_selected;

-- ==========================================
-- Пользователь 4: 444555666 (Tesla Model 3)
-- ==========================================

INSERT INTO users (tg_user_id, name, phone_number, tg_link)
VALUES (
    444555666,
    'Екатерина Смирнова',
    '+79164445556',
    '@kate_smirnova'
)
ON CONFLICT (tg_user_id) DO UPDATE SET
    name = EXCLUDED.name,
    phone_number = EXCLUDED.phone_number,
    tg_link = EXCLUDED.tg_link;

-- Автомобиль пользователя 444555666 (выбранный)
INSERT INTO cars (id, user_id, brand, model, license_plate, color, size, is_selected)
VALUES (
    4001,
    444555666,
    'Tesla',
    'Model 3',
    'Т123КХ777',
    'Синий',
    'D',
    true
)
ON CONFLICT (id) DO UPDATE SET
    brand = EXCLUDED.brand,
    model = EXCLUDED.model,
    license_plate = EXCLUDED.license_plate,
    color = EXCLUDED.color,
    size = EXCLUDED.size,
    is_selected = EXCLUDED.is_selected;

-- ==========================================
-- Пользователь 5: 555666777 (Volkswagen Polo)
-- ==========================================

INSERT INTO users (tg_user_id, name, phone_number, tg_link)
VALUES (
    555666777,
    'Дмитрий Волков',
    '+79165556667',
    '@dmitry_volkov'
)
ON CONFLICT (tg_user_id) DO UPDATE SET
    name = EXCLUDED.name,
    phone_number = EXCLUDED.phone_number,
    tg_link = EXCLUDED.tg_link;

-- Автомобиль пользователя 555666777 (выбранный)
INSERT INTO cars (id, user_id, brand, model, license_plate, color, size, is_selected)
VALUES (
    5001,
    555666777,
    'Volkswagen',
    'Polo',
    'О777ОО799',
    'Красный',
    'B',
    true
)
ON CONFLICT (id) DO UPDATE SET
    brand = EXCLUDED.brand,
    model = EXCLUDED.model,
    license_plate = EXCLUDED.license_plate,
    color = EXCLUDED.color,
    size = EXCLUDED.size,
    is_selected = EXCLUDED.is_selected;

-- ==========================================
-- Пользователь 6: 666777888 (Porsche Cayenne)
-- ==========================================

INSERT INTO users (tg_user_id, name, phone_number, tg_link)
VALUES (
    666777888,
    'Сергей Николаев',
    '+79166667778',
    '@sergey_nikolaev'
)
ON CONFLICT (tg_user_id) DO UPDATE SET
    name = EXCLUDED.name,
    phone_number = EXCLUDED.phone_number,
    tg_link = EXCLUDED.tg_link;

-- Автомобиль пользователя 666777888 (выбранный)
INSERT INTO cars (id, user_id, brand, model, license_plate, color, size, is_selected)
VALUES (
    6001,
    666777888,
    'Porsche',
    'Cayenne',
    'Н123МР777',
    'Черный',
    'J',
    true
)
ON CONFLICT (id) DO UPDATE SET
    brand = EXCLUDED.brand,
    model = EXCLUDED.model,
    license_plate = EXCLUDED.license_plate,
    color = EXCLUDED.color,
    size = EXCLUDED.size,
    is_selected = EXCLUDED.is_selected;

-- ==========================================
-- Пользователь 7: 777888999 (Lexus RX350)
-- ==========================================

INSERT INTO users (tg_user_id, name, phone_number, tg_link)
VALUES (
    777888999,
    'Ольга Кузнецова',
    '+79167778889',
    '@olga_kuznetsova'
)
ON CONFLICT (tg_user_id) DO UPDATE SET
    name = EXCLUDED.name,
    phone_number = EXCLUDED.phone_number,
    tg_link = EXCLUDED.tg_link;

-- Автомобиль пользователя 777888999 (выбранный)
INSERT INTO cars (id, user_id, brand, model, license_plate, color, size, is_selected)
VALUES (
    7001,
    777888999,
    'Lexus',
    'RX350',
    'К888КК199',
    'Белый',
    'E',
    true
)
ON CONFLICT (id) DO UPDATE SET
    brand = EXCLUDED.brand,
    model = EXCLUDED.model,
    license_plate = EXCLUDED.license_plate,
    color = EXCLUDED.color,
    size = EXCLUDED.size,
    is_selected = EXCLUDED.is_selected;

-- ==========================================
-- Менеджеры компаний (из SellerService)
-- ==========================================

-- Менеджер компании 1: Автомойка Премиум
INSERT INTO users (tg_user_id, name, phone_number, tg_link)
VALUES (
    777777777,
    'Менеджер Автомойка Премиум',
    '+79167777777',
    '@manager_premium'
)
ON CONFLICT (tg_user_id) DO UPDATE SET
    name = EXCLUDED.name,
    phone_number = EXCLUDED.phone_number,
    tg_link = EXCLUDED.tg_link;

-- Менеджер компании 2: СТО Профи
INSERT INTO users (tg_user_id, name, phone_number, tg_link)
VALUES (
    888888888,
    'Менеджер СТО Профи',
    '+79168888888',
    '@manager_sto'
)
ON CONFLICT (tg_user_id) DO UPDATE SET
    name = EXCLUDED.name,
    phone_number = EXCLUDED.phone_number,
    tg_link = EXCLUDED.tg_link;

-- Менеджер компании 3: Детейлинг Центр
INSERT INTO users (tg_user_id, name, phone_number, tg_link)
VALUES (
    999999000,
    'Менеджер Детейлинг Центр',
    '+79169999990',
    '@manager_detailing'
)
ON CONFLICT (tg_user_id) DO UPDATE SET
    name = EXCLUDED.name,
    phone_number = EXCLUDED.phone_number,
    tg_link = EXCLUDED.tg_link;

-- ==========================================
-- Пользователь без выбранного автомобиля
-- ==========================================
-- Для тестирования негативного сценария TC-2.7

INSERT INTO users (tg_user_id, name, phone_number, tg_link)
VALUES (
    999999999,
    'Пользователь Без Авто',
    '+79169999999',
    '@no_car_user'
)
ON CONFLICT (tg_user_id) DO UPDATE SET
    name = EXCLUDED.name,
    phone_number = EXCLUDED.phone_number,
    tg_link = EXCLUDED.tg_link;

-- У этого пользователя НЕТ автомобилей!

-- ==========================================
-- Сброс последовательностей ID (если нужно)
-- ==========================================
SELECT setval('cars_id_seq', (SELECT MAX(id) FROM cars));

-- ==========================================
-- ИТОГО: 11 пользователей, 7 автомобилей
-- ==========================================
-- 7 обычных пользователей с автомобилями
-- 3 менеджера компаний
-- 1 пользователь без автомобиля (для негативных тестов)
