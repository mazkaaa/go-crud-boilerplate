# go-crud-boilerplate

## Stack

- **Framework**: Echo v5 (`github.com/labstack/echo/v5`)
- **Database**: PostgreSQL via `pgx/v5` (`pgxpool` connection pool), raw SQL
- **Migrations**: `golang-migrate/v4` (`file://` source, runs on boot)
- **Auth**: bcrypt via `golang.org/x/crypto`
- **Env**: `github.com/joho/godotenv` — `.env` required at startup

## Setup

1. Copy `.env.examples` to `.env` and fill in PostgreSQL credentials.
2. `go run server.go` — runs migrations + seed on boot, then starts on `:1323`.

Seeds an `admin` role + admin user from `SEED_NAME`, `SEED_EMAIL`, `SEED_PASSWORD` env vars if the roles table is empty. `FRONTEND_URL` env var is declared but unused.

## Architecture

```
server.go          — entrypoint, migrator, CORS/logger/recover middleware
migrations/
  000001_*.sql     — versioned DDL (CREATE TABLE users, roles)
config/
  database.go      — pgxpool connection (postgres:// DSN from env)
  seed.go          — roles + admin seed (raw SQL INSERT)
models/
  user.go          — UUID PK via gen_random_uuid(), json:"-" on password & role
  role.go          — UUID PK, has-many Users relationship
  response.go      — APIResponse struct
controllers/
  user.go          — raw SQL via pgxpool (SELECT/INSERT with RETURNING, manual Scan)
  role.go          — raw SQL via pgxpool, two-query eager-load in GetDetailRole
routes/
  routes.go        — GET/POST /users, GET/POST/DELETE /roles, GET /roles/:id
utils/
  response_helper.go  — uniform JSON envelope {data, message, status}
  hash_password.go    — bcrypt.DefaultCost
```

- All controllers use `c.FormValue()` — requests must be form-encoded, not JSON.
- All handlers use `c *echo.Context` (echo v5 concrete struct, not interface).
- `RoleID` in `User` is `*string` (nullable, matches `ON DELETE SET NULL` FK).
- Query params are parameterized with `$1, $2...` — safe from SQL injection.

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
If the body message is long, use bullet points or paragraphs. Footer can include breaking change notes or issue references.

## Commands

| Action | Command |
|---|---|
| Run server | `go run server.go` |
| Build | `go build -o server server.go` |
