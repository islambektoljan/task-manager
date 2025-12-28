# API Gateway (KrakenD)

## 📖 Overview

This service serves as the single entry point for the entire Task Manager microservices architecture. It is built on [KrakenD](https://www.krakend.io/) — a high-performance API Gateway.

**Key Features:**
1. **Routing:** Redirects client requests to the appropriate microservices (Auth, Task, User, Submission).
2. **Aggregation:** Capable of merging responses from multiple services (currently configured as `no-op` for direct proxying).
3. **Security (CORS):** Manages Cross-Origin Resource Sharing rules.
4. **Header Manipulation:** Forwards authorization tokens and content headers.
5. **Architecture Hiding:** The client only interacts with port `8000`, hiding internal service ports (`8081`, `8082`, etc.) from direct external access in production.

---

## ⚙️ Configuration

The configuration is located in the `krakend.json` file.

### Global Settings
- **Port:** `8000`
- **Version:** `3` (KrakenD configuration version)
- **Logging:** Enabled at `DEBUG` level with `[KRAKEND]` prefix.

### CORS (Cross-Origin Resource Sharing)
Configured to allow interaction with the Frontend application:
```json
"github_com/devopsfaith/krakend-cors": {
  "allow_origins": ["http://localhost:3000", "http://127.0.0.1:3000"],
  "allow_methods": ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"],
  "allow_headers": ["Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"],
  "allow_credentials": true,
  "max_age": "12h"
}
```

---

## 🛣️ Endpoints

All requests go to `http://localhost:8000`. KrakenD routes them according to the table below.

### 🔐 Auth Service
*Backend Host:* `http://auth-service:8081`

| Gateway Endpoint | Method | Description | Backend Endpoint | Notes |
|------------------|--------|-------------|------------------|-------|
| `/register`      | POST   | Registration | `/register`      | Direct proxy |
| `/login`         | POST   | Login        | `/login`         | Direct proxy |
| `/logout`        | POST   | Logout       | `/logout`        | Forwards `Authorization`, `Content-Type` headers |
| `/refresh`       | POST   | Refresh Token| `/refresh`       | Forwards `Authorization`, `Content-Type` headers |
| `/health`        | GET    | Health Check | `/health`        | - |

### 📋 Task Service
*Backend Host:* `http://task-service:8082`

| Gateway Endpoint | Method | Description | Backend Endpoint | Notes |
|------------------|--------|-------------|------------------|-------|
| `/tasks`         | GET    | Get task list| `/tasks`        | Forwards `Authorization` |
| `/tasks`         | POST   | Create task  | `/tasks`        | Forwards `Authorization` |
| `/tasks/{taskId}`| GET    | Get task by ID|`/tasks/{taskId}`| - |
| `/tasks/{taskId}`| PUT    | Update task  | `/tasks/{taskId}`| - |
| `/tasks/{taskId}`| DELETE | Delete task  | `/tasks/{taskId}`| - |

> **Note:** Path parameters (e.g., `{taskId}`) are passed to the backend unchanged.

### 📮 Task Submission Service
*Backend Host:* `http://task-submission-service:8083`

| Gateway Endpoint | Method | Description | Backend Endpoint |
|------------------|--------|-------------|------------------|
| `/submissions`   | POST   | Submit solution | `/submissions` |
| `/submissions/{id}`| GET  | Get submission | `/submissions/{id}` |

### 👤 User Service
*Backend Host:* `http://user-service:8084`

| Gateway Endpoint | Method | Description | Backend Endpoint |
|------------------|--------|-------------|------------------|
| `/profile`       | GET    | Get profile  | `/profile` |
| `/profile`       | PUT    | Update profile| `/profile` |

---

## 🛠️ Technical Details

### Encoding: `no-op`
All endpoints use `output_encoding: "no-op"` and `encoding: "no-op"`.
This means **KrakenD does not attempt to parse or modify the request/response body**. It acts as a transparent proxy.
*   **Pros:** High performance, support for any data format (JSON, XML, Binary).
*   **Cons:** Cannot use KrakenD data manipulation features (field filtering, renaming). This is a deliberate choice for simplicity in this project.

### Network
The service must be in the same Docker network as the microservices (`app-network`) to access them via hostnames (`auth-service`, `task-service`, etc.).

---

## 🚀 Running

The service runs as part of `docker-compose`.

```yaml
  api-gateway:
    image: krakend:latest
    ports:
      - "8000:8000"
    volumes:
      - ./api-gateway/krakend.json:/etc/krakend/krakend.json
    depends_on:
      - auth-service
      - task-service
    networks:
      - app-network
```

To check configuration:
`krakend check -c krakend.json` (requires local KrakenD installation; done automatically in Docker).

---

## 🔍 Debugging

If requests are failing:
1. Check logs: `docker-compose logs -f api-gateway`.
2. Ensure backend services are accessible from the gateway container.
3. Check CORS headers in the browser (Network tab).
