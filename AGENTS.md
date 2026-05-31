# go-crud-boilerplate

## Stack

- **Framework**: Echo v4 (`github.com/labstack/echo/v4`)
- **ORM**: GORM (`gorm.io/gorm`)
- **Database**: PostgreSQL (`gorm.io/driver/postgres`, `pgx/v5`)
- **Auth**: bcrypt via `golang.org/x/crypto`, JWT lib present but not wired
- **Env**: `github.com/joho/godotenv` — `.env` required at startup

## Setup

1. Copy `.env.examples` to `.env` and fill in PostgreSQL credentials.
2. `go run server.go` — starts on `:1323`, runs auto-migrate + seed on boot.

Seeds an `admin` role + admin user from `SEED_NAME`, `SEED_EMAIL`, `SEED_PASSWORD` env vars if the roles table is empty. `FRONTEND_URL` env var is declared but unused.

## Architecture

```
server.go          — entrypoint, CORS, logger middleware
config/
  database.go      — GORM connection (DSN from env, sslmode=disable)
  seed.go          — roles + admin seed
models/
  user.go          — UUID PK via gen_random_uuid(), json:"-" on password & role
  role.go          — UUID PK, has-many Users relationship
  response.go      — APIResponse struct
controllers/
  user.go          — FormValue-based (not JSON body), hash_password + Role lookup
  role.go          — name uniqueness check, Preload("Users") on detail
routes/
  routes.go        — GET/POST /users, GET/POST /roles, GET /roles/:id
utils/
  response_helper.go  — uniform JSON envelope {data, message, status}
  hash_password.go    — bcrypt.DefaultCost
```

- All controllers use `c.FormValue()` — requests must be form-encoded, not JSON.
- `DeleteRole` is defined in `controllers/role.go` but **not registered** in routes.
- All handler errors use `log.Fatal` (crashes the server) rather than graceful error returns.

## Testing

No tests exist. No CI, linter, formatter, or Docker config.

## Git

All commits must follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]
[optional footer]
```

Types: `feat`, `fix`, `build`, `chore`, `ci`, `docs`, `perf`, `refactor`, `revert`, `style`, `test`.
Scope is optional (e.g., `auth`, `routes`, `db`, `users`, `roles`). Start description with lowercase, no period. Use imperative mood.

## Commands

| Action | Command |
|---|---|
| Run server | `go run server.go` |
| Build | `go build -o server server.go` |
