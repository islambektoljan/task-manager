# Auth Service

## 🔐 Overview

**Auth Service** is a critical microservice responsible for the security of the entire Task Manager system. It manages user registration, authentication (login), JWT token issuance and validation, and session management.

**Key Features:**
*   New user registration with password hashing.
*   Authentication and JWT (JSON Web Tokens) issuance.
*   Role-based access control (`admin`, `user`).
*   Secure logout using a token blacklist.
*   Token refreshing (Refresh Token flow).
*   Prometheus and Expvar metrics.

---

## 🛠️ Technical Stack

*   **Language:** Go 1.21+
*   **Web Framework:** [Gin Gonic](https://github.com/gin-gonic/gin)
*   **Database:** PostgreSQL (Driver: `pgx` via GORM)
*   **ORM:** [GORM](https://gorm.io/)
*   **Cache:** Redis (for Token Blacklist)
*   **Auth:** `golang-jwt/jwt/v4`
*   **Security:** `bcrypt` (for passwords)

---

## ⚙️ Configuration (.env)

Configured via environment variables:

```bash
# Service port
PORT=8081

# PostgreSQL Connection
# Format: user:password@host:port/dbname
DB_URL=postgres://postgres:password@postgres:5432/taskmanager?sslmode=disable

# Redis Connection
REDIS_URL=redis://redis:6379

# Secret key for JWT signing
# IMPORTANT: Must be long and secure in production
JWT_SECRET=your-super-secret-jwt-key-here-make-it-very-long-and-secure
```

---

## 💾 Database

Uses the `auth_schema` in PostgreSQL.

### Table `users`

| Field | Type | Description |
|-------|------|-------------|
| `id` | UUID | Primary Key (auto-generated) |
| `email` | VARCHAR(255) | Unique user email |
| `password` | VARCHAR(255) | Password hash (Bcrypt) |
| `role` | VARCHAR(50) | Role (`user` or `admin`) |
| `created_at`| TIMESTAMP | Creation date |
| `updated_at`| TIMESTAMP | Update date |

---

## 🔌 API Endpoints

### 1. Register
`POST /register`

Creates a new user.

**Request:**
```json
{
  "email": "newuser@example.com",
  "password": "strongpassword123",
  "role": "user" 
}
```
*(Role "admin" can only be created if permitted by logic; defaults to "user")*

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOi...",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "newuser@example.com",
    "role": "user"
  }
}
```

### 2. Login
`POST /login`

Authenticates user and issues a token.

**Request:**
```json
{
  "email": "newuser@example.com",
  "password": "strongpassword123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOi...",
    "user_id": "...",
    "email": "...",
    "role": "..."
  }
}
```

### 3. Logout
`POST /logout`
*Requires Header:* `Authorization: Bearer <token>`

Invalidates the current token by adding it to the "Blacklist" in Redis until its expiration.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "Successfully logged out"
  }
}
```

### 4. Refresh Token
`POST /refresh`
*Requires Header:* `Authorization: Bearer <token>`

Allows obtaining a new token if the old one is valid (or near expiry, depending on client logic). The old token is added to the Blacklist.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "token": "new.jwt.token...",
    ...
  }
}
```

### 5. Health Check
`GET /health`

Checks DB and Redis availability.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "service": "auth-service"
  }
}
```

---

## 🛡️ Security & JWT

### Token Structure (Payload)
The token contains the following Claims:
*   `user_id`: User UUID
*   `role`: User Role
*   `exp`: Expiration time (Unix timestamp)
*   `iat`: Issued at

### Blacklist (Redis)
To implement instant logout (which is impossible with stateless JWT), Redis is used.
*   On `/logout`, the token is stored in Redis.
*   Key: `blacklist:<token_string>`
*   Key TTL equals the remaining token lifetime.
*   `AuthMiddleware` checks for the token in Redis on every request.

---

## 📊 Monitoring

The service exports metrics for Prometheus and standard Go Expvar metrics.

*   `/metrics` — Prometheus metrics (requests, errors, latency).
*   `/debug/vars` — Internal Go runtime metrics (GC, goroutines, memory).

---

## 🚀 Running and Development

### Local Run (without Docker)

1.  Ensure PostgreSQL and Redis are running locally.
2.  Create an `.env` file based on the example above.
3.  Run:
    ```bash
    go mod download
    go run main.go
    ```

### Run Tests
```bash
go test ./...
```
