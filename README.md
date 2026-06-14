# go-crud-boilerplate

A production-ready Go REST API boilerplate following Clean Architecture (Handler → Service → Repository) with PostgreSQL, Echo v5, and dependency injection.

Designed for developers who want a decoupled, testable backend without framework lock-in. Swap Postgres for any SQL database by implementing the interface in `internal/repository/`.

---

## Quick Start

```bash
# 1. Clone and enter the project
git clone <repo-url> && cd go-crud-boilerplate

# 2. Copy env template and fill in your PostgreSQL credentials
cp .env.examples .env

# 3. Run migrations, seed admin user, and start the server
make run
```

The server starts on `:1323`. On first boot it creates an `admin` role and an admin user from the `SEED_*` env vars.

---

## Tech Stack

| Layer | Choice |
|---|---|
| **Framework** | [Echo v5](https://github.com/labstack/echo/v5) |
| **Database** | PostgreSQL 16+ via [pgx/v5](https://github.com/jackc/pgx/v5) (connection pool) |
| **Migrations** | [golang-migrate/v4](https://github.com/golang-migrate/migrate) (`file://` source, runs on boot) |
| **Auth** | bcrypt via [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) |
| **Env** | [godotenv](https://github.com/joho/godotenv) — `.env` file required at startup |
| **Language** | Go 1.26 |

---

## Architecture

The project follows **Clean Architecture** (3-layer) with dependency injection wired in `cmd/server/main.go`.

```
Handler (HTTP)  →  Service (Business Logic)  →  Repository (Data Access)
     │                      │                          │
     │                      │                          │
  echo.Context           domain.UserRepo            pgxpool.Pool
  JSON bind             domain.RoleRepo            raw SQL
```

### Dependency flow

```
main.go
  → Config (env vars)
  → pgxpool.Pool          (database connection)
  → repository.UserRepo   (implements domain.UserRepository)
  → repository.RoleRepo   (implements domain.RoleRepository)
  → service.UserService   (depends on domain.UserRepository + domain.RoleRepository)
  → service.RoleService   (depends on domain.RoleRepository)
  → handler.UserHandler   (depends on service)
  → handler.RoleHandler   (depends on service)
  → router.Register       (injects handlers into Echo)
```

### Layer rules

| Layer | Dependencies | Responsibility |
|---|---|---|
| `domain/` | stdlib only | Entities, repository interfaces (ports), domain errors |
| `repository/` | `domain/`, pgxpool | PostgreSQL implementations of domain interfaces |
| `service/` | `domain/`, `pkg/hash` | Validation, orchestration, business rules |
| `handler/` | `service/`, `dto/` | Parse JSON, call service, send HTTP response |
| `router/` | `handler/`, Echo | Route registration (handlers injected as dependencies) |
| `dto/` | stdlib | Request/response structs, separate from domain entities |
| `config/` | pgxpool, migrate, godotenv, crypto | Env loading, DB connection, migrations, seed |

---

## Project Layout

```
cmd/server/main.go              — Bootstrap, DI wiring, middleware, server start
internal/
  domain/                       — Pure domain (no external imports)
    user.go                     — User entity, UserRepository interface
    role.go                     — Role entity, RoleRepository interface
    errors.go                   — ErrNotFound, ErrConflict, ErrInvalidInput
  service/
    user_service.go             — CreateUser: validate → check role → hash pw → insert
    role_service.go             — Create/Delete/GetDetail roles
    user_service_test.go        — Unit tests with mock repositories
    role_service_test.go        — Unit tests with mock repositories
  handler/
    user_handler.go             — GET/POST /users (JSON bind, service call, response)
    role_handler.go             — GET/POST/DELETE /roles, GET /roles/:id
  repository/
    user_repo.go                — pgxpool queries for UserRepository
    role_repo.go                — pgxpool queries for RoleRepository
  dto/
    user_dto.go                 — CreateUserRequest, UserResponse
    role_dto.go                 — CreateRoleRequest, RoleWithUsersResponse
  config/
    config.go                   — Config struct, env load, DB connect, migrations, seed
  router/
    router.go                   — Route registration (handlers injected)
pkg/
  response/
    response.go                 — APIResponse envelope {data, message, status}
  hash/
    hash.go                     — bcrypt password hashing
migrations/
  000001_create_tables.up.sql   — PostgreSQL DDL (users + roles)
  000001_create_tables.down.sql
Makefile                         — run, build, test, fmt, vet
```

---

## API Endpoints

All endpoints accept JSON request bodies and return JSON responses wrapped in the standard envelope:

```jsonc
{
  "data":     /* response payload or null */,
  "message":  "Success",
  "status":   200
}
```

### Users

#### `GET /users`

List all users.

**Response `200`**
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Admin User",
      "email": "admin@example.com",
      "role_id": "uuid",
      "created_at": "2026-06-15T00:00:00Z",
      "updated_at": "2026-06-15T00:00:00Z"
    }
  ],
  "message": "Success",
  "status": 200
}
```

#### `POST /users`

Create a new user.

**Request**
```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "password": "secret123",
  "role_id": "existing-role-uuid"
}
```

**Response `201`** — same user object as above.

**Errors:** `400` (invalid input), `404` (role not found), `409` (email conflict).

### Roles

#### `GET /roles`

List all roles.

**Response `200`**
```json
{
  "data": [
    { "id": "uuid", "name": "admin", "created_at": "...", "updated_at": "..." }
  ],
  "message": "Success",
  "status": 200
}
```

#### `POST /roles`

Create a new role.

**Request**
```json
{
  "name": "editor"
}
```

**Response `201`** — role object with generated UUID.

**Errors:** `400` (empty name), `409` (duplicate name).

#### `GET /roles/:id`

Get a role with its associated users.

**Response `200`**
```json
{
  "data": {
    "id": "uuid",
    "name": "admin",
    "users": [
      { "id": "uuid", "name": "Admin User", "email": "admin@example.com", "role_id": "uuid", "created_at": "...", "updated_at": "..." }
    ],
    "created_at": "...",
    "updated_at": "..."
  },
  "message": "Success",
  "status": 200
}
```

#### `DELETE /roles/:id`

Delete a role.

**Response `200`** — `data` is `null`, `message: "Role deleted successfully"`.

**Errors:** `404` (role not found).

---

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `DB_HOST` | Yes | — | PostgreSQL host |
| `DB_USER` | Yes | — | PostgreSQL user |
| `DB_PASSWORD` | Yes | — | PostgreSQL password |
| `DB_NAME` | Yes | — | PostgreSQL database name |
| `DB_PORT` | Yes | — | PostgreSQL port |
| `SEED_NAME` | Yes* | — | Admin user name (required for first boot) |
| `SEED_EMAIL` | Yes* | — | Admin user email (required for first boot) |
| `SEED_PASSWORD` | Yes* | — | Admin user password (required for first boot) |
| `SERVER_PORT` | No | `1323` | HTTP server listen port |

\* *Required only when the roles table is empty (first boot). Once seeded, they are unused.*

---

## Commands

| Action | Command |
|---|---|
| Start server | `make run` or `go run ./cmd/server` |
| Build binary | `make build` or `go build -o bin/server ./cmd/server` |
| Run tests | `make test` or `go test ./... -v` |
| Format code | `make fmt` or `go fmt ./...` |
| Static analysis | `make vet` or `go vet ./...` |

---

## Development — Adding a New Entity

To add a new feature (e.g., `articles`), follow these steps:

### 1. Domain layer (`internal/domain/article.go`)

Define the entity and the repository interface that the service will depend on:

```go
type Article struct {
	ID        string
	Title     string
	Body      string
	AuthorID  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ArticleRepository interface {
	FindAll(ctx context.Context) ([]Article, error)
	Create(ctx context.Context, title, body, authorID string) (Article, error)
}
```

### 2. Repository implementation (`internal/repository/article_repo.go`)

Implement the interface using pgxpool:

```go
type ArticleRepo struct {
	pool *pgxpool.Pool
}

func NewArticleRepo(pool *pgxpool.Pool) *ArticleRepo {
	return &ArticleRepo{pool: pool}
}

func (r *ArticleRepo) Create(ctx context.Context, title, body, authorID string) (domain.Article, error) {
	// INSERT ... RETURNING id, title, body, author_id, created_at, updated_at
}
```

### 3. Service layer (`internal/service/article_service.go`)

Add business logic:

```go
type ArticleService struct {
	repo domain.ArticleRepository
}

func (s *ArticleService) CreateArticle(ctx context.Context, title, body, authorID string) (domain.Article, error) {
	if title == "" {
		return domain.Article{}, domain.ErrInvalidInput
	}
	return s.repo.Create(ctx, title, body, authorID)
}
```

### 4. DTOs (`internal/dto/article_dto.go`)

```go
type CreateArticleRequest struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	AuthorID string `json:"author_id"`
}
```

### 5. Handler (`internal/handler/article_handler.go`)

```go
type ArticleHandler struct {
	svc *service.ArticleService
}

func (h *ArticleHandler) CreateArticle(c *echo.Context) error {
	var req dto.CreateArticleRequest
	if err := c.Bind(&req); err != nil {
		return response.Send(c, nil, "Invalid request", 400)
	}
	article, err := h.svc.CreateArticle(c.Request().Context(), req.Title, req.Body, req.AuthorID)
	if err != nil {
		return handleServiceError(c, err)
	}
	return response.Send(c, article, "Created", 201)
}
```

### 6. Wire in `cmd/server/main.go`

```go
articleRepo := repository.NewArticleRepo(pool)
articleSvc := service.NewArticleService(articleRepo)
articleHandler := handler.NewArticleHandler(articleSvc)
router.Register(e, userHandler, roleHandler, articleHandler) // update router too
```

---

## Testing

The service layer is tested with **manual mock repositories** — no external dependencies, no database required.

```go
// mockArticleRepo implements domain.ArticleRepository
type mockArticleRepo struct {
	articles []domain.Article
	err      error
}

func TestCreateArticle(t *testing.T) {
	repo := &mockArticleRepo{}
	svc := service.NewArticleService(repo)

	article, err := svc.CreateArticle(ctx, "Title", "Body", "author-1")
	// assert...
}
```

Run all tests:

```bash
make test    # or: go test ./... -v
```
