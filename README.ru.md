# Task Manager - Микросервисная Архитектура

## 📋 Обзор Проекта

**Task Manager** - это современное веб-приложение для управления задачами, построенное на микросервисной архитектуре. Проект демонстрирует лучшие практики разработки распределенных систем с использованием Go, React, PostgreSQL, Redis и KrakenD API Gateway.

### Ключевые Особенности

- ✅ **Микросервисная архитектура** с разделением ответственности
- 🔐 **JWT-аутентификация** с поддержкой ролей
- 🚀 **API Gateway** для маршрутизации запросов (KrakenD)
- 💾 **PostgreSQL** с использованием схем для изоляции данных
- ⚡ **Redis** для кэширования и управления сессиями
- 🐳 **Docker Compose** для простого развертывания
- 🎨 **React + TypeScript** frontend с современным UI
- 📊 **Health checks** для мониторинга состояния сервисов

---

## 🏗️ Архитектура Системы

### Архитектурная Диаграмма

```mermaid
graph TB
    Client["Клиент (React Frontend)"] --> Gateway["API Gateway<br/>(KrakenD)<br/>:8000"]
    
    Gateway --> Auth["Auth Service<br/>(Go)<br/>:8081"]
    Gateway --> Task["Task Service<br/>(Go)<br/>:8082"]
    Gateway --> Submission["Task Submission Service<br/>(Go)<br/>:8083"]
    Gateway --> User["User Service<br/>(Go)<br/>:8084"]
    
    Auth --> PostgreSQL["PostgreSQL<br/>:5432"]
    Auth --> Redis["Redis<br/>:6379"]
    Task --> PostgreSQL
    Submission --> PostgreSQL
    User --> PostgreSQL
    
    PostgreSQL -.-> AuthSchema["auth_schema"]
    PostgreSQL -.-> TaskSchema["task_schema"]
    PostgreSQL -.-> SubmissionSchema["submission_schema"]
    PostgreSQL -.-> UserSchema["user_schema"]
```

### Компоненты Системы

| Компонент | Технология | Порт | Статус | Описание |
|-----------|-----------|------|--------|----------|
| **Frontend** | React + TypeScript + Vite | 3000 | ✅ Реализован | Пользовательский интерфейс |
| **API Gateway** | KrakenD | 8000 | ✅ Настроен | Единая точка входа для API |
| **Auth Service** | Go + Gin + GORM | 8081 | ✅ Реализован | Аутентификация и авторизация |
| **Task Service** | Go + Gin + GORM | 8082 | ✅ Реализован | Управление задачами |
| **Task Submission Service** | Go + Gin + GORM | 8083 | 🔧 Заглушка | Отправка задач (в разработке) |
| **User Service** | Go + Gin + GORM | 8084 | 🔧 Заглушка | Управление профилями (в разработке) |
| **PostgreSQL** | PostgreSQL 13 | 5432 | ✅ Настроен | База данных |
| **Redis** | Redis 7 Alpine | 6379 | ✅ Настроен | Кэш и управление токенами |

---

## 🔐 Auth Service (Сервис Аутентификации)

### Назначение
Отвечает за регистрацию пользователей, аутентификацию, управление JWT токенами и сессиями.

### Технологический Стек
- **Framework**: Gin (HTTP web framework)
- **ORM**: GORM (Go Object Relational Mapping)
- **БД**: PostgreSQL (схема `auth_schema`)
- **Кэш**: Redis (для blacklist токенов)
- **Аутентификация**: JWT (golang-jwt/jwt)
- **Шифрование**: bcrypt (для паролей)

### API Endpoints

#### Публичные Endpoints

| Method | Endpoint | Описание | Параметры |
|--------|----------|----------|-----------|
| POST | `/register` | Регистрация нового пользователя | `email`, `password`, `role` (опционально) |
| POST | `/login` | Вход в систему | `email`, `password` |
| GET | `/health` | Health check статус сервиса | - |

#### Защищенные Endpoints (требуют JWT токен)

| Method | Endpoint | Описание | Headers |
|--------|----------|----------|---------|
| POST | `/logout` | Выход из системы (blacklist токена) | `Authorization: Bearer <token>` |
| POST | `/refresh` | Обновление JWT токена | `Authorization: Bearer <token>` |

### Структура Данных

#### User Model
```go
type User struct {
    ID        uuid.UUID  // UUID пользователя
    Email     string     // Уникальный email
    Password  string     // Хешированный пароль (bcrypt)
    Role      string     // Роль: 'user' или 'admin'
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Схема БД
```sql
-- Таблица: auth_schema.users
CREATE TABLE auth_schema.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Особенности Реализации

#### 1. JWT Токены
- **Время жизни**: 24 часа
- **Алгоритм**: HS256
- **Payload**: `user_id`, `role`, `exp`, `iat`

#### 2. Управление Сессиями (Redis)
- **Blacklist токенов**: При логауте токен добавляется в Redis с TTL = оставшееся время жизни
- **Проверка**: Middleware проверяет каждый токен на наличие в blacklist
- **Ключи**: `blacklist:<token>`

#### 3. Безопасность
- ✅ CORS настроен для `localhost:3000`
- ✅ Security headers (middleware)
- ✅ Request logging
- ✅ Password hashing (bcrypt, cost 10)
- ✅ Email validation
- ✅ Graceful shutdown

### Мониторинг
- **Metrics**: `/metrics` (Prometheus-совместимые метрики)
- **Expvar**: `/debug/vars` (внутренние метрики Go)
- **Health**: `/health` (проверка DB и Redis)

---

## 📋 Task Service (Сервис Задач)

### Назначение
Управление жизненным циклом задач: создание, чтение, обновление, удаление (CRUD операции).

### Технологический Стек
- **Framework**: Gin
- **ORM**: GORM
- **БД**: PostgreSQL (схема `task_schema`)

### API Endpoints

Все endpoints требуют JWT аутентификацию через header `Authorization: Bearer <token>`.

| Method | Endpoint | Описание | Query Params |
|--------|----------|----------|--------------|
| GET | `/tasks` | Получить список задач | `status`, `priority`, `page`, `limit` |
| POST | `/tasks` | Создать новую задачу | - |
| GET | `/tasks/:id` | Получить задачу по ID | - |
| PUT | `/tasks/:id` | Обновить задачу | - |
| PATCH | `/tasks/:id/status` | Обновить только статус задачи | - |
| DELETE | `/tasks/:id` | Удалить задачу | - |
| GET | `/health` | Health check | - |

### Структура Данных

#### Task Model
```go
type Task struct {
    ID          uuid.UUID  // UUID задачи
    Title       string     // Заголовок (обязательно)
    Description string     // Описание
    Status      string     // Статус: pending, in_progress, completed, cancelled
    Priority    string     // Приоритет: low, medium, high, urgent
    DueDate     *time.Time // Срок выполнения (опционально)
    CreatedBy   uuid.UUID  // ID создателя (пользователя)
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Схема БД
```sql
-- Таблица: task_schema.tasks
CREATE TABLE task_schema.tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    priority VARCHAR(50) DEFAULT 'medium',
    due_date TIMESTAMP,
    created_by UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для оптимизации
CREATE INDEX idx_tasks_created_by ON task_schema.tasks(created_by);
CREATE INDEX idx_tasks_status ON task_schema.tasks(status);
CREATE INDEX idx_tasks_due_date ON task_schema.tasks(due_date);
CREATE INDEX idx_tasks_priority ON task_schema.tasks(priority);
```

### Бизнес-Логика

#### 1. Фильтрация и Пагинация
```http
GET /tasks?status=pending&priority=high&page=1&limit=20
```
- Поддержка фильтрации по `status` и `priority`
- Пагинация через `page` и `limit`
- По умолчанию: page=1, limit=10

#### 2. Авторизация
- Пользователи видят только **свои задачи** (по `created_by`)
- JWT токен декодируется в middleware для получения `user_id`

#### 3. Валидация Статусов
Допустимые значения:
- `pending` - В ожидании
- `in_progress` - В процессе
- `completed` - Завершено
- `cancelled` - Отменено

#### 4. Валидация Приоритетов
Допустимые значения:
- `low` - Низкий
- `medium` - Средний (по умолчанию)
- `high` - Высокий
- `urgent` - Срочный

---

## 🌐 API Gateway (KrakenD)

### Назначение
Единая точка входа для всех клиентов. Маршрутизация запросов к соответствующим микросервисам.

### Конфигурация

#### CORS Settings
```json
{
  "allow_origins": ["http://localhost:3000", "http://127.0.0.1:3000"],
  "allow_methods": ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"],
  "allow_headers": ["Origin", "Content-Type", "Authorization", "Accept"],
  "allow_credentials": true,
  "max_age": "12h"
}
```

#### Маршрутизация Endpoints

```mermaid
graph LR
    Client["Client :3000"] -->|/register| Gateway["Gateway :8000"]
    Client -->|/login| Gateway
    Client -->|/tasks| Gateway
    Client -->|/profile| Gateway
    
    Gateway -->|/register| Auth["Auth :8081"]
    Gateway -->|/login| Auth
    Gateway -->|/logout| Auth
    Gateway -->|/refresh| Auth
    
    Gateway -->|/tasks| Task["Task :8082"]
    Gateway -->|/tasks/:id| Task
    
    Gateway -->|/submissions| Submission["Submission :8083"]
    
    Gateway -->|/profile| User["User :8084"]
```

### Основные Routes

| Frontend Path | Backend Service | Backend Path | Auth Required |
|--------------|-----------------|--------------|---------------|
| `/register` | `auth-service:8081` | `/register` | ❌ |
| `/login` | `auth-service:8081` | `/login` | ❌ |
| `/logout` | `auth-service:8081` | `/logout` | ✅ |
| `/refresh` | `auth-service:8081` | `/refresh` | ✅ |
| `/tasks` | `task-service:8082` | `/tasks` | ✅ |
| `/tasks/:id` | `task-service:8082` | `/tasks/:id` | ✅ |
| `/submissions` | `task-submission-service:8083` | `/submissions` | ❌ |
| `/profile` | `user-service:8084` | `/profile` | ❌ |

### Особенности
- **No-op encoding**: Используется для auth endpoints, чтобы передавать ответы как есть
- **Header forwarding**: Authorization headers автоматически передаются в backend
- **Health checks**: Доступны через `/health`

---

## 💾 База Данных (PostgreSQL)

### Архитектура БД

Используется **схемная изоляция** для разделения данных разных микросервисов.

```mermaid
graph TB
    PostgreSQL["PostgreSQL Database: taskmanager"]
    
    PostgreSQL --> AuthSchema["auth_schema"]
    PostgreSQL --> TaskSchema["task_schema"]
    PostgreSQL --> SubmissionSchema["submission_schema"]
    PostgreSQL --> UserSchema["user_schema"]
    
    AuthSchema --> UsersTable["users"]
    TaskSchema --> TasksTable["tasks"]
```

### Схемы и Таблицы

#### 1. auth_schema
```sql
-- Пользователи
CREATE TABLE auth_schema.users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

#### 2. task_schema
```sql
-- Задачи
CREATE TABLE task_schema.tasks (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    priority VARCHAR(50) DEFAULT 'medium',
    due_date TIMESTAMP,
    created_by UUID NOT NULL,  -- Ссылка на auth_schema.users.id
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

#### 3. submission_schema
*(В разработке)*

#### 4. user_schema
*(В разработке)*

### Миграции

Миграции выполняются автоматически при запуске каждого сервиса через GORM AutoMigrate или SQL-файлы в папках `migrations/`.

**Структура миграций:**
```
migrations/
├── 001_create_users.up.sql      # Создание таблицы users
├── 001_create_users.down.sql    # Откат миграции
├── 002_add_role_to_users.up.sql # Добавление поля role
└── 002_add_role_to_users.down.sql
```

---

## 🎨 Frontend (React + TypeScript)

### Технологический Стек
- **Framework**: React 19
- **Language**: TypeScript
- **Build Tool**: Vite
- **Routing**: React Router v7
- **HTTP Client**: Axios
- **Styling**: Tailwind CSS 4
- **State Management**: React Context API

### Структура Проекта
```
frontend/src/
├── components/          # Переиспользуемые компоненты
│   └── common/
│       └── PrivateRoute.tsx  # Защищенные маршруты
├── contexts/           # React Context для состояния
│   ├── AuthContext.tsx      # Контекст аутентификации
│   └── TaskContext.tsx      # Контекст задач
├── hooks/              # Кастомные React hooks
├── pages/              # Страницы приложения
│   ├── Login.tsx            # Страница входа
│   ├── Register.tsx         # Страница регистрации
│   ├── Dashboard.tsx        # Главная страница с задачами
│   └── TaskPage.tsx         # Страница отдельной задачи
├── services/           # API клиенты
│   ├── api.ts              # Базовая конфигурация Axios
│   ├── authService.ts      # Сервис аутентификации
│   └── taskService.ts      # Сервис задач
├── types/              # TypeScript типы
├── App.tsx             # Главный компонент
└── main.tsx            # Точка входа
```

### Маршруты (Routes)

| Path | Component | Защищен | Описание |
|------|-----------|---------|----------|
| `/login` | `Login` | ❌ | Вход в систему |
| `/register` | `Register` | ❌ | Регистрация |
| `/dashboard` | `Dashboard` | ✅ | Список всех задач |
| `/tasks/:id` | `TaskPage` | ✅ | Детали конкретной задачи |
| `/` | Redirect | - | Перенаправление на `/dashboard` |

### Аутентификация (Frontend)

#### AuthContext
Управляет состоянием аутентификации:
- Хранение JWT токена в `localStorage`
- Автоматическое добавление токена в HTTP headers
- Logout функционал
- Проверка авторизации

#### PrivateRoute
Компонент-обертка для защищенных маршрутов:
```tsx
<PrivateRoute>
  <Dashboard />
</PrivateRoute>
```

### API Integration

#### Axios Instance (api.ts)
```typescript
const api = axios.create({
  baseURL: 'http://localhost:8000',  // API Gateway
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor для автоматического добавления токена
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

---

## 🐳 Docker & Deployment

### Docker Compose Architecture

```yaml
services:
  # Базы данных и кэш
  - postgres:13           # PostgreSQL БД
  - redis:7-alpine        # Redis Cache
  
  # API Gateway
  - api-gateway           # KrakenD :8000
  
  # Микросервисы
  - auth-service          # Go :8081
  - task-service          # Go :8082
  - task-submission-service  # Go :8083
  - user-service          # Go :8084
```

### Environment Variables

#### Auth Service
```bash
DB_URL=postgres://postgres:password@postgres:5432/taskmanager?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-here-make-it-very-long-and-secure
REDIS_URL=redis://redis:6379
PORT=8081
```

#### Task Service
```bash
DB_URL=postgres://postgres:password@postgres:5432/taskmanager?sslmode=disable
PORT=8082
```

### Health Checks

Все сервисы поддерживают health checks:

```yaml
# PostgreSQL
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U postgres -d taskmanager"]
  interval: 5s
  timeout: 5s
  retries: 5

# Redis
healthcheck:
  test: ["CMD", "redis-cli", "ping"]
  interval: 5s
  timeout: 3s
  retries: 5
```

### Docker Network

Все сервисы находятся в одной сети `app-network` (bridge driver), что позволяет им взаимодействовать по именам сервисов.

---

## 🚀 Установка и Запуск

### Предварительные Требования

- Docker 20.10+
- Docker Compose 2.0+
- Node.js 18+ (для локальной разработки frontend)
- Go 1.21+ (для локальной разработки backend)

### Быстрый Старт (Docker)

#### 1. Клонировать репозиторий
```bash
git clone <repository-url>
cd task-manager
```

#### 2. Настроить переменные окружения
```bash
# Опционально: создать .env файлы в каждом сервисе
cp auth-service/.env.example auth-service/.env
cp task-service/.env.example task-service/.env
```

#### 3. Запустить все сервисы
```bash
docker-compose up -d
```

#### 4. Проверить статус
```bash
docker-compose ps
```

Все сервисы должны быть в статусе `Up (healthy)`.

#### 5. Проверить health endpoints
```bash
# Auth Service
curl http://localhost:8000/health

# Task Service  
curl http://localhost:8082/health
```

### Доступ к Сервисам

| Сервис | URL | Описание |
|--------|-----|----------|
| Frontend | `http://localhost:3000` | React приложение |
| API Gateway | `http://localhost:8000` | Единая точка входа |
| Auth Service | `http://localhost:8081` | Прямой доступ (для отладки) |
| Task Service | `http://localhost:8082` | Прямой доступ (для отладки) |
| PostgreSQL | `localhost:5432` | БД (логин: postgres/password) |
| Redis | `localhost:6379` | Кэш |

### Локальная Разработка

#### Backend (Go)
```bash
# Auth Service
cd auth-service
go mod download
go run main.go

# Task Service
cd task-service
go mod download
go run main.go
```

#### Frontend (React)
```bash
cd frontend
npm install
npm run dev
```

Frontend будет доступен на `http://localhost:5173` (Vite dev server).

---

## 📝 Примеры Использования API

### 1. Регистрация пользователя
```bash
curl -X POST http://localhost:8000/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123"
  }'
```

**Ответ:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "role": "user"
  }
}
```

### 2. Вход в систему
```bash
curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123"
  }'
```

### 3. Создание задачи
```bash
TOKEN="your_jwt_token_here"

curl -X POST http://localhost:8000/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Завершить проект",
    "description": "Закончить реализацию всех микросервисов",
    "priority": "high",
    "due_date": "2024-12-31T23:59:59Z"
  }'
```

### 4. Получение списка задач
```bash
curl -X GET "http://localhost:8000/tasks?status=pending&page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

### 5. Обновление статуса задачи
```bash
TASK_ID="123e4567-e89b-12d3-a456-426614174000"

curl -X PATCH http://localhost:8000/tasks/$TASK_ID/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```

### 6. Выход из системы (Logout)
```bash
curl -X POST http://localhost:8000/logout \
  -H "Authorization: Bearer $TOKEN"
```

---

## 🔍 Мониторинг и Отладка

### Логи Docker
```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f auth-service
docker-compose logs -f task-service
```

### Метрики

#### Auth Service Metrics
```bash
# Prometheus метрики
curl http://localhost:8081/metrics

# Go expvar метрики
curl http://localhost:8081/debug/vars
```

### Подключение к БД
```bash
# Через Docker
docker-compose exec postgres psql -U postgres -d taskmanager

# Локально (если PostgreSQL установлен)
psql -h localhost -U postgres -d taskmanager
```

#### SQL Queries
```sql
-- Посмотреть всех пользователей
SELECT * FROM auth_schema.users;

-- Посмотреть все задачи
SELECT * FROM task_schema.tasks;

-- Посмотреть задачи с пользователями
SELECT 
    t.id, t.title, t.status, t.priority,
    u.email as created_by_email
FROM task_schema.tasks t
JOIN auth_schema.users u ON t.created_by = u.id;
```

### Подключение к Redis
```bash
# Через Docker
docker-compose exec redis redis-cli

# Команды Redis
KEYS *                    # Все ключи
KEYS blacklist:*          # Все токены в blacklist
GET blacklist:<token>     # Проверить конкретный токен
TTL blacklist:<token>     # Время жизни токена
```

---

## 🛠️ Статус Разработки

### ✅ Реализовано

- [x] Микросервисная архитектура
- [x] API Gateway (KrakenD)
- [x] Auth Service (полностью)
  - [x] Регистрация и логин
  - [x] JWT токены
  - [x] Redis blacklist
  - [x] Роли пользователей
  - [x] Health checks
- [x] Task Service (полностью)
  - [x] CRUD операции
  - [x] Фильтрация и пагинация
  - [x] Авторизация по пользователю
  - [x] Валидация данных
- [x] PostgreSQL с миграциями
- [x] Redis интеграция
- [x] Docker Compose конфигурация
- [x] Frontend (React + TypeScript)
  - [x] Аутентификация
  - [x] Роутинг
  - [x] Context API
  - [x] Axios интеграция

### 🔧 В Разработке

- [ ] User Service
  - [ ] Профиль пользователя
  - [ ] Обновление данных
  - [ ] Статистика пользователя
- [ ] Task Submission Service
  - [ ] Отправка решений задач
  - [ ] Проверка решений
  - [ ] История отправок
- [ ] Frontend улучшения
  - [ ] UI/UX дизайн
  - [ ] Формы создания/редактирования задач
  - [ ] Уведомления (Toast)
  - [ ] Drag-and-drop для задач
- [ ] Тестирование
  - [ ] Unit тесты (Go)
  - [ ] Integration тесты
  - [ ] E2E тесты (Frontend)

### 💡 Планируется

- [ ] Kubernetes deployment
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Prometheus + Grafana мониторинг
- [ ] ELK Stack для логирования
- [ ] Message queue (RabbitMQ/Kafka)
- [ ] WebSocket для real-time уведомлений
- [ ] File upload для задач
- [ ] Email notifications
- [ ] OAuth2 интеграция (Google, GitHub)
- [ ] Rate limiting
- [ ] API documentation (Swagger/OpenAPI)

---

## 📚 Технологии и Библиотеки

### Backend (Go)

| Библиотека | Версия | Назначение |
|------------|--------|------------|
| gin-gonic/gin | v1.9+ | HTTP web framework |
| gorm.io/gorm | v1.25+ | ORM для работы с БД |
| golang-jwt/jwt | v4 | JWT токены |
| redis/go-redis | v9 | Redis клиент |
| google/uuid | latest | UUID генерация |
| bcrypt | latest | Хеширование паролей |
| godotenv | latest | Загрузка .env файлов |

### Frontend (React)

| Библиотека | Версия | Назначение |
|------------|--------|------------|
| react | ^19.1 | UI framework |
| react-router-dom | ^7.9 | Роутинг |
| axios | ^1.13 | HTTP клиент |
| typescript | ~5.9 | Типизация |
| vite | ^7.1 | Build tool |
| tailwindcss | ^4.1 | CSS framework |

### Infrastructure

| Технология | Версия | Назначение |
|------------|--------|------------|
| PostgreSQL | 13 | База данных |
| Redis | 7 Alpine | Кэш и сессии |
| KrakenD | latest | API Gateway |
| Docker | 20.10+ | Контейнеризация |
| Docker Compose | 2.0+ | Оркестрация |

---

## 🤝 Contributing

Проект открыт для улучшений! Если вы хотите внести свой вклад:

1. Fork репозиторий
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

---

## 📄 Лицензия

Этот проект создан в образовательных целях.

---

## 📞 Контакты и Поддержка

Если у вас возникли вопросы или проблемы, пожалуйста:
- Создайте Issue в репозитории
- Проверьте логи Docker: `docker-compose logs -f`
- Убедитесь, что все сервисы запущены: `docker-compose ps`

---

## 🎯 Заключение

**Task Manager** - это полноценный пример современной микросервисной архитектуры с использованием актуальных технологий. Проект демонстрирует:

- ✅ Правильное разделение ответственности между сервисами
- ✅ Безопасную аутентификацию с JWT и Redis
- ✅ Использование API Gateway для централизованной маршрутизации
- ✅ Схемную изоляцию данных в PostgreSQL
- ✅ Контейнеризацию с Docker
- ✅ Современный React frontend с TypeScript

Проект готов для дальнейшего развития и добавления новых функций! 🚀
