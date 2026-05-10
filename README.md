# 🏋️ Workout Tracker API

REST API для отслеживания тренировок, написанный на Go. Пет-проект для демонстрации современного стека технологий backend разработки.

## 🚀 Стек технологий

| Слой | Технология |
|---|---|
| Язык | Go 1.25 |
| Роутер | chi |
| База данных | PostgreSQL |
| Кэш | Redis |
| Авторизация | JWT |
| Документация | Swagger / OpenAPI |
| Контейнеризация | Docker + docker-compose |
| CI/CD | GitHub Actions + Render |

## 📦 Возможности

- 🔐 Регистрация и авторизация через JWT
- 💪 Управление упражнениями с полнотекстовым поиском
- 📋 Создание тренировок и добавление подходов
- 📊 Статистика прогресса с кэшированием в Redis
- 🚦 Rate limiting — защита от спама
- 📨 Фоновый воркер для еженедельных отчётов
- 🔄 Graceful shutdown

## 🗂️ Структура проекта

```
workout-tracker/
├── cmd/
│   └── api/
│       └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── cache/
│   │   └── redis.go
│   ├── handler/
│   │   ├── response/
│   │   │   └── response.go
│   │   ├── auth.go
│   │   ├── exercise.go
│   │   ├── exercise_test.go
│   │   ├── stats.go
│   │   ├── workout.go
│   │   └── workout_test.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── logger.go
│   │   └── ratelimit.go
│   ├── model/
│   │   └── model.go
│   ├── repository/
│   │   ├── mocks/
│   │   │   ├── exercise_mock.go
│   │   │   ├── set_mock.go
│   │   │   └── workout_mock.go
│   │   ├── interfaces.go
│   │   └── repository.go
│   └── worker/
│       └── report.go
├── migrations/
│   ├── 001_init.sql
│   └── 002_search.sql
├── .github/
│   └── workflows/
│       └── ci.yml
├── docker-compose.yml
├── Dockerfile
├── .env.example
└── README.md
```

## ⚡ Быстрый старт

### Локально

**1. Клонируй репозиторий**
```bash
git clone https://github.com/Alexandr20i/workout-tracker.git
cd workout-tracker
```

**2. Создай .env файл**
```bash
cp .env.example .env
```

**3. Создай базу данных**
```bash
psql -U postgres -c "CREATE DATABASE workout_tracker;"
psql -U postgres -d workout_tracker -f migrations/001_init.sql
psql -U postgres -d workout_tracker -f migrations/002_search.sql
```

**4. Запусти Redis**
```bash
docker run -d -p 6379:6379 redis:7-alpine
```

**5. Запусти сервер**
```bash
go mod tidy
go run ./cmd/api
```

Сервер запустится на `http://localhost:8080`

---

### Через Docker

```bash
docker-compose up --build
```

Поднимает сразу три сервиса: API, PostgreSQL, Redis.

---

## 📖 API документация

После запуска открой в браузере:

http://localhost:8080/swagger/

Там можно интерактивно тестировать все эндпоинты прямо из браузера.

## 🔑 Эндпоинты

### Auth
| Метод | URL | Описание |
|-------|-----|----------|
| POST | `/auth/register` | Регистрация |
| POST | `/auth/login` | Вход |

### Упражнения
| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/exercises` | Список упражнений |
| GET | `/exercises/search?q=...` | Полнотекстовый поиск |
| POST | `/exercises` | Создать упражнение |
| DELETE | `/exercises/{id}` | Удалить упражнение |

### Тренировки
| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/workouts` | История тренировок |
| POST | `/workouts` | Создать тренировку |
| GET | `/workouts/{id}` | Детали тренировки |
| DELETE | `/workouts/{id}` | Удалить тренировку |
| POST | `/workouts/{id}/sets` | Добавить подход |
| DELETE | `/workouts/{id}/sets/{setId}` | Удалить подход |

### Статистика
| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/stats/summary` | Общая статистика |
| GET | `/stats/progress?exercise_id=1` | Прогресс по упражнению |

> Все эндпоинты кроме `/auth/*` требуют заголовок `Authorization: Bearer <token>`

## 💡 Примеры запросов

**Регистрация**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret123","name":"Иван"}'
```

**Создать упражнение**
```bash
curl -X POST http://localhost:8080/exercises \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Приседания","description":"базовое упражнение","muscle_group":"Ноги"}'
```

**Создать тренировку и добавить подход**
```bash
# Создать тренировку
curl -X POST http://localhost:8080/workouts \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Тренировка ног","date":"2026-05-10"}'

# Добавить подход
curl -X POST http://localhost:8080/workouts/1/sets \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"exercise_id":1,"reps":10,"weight_kg":80}'
```

**Поиск упражнений**
```bash
curl "http://localhost:8080/exercises/search?q=приседания" \
  -H "Authorization: Bearer <TOKEN>"
```

## 🧪 Тесты

```bash
# Запустить тесты
go test ./internal/handler/... -v

# С покрытием
go test ./internal/handler/... -cover
```

## 🌍 Деплой

Проект задеплоен на **Render**:

https://workout-tracker-api-g9xg.onrender.com

> На бесплатном плане сервис засыпает после 15 минут неактивности.
> Первый запрос после сна может занять ~30 секунд.

## 🔄 CI/CD

При каждом push в `main` GitHub Actions автоматически:
1. Устанавливает зависимости
2. Запускает все тесты
3. Проверяет линтером

## 📝 Переменные окружения

| Переменная | Описание | Пример |
|---|---|---|
| `DB_HOST` | Хост PostgreSQL | `localhost` |
| `DB_PORT` | Порт PostgreSQL | `5432` |
| `DB_USER` | Пользователь БД | `postgres` |
| `DB_PASSWORD` | Пароль БД | `postgres` |
| `DB_NAME` | Имя БД | `workout_tracker` |
| `JWT_SECRET` | Секрет для JWT | `your-secret-key` |
| `JWT_EXPIRATION_HOURS` | Время жизни токена | `24` |
| `SERVER_PORT` | Порт сервера | `8080` |
| `REDIS_ADDR` | Адрес Redis | `localhost:6379` |
| `DATABASE_URL` | Полный URL БД (Render) | `postgres://...` |
| `REDIS_URL` | Полный URL Redis (Render) | `redis://...` |