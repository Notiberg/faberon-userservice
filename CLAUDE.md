# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SMC-UserService - сервис управления пользователями и их автомобилями для приложения автомойки. Использует Clean Architecture с разделением на domain, service, infrastructure и handlers.

### Tech Stack
- **Language**: Go 1.24+
- **Architecture**: Clean Architecture (Domain, Service, Repository, Handlers)
- **Database**: PostgreSQL 16 + sqlx + golang-migrate
- **HTTP Router**: Gorilla Mux
- **Query Builder**: Squirrel (psqlbuilder wrapper)
- **Authentication**: Simplified (X-User-ID + X-User-Role headers for MVP)
- **Monitoring**: Prometheus + Grafana
- **Logging**: Custom logger (console + file)
- **Containerization**: Docker Compose
- **Build Tool**: Makefile

## Development Commands

### Quick Start (Docker)
```bash
# Запуск всех сервисов (рекомендуется)
make docker-up

# Просмотр логов приложения
make docker-logs-app

# Остановка всех сервисов
make docker-down
```

### Local Development
```bash
# Запуск только инфраструктуры
make dev

# Запуск приложения локально
make run

# Сборка бинарного файла
make build
```

### Testing
```bash
make test
```

### Database Management
```bash
# Применить миграции
make migrate-up

# Откатить миграции
make migrate-down

# Полный сброс БД
make db-reset

# Загрузить тестовые фикстуры (11 пользователей, 7 автомобилей)
make fixtures
```

### Cleanup Commands
```bash
# Очистка артефактов и логов
make clean

# Остановка сервисов и удаление volumes
make docker-clean

# Удаление Docker образов проекта
make docker-prune

# Полная очистка (artifacts + logs + volumes + images)
make clean-all
```

### Database
- PostgreSQL работает на порту **5435** (не стандартный 5432)
- Connection string: `host=localhost port=5435 user=postgres password=postgres dbname=smc_userservice sslmode=disable`
- Миграции автоматически применяются при `docker-compose up`

### Monitoring
- **Prometheus**: http://localhost:9091 - сбор метрик
- **Grafana**: http://localhost:3001 - визуализация метрик
  - Логин: `admin`
  - Пароль: `admin`
  - Dashboard "SMC UserService - HTTP Metrics" автоматически загружается при старте

## Architecture

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`)
   - `user.go` - доменная модель User с db тегами для sqlx
     - Поля: TGUserID, Name, PhoneNumber, TGLink, RoleID, Role, CreatedAt
     - PhoneNumber, TGLink и все nullable поля используют указатели
   - `car.go` - доменная модель Car с db тегами для sqlx
     - Car.ID использует int64 (BIGSERIAL в БД)
     - IsSelected (bool) - флаг выбранного автомобиля
   - `role.go` - ролевая модель и проверки прав доступа
     - Роли: client, manager, superuser
     - Константы RoleID: RoleIDClient (1), RoleIDManager (2), RoleIDSuperUser (3)
     - Методы: CanAccessUser(), CanModifyUser()

2. **Service Layer** (`internal/service/user/`)
   - `user_service.go` - полная бизнес-логика для User и Car
     - User: CreateUser, UpdateUser, DeleteUser, GetUserByID, GetUserWithCars, GetSuperUsers
     - Car: CreateCar, UpdateCar (PATCH + role), DeleteCar (role), GetSelectedCar, SetSelectedCar
     - Логика выбранного автомобиля:
       - Первый созданный автомобиль автоматически становится выбранным
       - При удалении выбранного автомобиля первый из оставшихся становится выбранным
       - Только один автомобиль может быть выбран одновременно
     - Все методы для Car принимают параметр role для проверки прав доступа
     - Проверка прав доступа на основе ролей (role.CanModifyUser)
     - Обработка ошибок с wrapping через fmt.Errorf
   - `contract.go` - интерфейсы репозиториев и кастомные ошибки
     - ErrUserNotFound, ErrUserAlreadyExists, ErrCarNotFound, ErrCarAccessDenied
   - `models/models.go` - все DTO модели (User, Car, UserWithCars)
     - DTOs включают поле role
     - PhoneNumber является необязательным полем (*string) во всех DTO

3. **Infrastructure Layer** (`internal/infra/storage/`)
   - `user/repository.go` - UserRepository с обработкой ошибок
     - JOIN с таблицей roles для получения имени роли
     - Create/GetByTGID/Update/Delete с поддержкой role_id
     - GetSuperUsers - получение списка tg_user_id всех суперпользователей
   - `car/repository.go` - CarRepository с обработкой ошибок
     - GetSelectedByUserID - получение выбранного автомобиля
     - UnselectAllByUserID - снятие выбора со всех автомобилей пользователя
   - Используют psqlbuilder для построения SQL запросов
   - Все методы имеют собственные кастомные ошибки (ErrCreateUser, ErrGetUser, etc.)

4. **Handlers Layer** (`internal/handlers/`)
   - `api/helpers.go` - утилиты (RespondJSON, RespondError, DecodeJSON и специфичные ошибки)
   - `api/create_user/handler.go` - POST /users
   - `api/get_current_user/handler.go` - GET /users/me (возвращает UserWithCars с is_selected)
   - `api/update_current_user/handler.go` - PUT /users/me
   - `api/delete_current_user/handler.go` - DELETE /users/me
   - `api/create_car/handler.go` - POST /users/me/cars (первый автомобиль автоматически выбирается)
   - `api/update_car/handler.go` - PATCH /users/me/cars/{car_id} (извлекает role из context и передаёт в сервис)
   - `api/delete_car/handler.go` - DELETE /users/me/cars/{car_id} (извлекает role из context и передаёт в сервис)
   - `api/select_car/handler.go` - PUT /users/me/cars/{car_id}/select (установка автомобиля как выбранного)
   - `api/get_user_by_id/handler.go` - GET /internal/users/{tg_user_id} (межсервисное взаимодействие)
   - `api/get_selected_car/handler.go` - GET /internal/users/{tg_user_id}/cars/selected (межсервисное взаимодействие, получение выбранного автомобиля по user_id)
   - `api/get_superusers/handler.go` - GET /internal/users/superusers (межсервисное взаимодействие, получение списка всех суперпользователей)
   - `middleware/auth.go` - упрощённая аутентификация через X-User-ID и X-User-Role
     - Функции: UserIDAuth, GetUserIDFromContext, GetRoleFromContext, RequireSuperUser
   - `middleware/metrics.go` - Prometheus метрики middleware

5. **Configuration** (`internal/config/`)
   - `config.go` - загрузка config.toml с секциями logs, server, database

### Key Design Patterns

- **Repository Pattern**: Абстракция доступа к данным через интерфейсы
- **DTO Pattern**: Разделение доменных моделей и API контрактов
- **Dependency Injection**: Все зависимости через конструкторы
- **Error Wrapping**: Многоуровневая обработка ошибок (repository → service → handler)
- **Handler per Endpoint**: Каждый endpoint в отдельной папке для тестирования
- **Middleware Pattern**: Metrics и Auth middleware для всех endpoints

### Logging

Сервис использует кастомный логгер (`pkg/logger`) с многоуровневой записью:
- **INFO** - информационные сообщения (только консоль)
- **WARN** - предупреждения для 4xx ошибок (консоль + файл `logs/app.log`)
- **ERROR** - ошибки для 5xx ошибок (консоль + файл `logs/app.log`)

Все handlers логируют:
- Успешные операции (INFO): метод, endpoint, user_id, результат
- Ошибки валидации (WARN): метод, endpoint, user_id, детали ошибки
- Ошибки авторизации (WARN): метод, endpoint, причина
- Внутренние ошибки (ERROR): метод, endpoint, user_id, полный стек ошибки

### Monitoring & Observability

**Prometheus метрики:**
- **http_requests_total** - счётчик всех HTTP запросов (labels: method, endpoint, status)
- **http_request_duration_seconds** - гистограмма длительности запросов (labels: method, endpoint, status)
- **http_requests_in_flight** - gauge текущих активных запросов

Метрики доступны на endpoint `/metrics`

**Grafana Dashboard:**

Dashboard "SMC UserService - HTTP Metrics" включает:
1. **HTTP Request Rate** - частота запросов в секунду по endpoint'ам
2. **Requests In Flight** - текущее количество обрабатываемых запросов
3. **HTTP Request Duration (p95, p99)** - перцентили времени обработки запросов
4. **HTTP Status Codes Distribution** - распределение статус-кодов
5. **Requests by Endpoint** - запросы по endpoint'ам

**Как посмотреть метрики:**
1. Запустите сервис: `go run cmd/main.go`
2. Запустите мониторинг: `docker-compose up -d prometheus grafana`
3. Сделайте несколько запросов к API для генерации метрик
4. Откройте Grafana: http://localhost:3001 (admin/admin)
5. Dashboard автоматически загружен и доступен на главной странице

### Database Schema

- **roles table**:
  - id SERIAL PRIMARY KEY
  - name VARCHAR(20) UNIQUE (client, manager, superuser)
  - description TEXT
  - created_at TIMESTAMP

- **users table**:
  - tg_user_id BIGINT PRIMARY KEY
  - name (required), phone_number (nullable), tg_link (nullable), created_at
  - role_id INTEGER FK → roles(id) ON DELETE RESTRICT
  - Index на phone_number
  - Index на role_id

- **cars table**:
  - id BIGSERIAL PRIMARY KEY (автоинкремент)
  - user_id BIGINT FK → users(tg_user_id) ON DELETE CASCADE
  - brand, model, license_plate (required)
  - color, size (nullable)
    - size - класс автомобиля согласно европейской системе классов: A (мини), B (малые), C (средние/гольф-класс), D (большие средние), E (бизнес-класс), F (люксовые), J (внедорожники), M (минивэны), S (спорткары)
  - is_selected BOOLEAN NOT NULL DEFAULT false - флаг выбранного автомобиля
  - Indexes на user_id и license_plate
  - Уникальный индекс idx_cars_user_selected (user_id) WHERE is_selected = true - гарантирует только один выбранный автомобиль у пользователя

### Important Utilities

- **psqlbuilder** (`pkg/psqlbuilder/psqlbuilder.go`): обертки над squirrel
  - Автоматически использует Dollar placeholder ($1, $2) для PostgreSQL
  - Select(), Insert(), Update(), Delete()

## API Contract

API реализует OpenAPI спецификацию из `schemas/api/schema.yaml`:

### Public Endpoints
- `POST /users` - создание пользователя (с указанием роли, phone_number опционально)

### Internal Endpoints (для межсервисного взаимодействия)
- `GET /internal/users/superusers` - получение списка tg_user_id всех суперпользователей
- `GET /internal/users/{tg_user_id}` - получение пользователя с автомобилями по ID
- `GET /internal/users/{tg_user_id}/cars/selected` - получение текущего выбранного автомобиля пользователя по его ID

### Protected Endpoints (требуют X-User-ID и X-User-Role)
- `GET /users/me` - получение пользователя с автомобилями (включает is_selected для каждого автомобиля)
- `PUT /users/me` - частичное обновление профиля (обновляются только переданные поля)
- `DELETE /users/me` - удаление профиля
- `POST /users/me/cars` - добавление автомобиля (первый автомобиль автоматически становится выбранным)
- `PATCH /users/me/cars/{car_id}` - частичное обновление автомобиля
- `DELETE /users/me/cars/{car_id}` - удаление автомобиля (при удалении выбранного, первый из оставшихся становится выбранным)
- `PUT /users/me/cars/{car_id}/select` - установка автомобиля как выбранного (предыдущий автоматически снимается с выбора)

### Monitoring Endpoints
- `GET /metrics` - Prometheus метрики в формате OpenMetrics

### Authentication (MVP - Упрощённая версия)

Для доступа к защищённым endpoints используются два заголовка:
```
X-User-ID: <telegram_user_id>
X-User-Role: <client|manager|superuser>
```

**⚠️ Важно**: Это временное решение для MVP. В продакшене планируется:
- Отдельный SMC-AuthService для генерации JWT токенов
- Валидация Telegram InitData
- Refresh token механизм
- Полноценная JWT аутентификация

### Role-Based Access Control

**Client** (role_id=1):
- Может видеть и изменять **только свои данные**
- Может управлять **только своими автомобилями**
- Проверка: `targetUserID == requestUserID`
- При попытке доступа к чужим данным получает 403 Forbidden

**Manager** (role_id=2):
- То же, что и Client для данных пользователей
- Дополнительный доступ к настройкам компании (в других сервисах)

**Superuser** (role_id=3):
- **Полный доступ** ко всем данным
- Может изменять любых пользователей и автомобили
- `role.CanModifyUser()` всегда возвращает true
- Обходит все проверки на владение ресурсами

### Примеры проверки прав доступа

```go
// В user_service.go для операций с автомобилями:
func (s *Service) UpdateCar(ctx context.Context, tgID int64, carID int64, input models.UpdateCarInputDTO, role domain.Role) (*models.CarDTO, error) {
    car, err := s.carRepo.GetByID(ctx, carID)

    // Проверка прав: superuser может изменять любую машину, остальные - только свои
    if !role.CanModifyUser(car.UserID, tgID) {
        return nil, ErrCarAccessDenied
    }
    // ... обновление автомобиля
}
```

```go
// В domain/role.go:
func (r Role) CanModifyUser(targetUserID, requestUserID int64) bool {
    switch r {
    case RoleSuperUser:
        return true // Superuser может изменять любого
    case RoleManager, RoleClient:
        return targetUserID == requestUserID // Только свои данные
    default:
        return false
    }
}
```

## Configuration

`config.toml`:
```toml
[logs]
level = "info"

[server]
http_port = 8080

[database]
host = "localhost"
port = 5435  # Не стандартный порт!
user = "postgres"
password = "postgres"
dbname = "smc_userservice"
sslmode = "disable"
```

### Environment Variables Override

При запуске в Docker конфигурация БД переопределяется переменными окружения из `docker-compose.yml`:
- `DATABASE_HOST` → host (для Docker используется "postgres")
- `DATABASE_PORT` → port (для Docker используется 5432)
- `DATABASE_USER` → user
- `DATABASE_PASSWORD` → password
- `DATABASE_NAME` → dbname

При локальном запуске (`make run`) используются значения из `config.toml`.

## Project Structure

```
SMC-UserService/
├── cmd/
│   └── main.go                 # Entry point с routing и DI
├── internal/
│   ├── config/
│   │   └── config.go          # Config loader
│   ├── domain/
│   │   ├── user.go
│   │   └── car.go
│   ├── service/user/
│   │   ├── contract.go        # Interfaces и errors
│   │   ├── user_service.go    # Business logic
│   │   └── models/models.go   # DTOs
│   ├── infra/storage/
│   │   ├── user/repository.go
│   │   └── car/repository.go
│   └── handlers/
│       ├── middleware/
│       │   ├── auth.go
│       │   └── metrics.go
│       └── api/
│           ├── helpers.go
│           ├── create_user/handler.go
│           ├── get_current_user/handler.go
│           ├── update_current_user/handler.go
│           ├── delete_current_user/handler.go
│           ├── create_car/handler.go
│           ├── update_car/handler.go
│           └── delete_car/handler.go
├── pkg/psqlbuilder/
├── migrations/
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_cars_table.up.sql
│   ├── 002_create_cars_table.down.sql
│   ├── 003_add_roles.up.sql
│   └── 003_add_roles.down.sql
├── docker-compose.yml
├── config.toml
└── .gitignore
```