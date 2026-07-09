# Courses — E-Learning Marketplace

A full-stack e-learning platform where users can **buy** and **sell** courses. Instructors can create and manage course content, while students can discover, purchase, and enroll in courses. The frontend is built into the Go binary — a single, deployable server serves both the API and the SPA.

## Tech Stack

### Backend
| Component | Technology |
|---|---|
| Language | Go 1.26.4 |
| HTTP Router | `net/http` standard library (Go 1.22+ pattern routing) |
| Database | PostgreSQL via `pgx/v5` (connection pool) |
| Cache / OTP Store | Redis via `go-redis/v9` |
| Auth | JWT (`golang-jwt/jwt/v5`) — access (1h) + refresh (15d) tokens |
| Password Hashing | bcrypt (`golang.org/x/crypto`) |
| Object Storage | AWS S3 SDK v2 (`ap-south-1`) |
| Email | Resend API (`resend-go/v3`) |
| UUID | `google/uuid` |

### Frontend
| Component | Technology |
|---|---|
| Framework | Svelte 5 + SvelteKit 2 |
| Rendering | Client-side only (CSR SPA, `ssr = false`) |
| Language | TypeScript |
| CSS | Tailwind CSS 4 + `tw-animate-css` |
| UI Library | shadcn-svelte ("maia" style, olive base) |
| Icons | Lucide (`lucide-svelte`) |
| Package Manager | pnpm |
| Build | Vite + `@sveltejs/adapter-static` |

### Infrastructure
- **PostgreSQL** — Docker container (`courses-psql`)
- **Redis** — Docker Container (`Courses-redis`)
- **Resend** — transactional email service

## Project Structure

```
├── backend/                          # Go API server (net/http, pgx, Redis, JWT, S3, Resend)
│   ├── cmd/main.go                   # Entry point — wires dependencies & starts server on :8080
│   ├── assets.go                     # //go:embed dist/* — embeds SvelteKit build into binary
│   ├── dist/                         # SvelteKit build output (gitignored, rebuilt on deploy)
│   ├── Makefile                      # run, migrate-*, delete-all-users
│   ├── .air.toml                     # Hot-reload config (Air)
│   ├── internal/
│   │   ├── config/                   # Environment variable loading (godotenv)
│   │   │   └── config.go            # Exports package-level vars: DB, JWT, AWS, Redis, Email, Port
│   │   ├── database/                 # Database & external service clients
│   │   │   ├── connection.go        # ConnectDataBase() — creates pgxpool.Pool, pings & returns
│   │   │   ├── redis.go            # NewRedisClient() — creates & pings redis.Client
│   │   │   ├── aws.go              # GetS3Client() — creates s3.Client from AWS config
│   │   │   └── migrations/          # SQL migration files
│   │   │       ├── 000001_init.up.sql     # Creates 8 tables (User, Course, etc.) + Content_Type enum
│   │   │       ├── 000001_init.down.sql   # Drops all tables & enum
│   │   │       └── delete_all_user.sql    # TRUNCATE "User" CASCADE (dev utility)
│   │   ├── models/                   # Shared data structures
│   │   │   └── user_model.go        # User struct (User_id, Username, Email, HashedPassword, etc.)
│   │   ├── repository/               # Data access layer (PostgreSQL + Redis)
│   │   │   ├── user.go              # UserRepo — CheckIfEmailExists, GetPasswordForEmail, Add (with uniqueness check)
│   │   │   └── otp.go              # RedisOTPRepo — SaveOTP (with TTL), VerifyOTP (with auto-delete)
│   │   ├── auth/                     # Authentication business logic
│   │   │   ├── auth_service.go      # AuthService struct + UserDataRepo / OTPRepo interfaces
│   │   │   ├── sign_up.go           # SignUpUsingEmailAndPassword — password validation, bcrypt hash, DB insert
│   │   │   ├── login.go             # LoginWithEmailPassword — bcrypt compare, JWT token generation
│   │   │   └── tokens.go            # CreateRefreshToken (15d), CreateAccessToken (1h), VerifyAccessToken
│   │   ├── middleware/               # HTTP middleware
│   │   │   └── auth.go             # CheckAuth — extracts Bearer token, verifies JWT, enriches context
│   │   ├── handlers/                 # HTTP handlers & route registration
│   │   │   ├── routes.go            # RegisterRoutes — wires all routes into ServeMux
│   │   │   ├── home.go              # Health-check handler ("Hello! Server is Running")
│   │   │   └── auth/                 # Authentication endpoint handlers
│   │   │       ├── send_otp.go       # POST /api/auth/send-otp — email check, OTP generation, Redis save, Resend email
│   │   │       ├── verify_otp.go     # POST /api/auth/send-otp/verify — OTP verify, sign-up, response
│   │   │       ├── login.go          # POST /api/auth/login — password verify, set JWT cookies
│   │   │       └── me.go            # GET /api/auth/me — stub (returns current user profile)
│   │   └── mailer/                   # Email sending
│   │       ├── init.go              # MailSender interface + ResendMailer (SendOTPMail via Resend API)
│   │       ├── templates.go         # Embed & parse OTP HTML template (text/template)
│   │       └── templates/otp.html   # Styled OTP email HTML template
│   └── tmp/                          # Air hot-reload build artifacts
├── frontend/                         # SvelteKit SPA (CSR, embedded into Go binary)
│   ├── package.json                 # Dependencies: svelte, sveltekit, tailwindcss, shadcn-svelte, lucide-svelte
│   ├── pnpm-workspace.yaml          # pnpm workspace config
│   ├── vite.config.ts               # Vite config with adapter-static → backend/dist/
│   ├── tsconfig.json                # TypeScript configuration
│   ├── prettier.config.js           # Code formatting config
│   ├── components.json              # shadcn-svelte component configuration
│   ├── src/
│   │   ├── app.html                 # SPA HTML shell (sveltekit %sveltekit.head% / %sveltekit.body%)
│   │   ├── app.css                  # Global styles + Tailwind imports
│   │   ├── app.d.ts                 # App-level type declarations
│   │   ├── routes/
│   │   │   ├── +layout.svelte       # Root layout (wraps all pages)
│   │   │   ├── +layout.ts           # Layout config: ssr = false, prerender = false
│   │   │   ├── +page.svelte         # Home page
│   │   │   ├── layout.css           # Layout-specific styles
│   │   │   └── (auth)/              # Auth route group
│   │   │       ├── +layout.svelte   # Auth layout wrapper
│   │   │       ├── login/+page.svelte       # Login page
│   │   │       └── signup/+page.svelte      # Signup page
│   │   └── lib/
│   │       ├── index.ts             # Lib barrel exports
│   │       ├── utils.ts             # cn() helper (clsx + tailwind-merge)
│   │       ├── assets/favicon.svg   # Application favicon
│   │       └── components/ui/       # shadcn-svelte components
│   │           └── button/          # Button component (index.ts + button.svelte)
│   ├── static/                      # Static assets served as-is
│   │   └── robots.txt              # Crawler directives
│   └── .svelte-kit/                 # SvelteKit build cache (gitignored)
└── README.md                        # This file
```

## Getting Started

### Prerequisites

- Go 1.26+
- pnpm
- PostgreSQL (or Docker with `courses-psql` container)
- Redis
- (Optional) Resend API key & AWS credentials for full functionality

### Backend (API Server)

```bash
cd backend
cp .env.example .env        # Configure DB, Redis, JWT secret, AWS, Email
make run                    # Starts the Go server on :8080
```

For hot-reloading during development:
```bash
cd backend
air                         # Uses .air.toml — watches .go files & auto-restarts
```

### Frontend (Development Server)

The SvelteKit app is configured to write its production build to `backend/dist/`, but during development you can run the Vite dev server separately:

```bash
cd frontend
pnpm install                # Install all dependencies
pnpm run dev                # Starts Vite dev server on :5173
```

The Vite dev server is configured to proxy `/api` requests to `http://localhost:8080` (the Go backend). This means during development you work at `http://localhost:5173` and all API calls are forwarded to Go automatically.

> **Note:** The frontend uses `pnpm` (not npm). Ensure you have pnpm installed (`npm install -g pnpm` or `corepack enable`).

### Production Build (Single Binary)

The frontend is embedded into the Go binary. To build a single deployable server:

```bash
# 1. Build the frontend (output → backend/dist/)
cd frontend && pnpm run build

# 2. Build the Go binary (embeds backend/dist/)
cd backend && go build -o server ./cmd/main.go

# 3. Run the single binary (serves API + SPA on :8080)
./server
```

The resulting binary is fully self-contained — static assets are compiled in, no separate file serving needed.

## Database Migrations

Migrations live in `backend/internal/database/migrations/`. Run them from `backend/`:

```bash
make migrate-up             # Apply all pending migrations
make migrate-down           # Roll back the last migration
make migrate-force v=<N>    # Force migration version to N
```

## Makefile Reference

| Target | Purpose |
|---|---|
| `run` | `go run cmd/main.go` |
| `migrate-up` | Run all pending migrations |
| `migrate-down` | Roll back last migration |
| `migrate-force v=<N>` | Force-set migration version |
| `delete-all-users` | Truncate User table (via Docker) |
