# Frontend Application

## 💻 Overview

**Frontend** is the client side of the Task Manager application, developed with **React 19** using **TypeScript** and **Vite**. The application interacts with microservices through a unified API Gateway.

**Key Features:**
*   Registration and Authentication (JWT).
*   Protected Routes (Private Routes).
*   Dashboard with task list.
*   Create, edit, and delete tasks.
*   Filter tasks by status and priority.

---

## 🛠️ Technical Stack

*   **Build Tool:** [Vite](https://vitejs.dev/)
*   **Framework:** [React 19](https://react.dev/)
*   **Language:** [TypeScript](https://www.typescriptlang.org/)
*   **Routing:** [React Router v7](https://reactrouter.com/)
*   **Styling:** [Tailwind CSS 4](https://tailwindcss.com/)
*   **HTTP Client:** [Axios](https://axios-http.com/)

---

## 📂 Project Structure

```
src/
├── components/          # Reusable UI components
│   └── common/          # Common components (PrivateRoute, Layouts)
├── contexts/            # React Context (Global State)
│   ├── AuthContext.tsx  # Auth token and user management
│   └── TaskContext.tsx  # Task list management
├── hooks/               # Custom hooks
├── pages/               # Application pages
│   ├── Login.tsx        # Login
│   ├── Register.tsx     # Registration
│   ├── Dashboard.tsx    # Task list
│   └── TaskPage.tsx     # Task details
├── services/            # API logic
│   ├── api.ts           # Axios instance with interceptors
│   ├── authService.ts   # Auth API
│   └── taskService.ts   # Task API
├── types/               # TypeScript interfaces
├── App.tsx              # Routing and providers
└── main.tsx             # Entry point
```

---

## ⚙️ Configuration

The frontend expects the **API Gateway** to be running and accessible at `http://localhost:8000`.

API URL configuration is located in `src/services/api.ts`:
```typescript
const api = axios.create({
  baseURL: 'http://localhost:8000',
  // ...
});
```

---

## 🚀 Running and Development

### Install Dependencies
```bash
npm install
```

### Start Dev Server
```bash
npm run dev
```
The app will be available at: http://localhost:5173

### Build for Production
```bash
npm run build
```
Build output will be in the `dist/` folder.

---

## 🔐 Authentication

### Token Management
The application uses **JWT (Access Token)** for authorization.
*   On login, the token is stored in `localStorage`.
*   `AuthContext` initializes state by checking for the token.
*   Axios Interceptor (`src/services/api.ts`) automatically adds the `Authorization: Bearer <token>` header to all requests.

### Protected Routes
The `<PrivateRoute>` component checks authentication state. If the user is not logged in, they are redirected to `/login`.

---

## 🔌 API Integration

Frontend communicates only with the API Gateway (`localhost:8000`), which proxies requests to microservices:
*   `/login`, `/register` -> **Auth Service**
*   `/tasks` -> **Task Service**
*   `/profile` -> **User Service**

---

## 📝 Developing New Features

1.  **Add Page:** Create a component in `src/pages/` and add a route in `src/App.tsx`.
2.  **Add API Call:** Describe the method in the corresponding service in `src/services/`.
3.  **Add Type:** Update interfaces in `src/types/`.
