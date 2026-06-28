# BudgetFlow

A full-stack personal budgeting app built with **Go Echo** (backend) and **Vue 3** (frontend).

## Features

- 🔐 JWT authentication (register / login)
- 🗂️ Budget management with color labels, categories, spend tracking
- 💳 Income & expense transactions linked to budgets
- 📊 Dashboard with balance summary, budget progress bars, monthly bar chart
- 🔍 Transaction search and category/type filtering
- 💾 SQLite database — zero infrastructure needed

---

## Project Structure

```
budgetapp/
├── backend/
│   ├── main.go              # Echo server, routes, CORS, JWT middleware
│   ├── db/db.go             # SQLite schema (users, budgets, transactions)
│   ├── models/models.go     # Data types
│   ├── handlers/handlers.go # All route handlers
│   ├── vendor/              # Vendored dependencies (no internet needed)
│   └── go.mod
└── frontend/
    ├── src/
    │   ├── api/client.js           # Axios + auth interceptor
    │   ├── stores/auth.js          # Pinia auth store
    │   ├── stores/budgets.js       # Budget CRUD store
    │   ├── stores/transactions.js  # Transaction CRUD store
    │   ├── router/index.js         # Vue Router + auth guards
    │   ├── views/AuthView.vue      # Login / register
    │   ├── views/DashboardView.vue # Summary, budgets, chart
    │   ├── views/BudgetsView.vue   # Budget CRUD
    │   ├── views/TransactionsView.vue # Transaction CRUD with filters
    │   └── components/AppShell.vue # Sidebar layout
    └── package.json
```

---

## Running Locally

### Backend

```bash
cd backend
go run main.go
# Server starts on http://localhost:8080
```

The backend uses vendored dependencies — **no `go mod tidy` or internet needed**.

SQLite database file `budget.db` is created automatically on first run.

### Frontend (dev)

```bash
cd frontend
npm install
npm run dev
# Dev server on http://localhost:5173
```

### Frontend (production)

```bash
cd frontend
npm run build
# Serve the dist/ folder with any static file server
npx serve dist
```

---

## API Reference

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/auth/register` | — | Register new user |
| POST | `/api/auth/login` | — | Login, returns JWT |
| GET | `/api/me` | ✓ | Current user info |
| GET | `/api/budgets` | ✓ | List budgets with spend totals |
| POST | `/api/budgets` | ✓ | Create budget |
| PUT | `/api/budgets/:id` | ✓ | Update budget |
| DELETE | `/api/budgets/:id` | ✓ | Delete budget |
| GET | `/api/transactions` | ✓ | List transactions (latest 100) |
| POST | `/api/transactions` | ✓ | Create transaction |
| PUT | `/api/transactions/:id` | ✓ | Update transaction |
| DELETE | `/api/transactions/:id` | ✓ | Delete transaction |
| GET | `/api/summary` | ✓ | Balance + budget stats + monthly data |

All protected routes require `Authorization: Bearer <token>` header.

---

## Tech Stack

| Layer | Tech |
|-------|------|
| Backend | Go 1.22, Echo v4, SQLite (go-sqlite3), JWT (golang-jwt/jwt v5) |
| Frontend | Vue 3, Vite, Pinia, Vue Router, Axios |
| Database | SQLite (file: `backend/budget.db`) |
