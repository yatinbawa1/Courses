# Courses — E-Learning Marketplace

A full-stack e-learning platform where users can **buy** and **sell** courses. Instructors can create and manage course content, while students can discover, purchase, and enroll in courses.

## Tech Stack

- **Backend:** Go (PostgreSQL, JWT auth)
- **Frontend:** Svelt Kit (TypeScript, Tailwind CSS)

## Project Structure

```
├── backend/                    # Go API server (net/http, pgx, Redis, JWT)
│   ├── cmd/main.go             # Entry point — wires dependencies & starts server
│   ├── internal/
│   │   ├── auth/               # Business logic: sign-up, JWT tokens, password hashing
│   │   ├── config/             # ENV loading & global config (DB, JWT, AWS, Redis, Email)
│   │   ├── database/           # PostgreSQL pool, Redis client, AWS S3 client
│   │   │   └── migrations/     # SQL migration files (init up/down)
│   │   ├── handlers/           # HTTP handlers & route registration (ServeMux)
│   │   ├── mailer/             # Resend email client + OTP HTML template
│   │   ├── models/             # Shared data structs (User, etc.)
│   │   └── repository/         # Data access: PostgreSQL (users), Redis (OTP)
│   ├── .air.toml               # Hot-reload config (Air)
│   └── Makefile                # run, migrate-up/down, delete-all-users
├── frontend/                   # SveltKit application
│   └── app/                    # App router pages
```

## Getting Started

### Backend

```bash
cd backend
cp .env.example .env   # Configure your database and auth secrets
make run
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```
