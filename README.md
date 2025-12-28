# Task Manager - Microservices Architecture

## 📋 Project Overview

**Task Manager** is a modern web application for task management built on microservices architecture. The project demonstrates best practices in developing distributed systems using Go, React, PostgreSQL, Redis, and KrakenD API Gateway.

### Key Features

- ✅ **Microservices architecture** with separation of concerns
- 🔐 **JWT authentication** with role support
- 🚀 **API Gateway** for request routing (KrakenD)
- 💾 **PostgreSQL** using schemas for data isolation
- ⚡ **Redis** for caching and session management
- 🐳 **Docker Compose** for easy deployment
- 🎨 **React + TypeScript** frontend with modern UI
- 📊 **Health checks** for service monitoring

---

## 🏗️ System Architecture

### Architecture Diagram

```mermaid
graph TB
    Client["Client (React Frontend)"] --> Gateway["API Gateway<br/>(KrakenD)<br/>:8000"]
    
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

### System Components

| Component | Technology | Port | Status | Description |
|-----------|-----------|------|--------|-------------|
| **Frontend** | React + TypeScript + Vite | 3000 | ✅ Implemented | User interface |
| **API Gateway** | KrakenD | 8000 | ✅ Configured | Single entry point for API |
| **Auth Service** | Go + Gin + GORM | 8081 | ✅ Implemented | Authentication and authorization |
| **Task Service** | Go + Gin + GORM | 8082 | ✅ Implemented | Task management |
| **Task Submission Service** | Go + Gin + GORM | 8083 | 🔧 Skeleton | Task submissions (in development) |
| **User Service** | Go + Gin + GORM | 8084 | 🔧 Skeleton | Profile management (in development) |
| **PostgreSQL** | PostgreSQL 13 | 5432 | ✅ Configured | Database |
| **Redis** | Redis 7 Alpine | 6379 | ✅ Configured | Cache and token management |

---

## 🔐 Auth Service

### Purpose
Handles user registration, authentication, JWT token management, and sessions.

### Technology Stack
- **Framework**: Gin (HTTP web framework)
- **ORM**: GORM (Go Object Relational Mapping)
- **Database**: PostgreSQL (`auth_schema` schema)
- **Cache**: Redis (for token blacklist)
- **Authentication**: JWT (golang-jwt/jwt)
- **Encryption**: bcrypt (for passwords)

### API Endpoints

#### Public Endpoints

| Method | Endpoint | Description | Parameters |
|--------|----------|-------------|------------|
| POST | `/register` | Register new user | `email`, `password`, `role` (optional) |
| POST | `/login` | User login | `email`, `password` |
| GET | `/health` | Service health check | - |

#### Protected Endpoints (require JWT token)

| Method | Endpoint | Description | Headers |
|--------|----------|-------------|---------|
| POST | `/logout` | Logout (blacklist token) | `Authorization: Bearer <token>` |
| POST | `/refresh` | Refresh JWT token | `Authorization: Bearer <token>` |

### Features

- ✅ JWT tokens with 24-hour expiration
- ✅ Redis-based token blacklist for secure logout
- ✅ Role-based access control (user/admin)
- ✅ Password hashing with bcrypt
- ✅ CORS configuration
- ✅ Request logging and metrics
- ✅ Graceful shutdown

---

## 📋 Task Service

### Purpose
Manages task lifecycle: create, read, update, delete (CRUD operations).

### API Endpoints

All endpoints require JWT authentication via `Authorization: Bearer <token>` header.

| Method | Endpoint | Description | Query Params |
|--------|----------|-------------|--------------|
| GET | `/tasks` | Get task list | `status`, `priority`, `page`, `limit` |
| POST | `/tasks` | Create new task | - |
| GET | `/tasks/:id` | Get task by ID | - |
| PUT | `/tasks/:id` | Update task | - |
| PATCH | `/tasks/:id/status` | Update task status only | - |
| DELETE | `/tasks/:id` | Delete task | - |
| GET | `/health` | Health check | - |

### Task Model
```go
type Task struct {
    ID          uuid.UUID  // Task UUID
    Title       string     // Title (required)
    Description string     // Description
    Status      string     // Status: pending, in_progress, completed, cancelled
    Priority    string     // Priority: low, medium, high, urgent
    DueDate     *time.Time // Due date (optional)
    CreatedBy   uuid.UUID  // Creator ID (user)
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Features

- ✅ Filtering by status and priority
- ✅ Pagination support
- ✅ User-based authorization (users see only their tasks)
- ✅ Input validation
- ✅ Database indexes for performance

---

## 🌐 API Gateway (KrakenD)

### Purpose
Single entry point for all clients. Routes requests to appropriate microservices.

### Route Mapping

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

### CORS Configuration
- Allowed origins: `http://localhost:3000`, `http://127.0.0.1:3000`
- Allowed methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
- Credentials: Enabled

---

## 💾 Database (PostgreSQL)

### Schema Architecture

Uses **schema isolation** to separate data between microservices.

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

### Schemas and Tables

#### 1. auth_schema
```sql
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
CREATE TABLE task_schema.tasks (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    priority VARCHAR(50) DEFAULT 'medium',
    due_date TIMESTAMP,
    created_by UUID NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Indexes for optimization
CREATE INDEX idx_tasks_created_by ON task_schema.tasks(created_by);
CREATE INDEX idx_tasks_status ON task_schema.tasks(status);
CREATE INDEX idx_tasks_due_date ON task_schema.tasks(due_date);
CREATE INDEX idx_tasks_priority ON task_schema.tasks(priority);
```

---

## 🎨 Frontend (React + TypeScript)

### Technology Stack
- **Framework**: React 19
- **Language**: TypeScript
- **Build Tool**: Vite
- **Routing**: React Router v7
- **HTTP Client**: Axios
- **Styling**: Tailwind CSS 4
- **State Management**: React Context API

### Project Structure
```
frontend/src/
├── components/       # Reusable components
├── contexts/         # React Context for state
│   ├── AuthContext.tsx
│   └── TaskContext.tsx
├── pages/            # Application pages
│   ├── Login.tsx
│   ├── Register.tsx
│   ├── Dashboard.tsx
│   └── TaskPage.tsx
├── services/         # API clients
│   ├── api.ts
│   ├── authService.ts
│   └── taskService.ts
└── types/            # TypeScript types
```

### Routes

| Path | Component | Protected | Description |
|------|-----------|-----------|-------------|
| `/login` | `Login` | ❌ | Login page |
| `/register` | `Register` | ❌ | Registration |
| `/dashboard` | `Dashboard` | ✅ | All tasks list |
| `/tasks/:id` | `TaskPage` | ✅ | Task details |
| `/` | Redirect | - | Redirect to `/dashboard` |

---

## 🐳 Docker & Deployment

### Quick Start

```bash
# Clone repository
git clone <repository-url>
cd task-manager

# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

### Service URLs

| Service | URL | Description |
|---------|-----|-------------|
| Frontend | `http://localhost:3000` | React application |
| API Gateway | `http://localhost:8000` | Single entry point |
| Auth Service | `http://localhost:8081` | Direct access (debug) |
| Task Service | `http://localhost:8082` | Direct access (debug) |
| PostgreSQL | `localhost:5432` | Database (user: postgres/password) |
| Redis | `localhost:6379` | Cache |

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

---

## 📝 API Usage Examples

### 1. Register User
```bash
curl -X POST http://localhost:8000/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123"
  }'
```

**Response:**
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

### 2. Login
```bash
curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123"
  }'
```

### 3. Create Task
```bash
TOKEN="your_jwt_token_here"

curl -X POST http://localhost:8000/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete project",
    "description": "Finish implementing all microservices",
    "priority": "high",
    "due_date": "2024-12-31T23:59:59Z"
  }'
```

### 4. Get Tasks
```bash
curl -X GET "http://localhost:8000/tasks?status=pending&page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

### 5. Update Task Status
```bash
TASK_ID="123e4567-e89b-12d3-a456-426614174000"

curl -X PATCH http://localhost:8000/tasks/$TASK_ID/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```

---

## 🔍 Monitoring and Debugging

### Docker Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f auth-service
docker-compose logs -f task-service
```

### Database Access
```bash
# Via Docker
docker-compose exec postgres psql -U postgres -d taskmanager

# SQL queries
SELECT * FROM auth_schema.users;
SELECT * FROM task_schema.tasks;
```

### Redis Access
```bash
# Via Docker
docker-compose exec redis redis-cli

# Redis commands
KEYS *                    # All keys
KEYS blacklist:*          # All blacklisted tokens
GET blacklist:<token>     # Check specific token
TTL blacklist:<token>     # Token TTL
```

---

## 🛠️ Development Status

### ✅ Implemented

- [x] Microservices architecture
- [x] API Gateway (KrakenD)
- [x] Auth Service (complete)
- [x] Task Service (complete)
- [x] PostgreSQL with migrations
- [x] Redis integration
- [x] Docker Compose configuration
- [x] Frontend (React + TypeScript)

### 🔧 In Development

- [ ] User Service (profile management)
- [ ] Task Submission Service
- [ ] Frontend UI/UX improvements
- [ ] Unit and integration tests

### 💡 Planned

- [ ] Kubernetes deployment
- [ ] CI/CD pipeline
- [ ] Prometheus + Grafana monitoring
- [ ] WebSocket for real-time notifications
- [ ] OAuth2 integration
- [ ] API documentation (Swagger/OpenAPI)

---

## 📚 Technologies

### Backend
- **Gin** - HTTP web framework
- **GORM** - ORM for database operations
- **JWT** - Authentication tokens
- **Redis** - Caching and session storage
- **bcrypt** - Password hashing

### Frontend
- **React 19** - UI framework
- **TypeScript** - Type safety
- **Vite** - Build tool
- **React Router** - Navigation
- **Axios** - HTTP client
- **Tailwind CSS** - Styling

### Infrastructure
- **PostgreSQL 13** - Database
- **Redis 7** - Cache
- **KrakenD** - API Gateway
- **Docker** - Containerization
- **Docker Compose** - Orchestration

---

## 🎯 Conclusion

**Task Manager** is a complete example of modern microservices architecture using current technologies. The project demonstrates:

- ✅ Proper separation of concerns between services
- ✅ Secure authentication with JWT and Redis
- ✅ API Gateway for centralized routing
- ✅ Schema isolation in PostgreSQL
- ✅ Containerization with Docker
- ✅ Modern React frontend with TypeScript

The project is ready for further development and new features! 🚀
