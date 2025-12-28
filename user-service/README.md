# User Service

## 👤 Overview

**User Service** is responsible for managing user profiles.
Unlike Auth Service, which handles security (login/password), User Service stores extended user information (name, avatar, settings).

> ⚠️ **Status:** Currently in early development stage (Skeleton).

---

## 🛠️ Technical Stack

*   **Language:** Go 1.21+
*   **Web Framework:** [Gin Gonic](https://github.com/gin-gonic/gin)
*   **Database:** PostgreSQL (schema `user_schema`)
*   **ORM:** [GORM](https://gorm.io/)

---

## ⚙️ Configuration (.env)

```bash
PORT=8084
DB_URL=postgres://postgres:password@postgres:5432/taskmanager?sslmode=disable
```

---

## 🔌 API Endpoints (Planned)

### 1. Get Profile
`GET /profile`
*Requires Header:* `Authorization: Bearer <token>`

Returns information about the current user.

**Response:**
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "first_name": "John",
  "last_name": "Doe",
  "bio": "Software Engineer",
  "avatar": "https://..."
}
```

### 2. Update Profile
`PUT /profile`
*Requires Header:* `Authorization: Bearer <token>`

**Request:**
```json
{
  "first_name": "John",
  "last_name": "Doe Updated",
  "bio": "New bio info"
}
```

---

## 🚀 Development

The service is connected to the shared `app-network` and database, but the main business logic is yet to be implemented.

```bash
# Run (skeleton)
go run main.go
```
