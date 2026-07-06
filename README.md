# Courses — E-Learning Marketplace

A full-stack e-learning platform where users can **buy** and **sell** courses. Instructors can create and manage course content, while students can discover, purchase, and enroll in courses.

## Tech Stack

- **Backend:** Go (PostgreSQL, JWT auth)
- **Frontend:** Next.js (React, TypeScript, Tailwind CSS)

## Project Structure

```
├── backend/          # Go API server
│   ├── cmd/          # Entry point & migrations
│   └── internal/     # Auth, database, handlers, models
└── frontend/         # Next.js application
    └── app/          # App router pages
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
