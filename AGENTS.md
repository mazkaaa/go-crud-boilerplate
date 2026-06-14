# go-crud-boilerplate

## Stack

- **Framework**: Echo v5 (`github.com/labstack/echo/v5`)
- **Database**: PostgreSQL via `pgx/v5` (`pgxpool` connection pool), raw SQL
- **Migrations**: `golang-migrate/v4` (`file://` source, runs on boot)
- **Auth**: bcrypt via `golang.org/x/crypto`
- **Env**: `github.com/joho/godotenv` — `.env` required at startup

## Setup

1. Copy `.env.examples` to `.env` and fill in PostgreSQL credentials.
2. `go run ./cmd/server` (or `make run`) — runs migrations + seed on boot, then starts on `:1323`.

Seeds an `admin` role + admin user from `SEED_NAME`, `SEED_EMAIL`, `SEED_PASSWORD` env vars if the roles table is empty. `SERVER_PORT` env var defaults to `1323`.

## Architecture

Clean Architecture (3-layer: Handler → Service → Repository) with dependency injection in `cmd/server/main.go`.

```
cmd/server/main.go    — bootstrap, DI wiring, CORS/logger/recover middleware
migrations/
  000001_*.sql        — versioned DDL (CREATE TABLE users, roles)
internal/
  domain/             — entities + repository interfaces (no external deps)
    user.go           — User struct, UserRepository interface
    role.go           — Role struct, RoleRepository interface
    errors.go         — domain errors: ErrNotFound, ErrConflict, ErrInvalidInput
  service/            — business logic, validation, orchestration
    user_service.go   — UserService (UserRepository + RoleRepository injected)
    role_service.go   — RoleService (RoleRepository injected)
  handler/            — HTTP adapters (JSON bind, service call, response)
    user_handler.go   — GET/POST /users
    role_handler.go   — GET/POST/DELETE /roles, GET /roles/:id
  repository/         — PostgreSQL implementations of domain interfaces
    user_repo.go      — pgxpool queries, returns domain.User
    role_repo.go      — pgxpool queries, two-query eager-load in FindByIDWithUsers
  dto/                — request/response structs (CreateUserRequest, etc.)
  config/            — Config struct, env loading, DB connect, migrations, seed
  router/            — route registration (handlers injected as deps)
pkg/
  response/           — APIResponse envelope {data, message, status}
  hash/               — bcrypt.DefaultCost
Makefile               — run, build, test, fmt, vet targets
```

- All handlers accept JSON request bodies via `c.Bind()` — not form-encoded.
- Handlers use `c *echo.Context` (echo v5 concrete struct).
- `RoleID` in `User` is `*string` (nullable, matches `ON DELETE SET NULL` FK).
- Query params are parameterized with `$1, $2...` — safe from SQL injection.
- Repository interfaces are defined in `domain/` and implemented in `repository/`.
- Services depend on interfaces only (testable with mocks).

## Testing

Service layer has unit tests with manual mock repositories. Run with `make test` or `go test ./... -v`.

No CI, linter, formatter, or Docker config.

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
| Run server | `make run` or `go run ./cmd/server` |
| Build | `make build` or `go build -o bin/server ./cmd/server` |
| Test | `make test` or `go test ./... -v` |
| Format | `make fmt` or `go fmt ./...` |
| Vet | `make vet` or `go vet ./...` |
