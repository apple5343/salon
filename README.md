# Auto Salon API
Учебный проект 1 курса КФУ.

Backend-система управления автосалоном, написанная на Go. Проект реализует REST API для
работы с брендами, моделями автомобилей, клиентами, сотрудниками, поставщиками, продажами и
событиями, а также модуль аналитики и генератор тестовых данных (симуляцию).

## Стек

- **Язык:** Go 1.24
- **HTTP-фреймворк:** [Echo](https://echo.labstack.com/)
- **DI-контейнер:** [uber-go/fx](https://github.com/uber-go/fx)
- **База данных:** PostgreSQL 14
- **Кэш:** Redis 7
- **Миграции:** [golang-migrate](https://github.com/golang-migrate/migrate)
- **Аутентификация:** JWT (access + refresh токены)
- **Логирование:** zap
- **Контейнеризация:** Docker + docker-compose

## Архитектура

Проект построен по слоистой архитектуре (transport → service → repository → infrastructure):

```
cmd/                      Точки входа
internal/
  app/                    Сборка приложения (fx-модули)
  config/                 Конфигурация из env (через cleanenv)
  models/                 Доменные модели и DTO
  infrastructure/         Клиенты postgres и redis
  repository/             Доступ к данным (postgres + redis-кэш)
  service/                Бизнес-логика
  transport/
    http/                 REST-хендлеры (Echo) и мидлвари
    cli/                  Консольные команды
  pkg/
    simulation/           Генератор тестовых данных
  utils/                  Утилиты 
pkg/                      Переиспользуемые пакеты
migrations/               SQL-миграции (up/down)
config/                   Примеры env-файлов
datasets/                 Датасеты для симуляции 
tests/integration/        Интеграционные тесты 
```

## Функциональность

- **Бренды / Модели / Автомобили** — CRUD с кэшированием брендов и поставщиков в Redis.
- **Клиенты** — регистрация (сотрудником), профиль, обновление.
- **Сотрудники** — регистрация, авторизация (JWT), найм (`hire`), профиль, обновление.
- **Поставщики** — CRUD + удаление.
- **Продажи** — создание, список, получение по id, обновление.
- **События** — журнал событий (поступление машин и т.п.).
- **Аналитика** (защищено авторизацией):
  - `/analytics/sales` — продажи за период (`date_from`, `date_to`)
  - `/analytics/warehouse` — состояние склада
  - `/analytics/supply` — поставки
  - `/analytics/employee/:id` — метрики по сотруднику
  - `/analytics/supplier/:id` — метрики по поставщику
- **Симуляция** — заполнение БД реалистичными данными из датасетов (настраивается через env).

## Эндпоинты

| Метод | Путь | Auth |
|-------|------|------|
| GET | `/health` | — |
| POST | `/employees/auth/login` | — |
| POST | `/employees/auth/register` | admin |
| GET | `/employees/auth/access` | refresh token |
| GET | `/employees/auth/refresh` | refresh token |
| GET | `/employees/me` | user |
| POST | `/employees/:id/hire` | admin |
| GET | `/employees` | admin |
| GET | `/employees/:id` | admin |
| PUT | `/employees/:id` | admin |
| POST | `/clients/auth/register` | admin |
| GET | `/clients/:id` | admin |
| PUT | `/clients/:id` | admin |
| POST | `/brands` | admin |
| GET | `/brands` | — (soft) |
| GET | `/brands/:id` | — (soft) |
| PUT | `/brands/:id` | admin |
| POST | `/models` | admin |
| GET | `/models` | — |
| GET | `/models/:id` | — (soft) |
| PUT | `/models/:id` | admin |
| POST | `/cars` | admin |
| GET | `/cars` | — (soft) |
| GET | `/cars/:id` | — (soft) |
| PUT | `/cars/:id` | admin |
| POST | `/suppliers` | admin |
| GET | `/suppliers` | — (soft) |
| GET | `/suppliers/:id` | — (soft) |
| PUT | `/suppliers/:id` | admin |
| DELETE | `/suppliers/:id` | admin |
| POST | `/sales` | admin |
| GET | `/sales` | admin |
| GET | `/sales/:id` | admin |
| PATCH | `/sales/:id` | admin |
| GET | `/events` | admin |
| GET | `/events/:id` | admin |
| GET | `/analytics/sales` | admin |
| GET | `/analytics/warehouse` | admin |
| GET | `/analytics/supply` | admin |
| GET | `/analytics/employee/:id` | admin |
| GET | `/analytics/supplier/:id` | admin |

> `soft` — запрос проходит, даже если токен отсутствует/невалиден (публичное чтение).

## Быстрый старт

### Docker Compose

```bash
# Запустите Docker Compose
docker compose --env-file ./config/example.env up -d
```

Сервер будет доступен на `http://localhost:8080` (health-check — `GET /health`).

Будет создан админимтратор с логином `admin@gmail.com` и паролем `qwerty1234` (по умолчанию).

### Симуляция данных

```bash
# Заполните БД тестовыми данными
ACCESS_SECRET=example_access_secret REFRESH_SECRET=example_refresh_secret go run cmd/simulation/main.go
```

Основные параметры симуляции (env):

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `SIMULATION_START_DATE` | `2022-01-01` | начальная дата |
| `SIMULATION_DAYS_COUNT` | `365` | число дней симуляции |

Параметры, отвечающие за наполнение в первый день симуляции:

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `SIMULATION_BRANDS_PERCENT` | `50` | процент создания от общего числа (`datasets/brands.json`) |
| `SIMULATION_MODELS_PERCENT` | `50` | процент создания от общего числа (`datasets/models.json`) |
| `SIMULATION_SUPPLIERS_PERCENT` | `50` | процент создания от общего числа (`datasets/suppliers.json`) |
| `ADMINS_COUNT` | `3` | число администраторов |
| `EMPLOYEES_COUNT` | `5` | число сотрудников |

Параметры, отвечащие за наполнение в остальные дни симуляции:

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `SIMULATION_NEW_EMLOYEES_RATIO` | `0.01` | вероятность создания нового сотрудника |
| `SIMULATION_NEW_CARS_RATIO` | `0.01` | вероятность создания нового автомобиля |
| `SIMULATION_NEW_BRANDS_RATIO` | `0.01` | вероятность создания нового бренда |
| `SIMULATION_NEW_MODELS_RATIO` | `0.01` | вероятность создания нового модели |
| `SIMULATION_NEW_SUPPLIERS_RATIO` | `0.01` | вероятность создания нового поставщика |
| `SIMULATION_NEW_SALES_RATIO` | `0.01` | вероятность создания нового заказа |
## Переменные окружения

Обязательные (без них приложение не стартует):

- `ACCESS_SECRET` — секрет для подписи access-токенов
- `REFRESH_SECRET` — секрет для подписи refresh-токенов

Остальные имеют значения по умолчанию (см. `.env.example` и пакет `internal/config`).

## Тестирование

