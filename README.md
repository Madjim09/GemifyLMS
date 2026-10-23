# GemifyLMS

Сервис для геймификации обучения персонала компании: опыт (XP) и уровни, достижения, дерево навыков, квизы и рейтинг сотрудников.

## Возможности

- Регистрация и вход (JWT)
- XP и уровни, начисление за события обучения
- Достижения и дерево навыков
- Квизы с прохождением и историей результатов
- Таблица лидеров
- Приём событий обучения из внешних систем по webhook (формат xAPI)
- Демо-данные загружаются автоматически при первом запуске

## Технологии

| Часть | Стек |
|-------|------|
| Backend | Go 1.21, Gin, PostgreSQL 16, Redis 7 (кеш), JWT |
| Frontend | Next.js 14, React 18, TypeScript, Tailwind CSS |
| Запуск | Docker, Docker Compose |

Бэкенд построен по принципам Clean Architecture: `domain` → `usecase` → `repository` / `delivery`.

## Быстрый старт (Docker)

Самый простой способ: всё поднимается одной командой.

**Нужно установить:** [Docker](https://docs.docker.com/get-docker/) с Docker Compose v2.

```bash
git clone https://github.com/Madjim09/GemifyLMS.git
cd GemifyLMS
docker compose up --build
```

Первая сборка занимает несколько минут. Когда в логах появится `Server starting on port 8080`, сервис готов.

| Сервис | Адрес |
|--------|-------|
| Веб-интерфейс | http://localhost:3000 |
| API | http://localhost:8080 |
| Проверка API | http://localhost:8080/health |
| PostgreSQL | localhost:5432 |
| Redis | localhost:6379 |

Остановить: `Ctrl+C`, затем `docker compose down`.
Остановить и **удалить данные БД**: `docker compose down -v`.

## Демо-аккаунты

Откройте http://localhost:3000/auth. Пароль у всех: `demo123`

| Email | Роль | Уровень |
|-------|------|---------|
| novice@demo.com | Новичок | 1 |
| active@demo.com | Активный сотрудник | 4 |
| mid@demo.com | Средний уровень | 5 |
| top@demo.com | Топ-сотрудник | 10 |

Либо зарегистрируйте нового пользователя на той же странице.

## Запуск для разработки (без Docker для приложения)

PostgreSQL и Redis поднимаем через Docker, бэкенд и фронтенд запускаем локально.

**Нужно установить:** Docker, Go 1.21+, Node.js 20+.

### 1. База данных и Redis

```bash
docker compose up -d postgres redis
```

### 2. Backend

```bash
cd backend
cp .env.example .env
go run ./cmd/server
```

Windows (PowerShell): `copy .env.example .env`

При старте сервер сам применяет миграции и загружает демо-данные. API будет доступен на http://localhost:8080.

### 3. Frontend

В новом терминале:

```bash
cd frontend
npm install
npm run dev
```

Откройте http://localhost:3000.

## Переменные окружения

### Backend (`backend/.env`)

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `PORT` | `8080` | Порт API |
| `DB_HOST` | `localhost` | Хост PostgreSQL |
| `DB_PORT` | `5432` | Порт PostgreSQL |
| `DB_USER` | `postgres` | Пользователь БД |
| `DB_PASSWORD` | `postgres` | Пароль БД |
| `DB_NAME` | `gamification` | Имя БД |
| `DB_SSLMODE` | `disable` | Режим SSL |
| `REDIS_ADDR` | `localhost:6379` | Адрес Redis |
| `REDIS_PASSWORD` | пусто | Пароль Redis |
| `JWT_SECRET` | встроенное значение | Секрет подписи JWT |

Redis необязателен: если он недоступен, сервис работает без кеширования.

### Frontend

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080` | Адрес backend API |

> Значение `NEXT_PUBLIC_API_URL` вшивается в фронтенд на этапе сборки. Если меняете адрес API, пересоберите фронтенд (`docker compose build frontend`).

## API

Базовый путь: `/api/v1`

| Метод | Путь | Доступ | Описание |
|-------|------|--------|----------|
| POST | `/auth/register` | публичный | Регистрация |
| POST | `/auth/login` | публичный | Вход, возвращает JWT |
| GET | `/leaderboard` | публичный | Таблица лидеров |
| GET | `/achievements` | публичный | Все достижения |
| GET | `/skills` | публичный | Дерево навыков |
| GET | `/quizzes` | публичный | Список квизов |
| GET | `/quizzes/:id` | публичный | Квиз с вопросами |
| GET | `/users` | публичный | Список пользователей |
| GET | `/users/:id` | публичный | Профиль пользователя |
| GET | `/users/:id/achievements` | публичный | Достижения пользователя |
| GET | `/users/:id/skills` | публичный | Навыки пользователя |
| GET | `/users/:id/quiz-results` | публичный | Результаты квизов |
| POST | `/quizzes/:id/submit` | JWT | Отправить ответы квиза |

Дополнительно: `GET /health` (проверка состояния) и `POST /webhook/event` (приём xAPI-событий).

Для защищённых запросов передавайте заголовок `Authorization: Bearer <token>`.

Пример входа:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"top@demo.com","password":"demo123"}'
```

## Тесты

```bash
cd backend
go test ./...
```

## Структура проекта

```
GemifyLMS/
├── backend/
│   ├── cmd/server/          # точка входа
│   └── internal/
│       ├── domain/          # сущности
│       ├── usecase/         # бизнес-логика
│       ├── repository/      # работа с PostgreSQL
│       ├── delivery/http/   # хендлеры и роутер (Gin)
│       └── infrastructure/  # подключение к БД, миграции, seed, Redis
├── frontend/
│   └── src/
│       ├── app/             # страницы Next.js
│       ├── components/
│       ├── contexts/        # AuthContext
│       └── lib/api.ts       # клиент API
└── docker-compose.yml
```

## Частые проблемы

**Порт уже занят (`port is already allocated`).** Освободите порт 3000, 8080, 5432 или 6379 либо измените левую часть в `ports` в `docker-compose.yml`, например `"5433:5432"`.

**Backend не стартует: `Failed to connect to database`.** PostgreSQL ещё не поднялся или недоступен. При локальном запуске проверьте, что контейнер работает (`docker compose ps`) и значения в `backend/.env` верны.

**Фронтенд не видит API / ошибки CORS или Network Error.** Убедитесь, что backend отвечает на http://localhost:8080/health и что `NEXT_PUBLIC_API_URL` указывает на правильный адрес.

**Нужно сбросить данные.** `docker compose down -v` удалит том с БД; при следующем запуске демо-данные создадутся заново.

## Продакшен

Перед публикацией обязательно:

- задайте свой `JWT_SECRET` (длинная случайная строка);
- смените пароль PostgreSQL и не публикуйте порты БД и Redis наружу;
- уберите автоматическую загрузку демо-аккаунтов с известным паролем;
- ограничьте `AllowOrigins` в CORS (сейчас разрешены все источники).