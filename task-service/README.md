# Task Service

## 📋 Overview

**Task Service** is responsible for managing tasks within the system. It is a RESTful microservice providing full CRUD (Create, Read, Update, Delete) operations for tasks.

**Key Features:**
*   Create, read, update, and delete tasks.
*   Filter task lists (by status, priority).
*   Result pagination.
*   Data validation.
*   Data isolation: tasks are linked to users (User ID).

---

## 🛠️ Technical Stack

*   **Language:** Go 1.21+
*   **Web Framework:** [Gin Gonic](https://github.com/gin-gonic/gin)
*   **Database:** PostgreSQL (Driver: `pgx` via GORM)
*   **ORM:** [GORM](https://gorm.io/)

---

## ⚙️ Configuration (.env)

```bash
# Service port
PORT=8082

# PostgreSQL Connection
DB_URL=postgres://postgres:password@postgres:5432/taskmanager?sslmode=disable
```

---

## 💾 Database

Uses the `task_schema` in PostgreSQL.

### Table `tasks`

| Field | Type | Description |
|-------|------|-------------|
| `id` | UUID | Primary Key (auto-generated) |
| `title` | VARCHAR(255) | Task title (required) |
| `description` | TEXT | Detailed description |
| `status` | VARCHAR(50) | Status (see below) |
| `priority` | VARCHAR(50) | Priority (see below) |
| `due_date` | TIMESTAMP | Due date |
| `created_by` | UUID | Creator ID (link to User Service) |
| `created_at`| TIMESTAMP | Creation date |

**Indexes** created for fields: `created_by`, `status`, `priority`, `due_date`.

---

## 🔌 API Endpoints

All requests are validated for an `Authorization` header, forwarded by the API Gateway. Although Task Service doesn't verify the JWT signature (Gateway/Auth Service does), it extracts `user_id` from the token for access filtering.

### 1. Get Task List
`GET /tasks`

**Query Parameters:**
*   `status`: Filter by status (e.g., `pending`)
*   `priority`: Filter by priority (e.g., `high`)
*   `page`: Page number (default 1)
*   `limit`: Items per page (default 10)

**Example:** `GET /tasks?status=in_progress&priority=high&page=1&limit=5`

### 2. Create Task
`POST /tasks`

**Body:**
```json
{
  "title": "Fix critical bug",
  "description": "Error in production...",
  "priority": "urgent",
  "due_date": "2024-12-31T23:59:59Z"
}
```

### 3. Get Task by ID
`GET /tasks/:id`

Returns a task only if it belongs to the current user.

### 4. Update Task
`PUT /tasks/:id`

Full update of task fields.

**Body:**
```json
{
  "title": "New Title",
  "description": "Updated desc",
  "status": "in_progress",
  "priority": "medium"
}
```

### 5. Update Status
`PATCH /tasks/:id/status`

Quick status update.

**Body:**
```json
{
  "status": "completed"
}
```

### 6. Delete Task
`DELETE /tasks/:id`

---

## 📊 Business Logic

### Statuses (`status`)
*   `pending` (Default)
*   `in_progress`
*   `completed`
*   `cancelled`

### Priorities (`priority`)
*   `low`
*   `medium` (Default)
*   `high`
*   `urgent`

---

## 🚀 Running

### Local
```bash
go mod download
go run main.go
```

### Docker
The service runs in the `task-service` container and is available within the Docker network on port `8082`.
