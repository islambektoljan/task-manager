# Task Submission Service

## 📮 Overview

**Task Submission Service** is responsible for the task submission process. It allows users to submit work results (links, files, comments) for specific tasks.

> ⚠️ **Status:** Currently in early development stage (Skeleton).

---

## 🛠️ Technical Stack

*   **Language:** Go 1.21+
*   **Web Framework:** [Gin Gonic](https://github.com/gin-gonic/gin)
*   **Database:** PostgreSQL (schema `submission_schema`)
*   **ORM:** [GORM](https://gorm.io/)

---

## ⚙️ Configuration (.env)

```bash
PORT=8083
DB_URL=postgres://postgres:password@postgres:5432/taskmanager?sslmode=disable
```

---

## 💾 Database

The service plans to use the `submissions` table in `submission_schema`.

| Field | Type | Description |
|-------|------|-------------|
| `id` | UUID | Submission ID |
| `task_id` | UUID | Task ID (linked to Task Service) |
| `user_id` | UUID | User ID (author) |
| `solution` | TEXT | Solution text or link |
| `status` | VARCHAR | Status (`submitted`, `reviewed`) |
| `submitted_at` | TIMESTAMP | Submission date |
| `score` | INT | Score (0-100) |
| `comments` | TEXT | Reviewer comments |

---

## 🔌 API Endpoints (Planned)

### 1. Submit Solution
`POST /submissions`
*Requires Header:* `Authorization: Bearer <token>`

**Request:**
```json
{
  "task_id": "uuid-task-id",
  "solution": "https://github.com/my-solution"
}
```

**Response:**
```json
{
  "id": "uuid-submission-id",
  "status": "submitted",
  "submitted_at": "2023-10-27T10:00:00Z"
}
```

---

## 🚀 Development

The service is connected to the shared `app-network` and ready for business logic implementation.

```bash
# Run (skeleton)
go run main.go
```
