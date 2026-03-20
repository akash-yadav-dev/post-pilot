# Post Pilot

Post Pilot is a monorepo for a social publishing platform with:

- Go API service
- Next.js web app
- Go database migrator
- Worker app scaffold
- Shared packages for database, cache, logging, and security

## Monorepo Layout

- apps/api: Gin HTTP API (auth, users, posts)
- apps/web: Next.js frontend
- apps/migrator: migration CLI (golang-migrate)
- apps/worker: background worker scaffold
- packages: shared Go packages (cache, database, logger, security, queue, utils)
- deployments/docker: local dev Docker stack
- deployments/k8s: Kubernetes manifests
- scripts: helper scripts for lint and migrations
- docs: design/docs placeholders

## Tech Stack

- Go 1.26.1 workspace (single module)
- Gin for HTTP API
- PostgreSQL as primary database
- Redis for token revocation and rate limiting
- Next.js 16 + React 19 frontend
- golang-migrate for schema migrations

## Implemented Backend Surface

Base groups are versioned under /api/v1.

### Health

- GET /health

### Auth

- POST /api/v1/auth/register
- POST /api/v1/auth/login
- POST /api/v1/auth/refresh
- POST /api/v1/auth/logout
- GET /api/v1/auth/me

Auth module includes:

- JWT access/refresh tokens
- Redis-backed refresh token store
- IP and email-aware login/register/refresh rate limits
- DB-backed audit logging hooks

### Users

- POST /api/v1/users
- GET /api/v1/users/:id
- PATCH /api/v1/users/:id
- DELETE /api/v1/users/:id

Users persistence is implemented against PostgreSQL.

### Posts

- POST /api/v1/posts
- GET /api/v1/posts
- GET /api/v1/posts/:id
- PATCH /api/v1/posts/:id
- DELETE /api/v1/posts/:id

Note: routes and service wiring exist, but post repository implementation is still stubbed (panic placeholders).

## Frontend (apps/web)

Implemented UI routes:

- /
- /login
- /signin
- /dashboard

Frontend includes:

- cookie-based token storage in browser
- route protection middleware for login/dashboard navigation
- auth service calls for login/register/me/logout

## Worker Status

Worker entrypoint exists but currently prints startup text only.
Job files are scaffolds with package declarations and no processing logic yet.

## Database Migrations

Migration files live in packages/database/migrations and include schemas for:

- users
- auth_accounts
- social_accounts
- posts
- post_targets
- jobs
- job_runs
- analytics_events
- scheduler_locks
- audit_logs

Use the migrator app via script:

```bash
./scripts/migrate.sh up
./scripts/migrate.sh version
./scripts/migrate.sh steps --steps 1
```

Destructive commands require explicit confirmation:

```bash
./scripts/migrate.sh down --yes
```

## Environment Variables

Copy and adjust .env.example for local values.

High-impact groups:

- App: APP_PORT, APP_BASE_URL
- Database: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
- JWT: JWT_ACCESS_SECRET_KEY, JWT_REFRESH_SECRET_KEY, JWT_EXPIRY, JWT_REFRESH_EXPIRY
- Auth hardening: BCRYPT_COST, AUTH_RATE_LIMIT_*
- Redis: REDIS_ADDR plus pool/timeouts/tls keys
- CORS: CORS_ALLOWED_ORIGINS

Important runtime behavior: API startup currently requires Redis connectivity because Redis client initialization happens during router bootstrap.

## Local Development (without Docker)

Prerequisites:

- Go 1.26.1
- Node.js 22+
- PostgreSQL
- Redis

### 1) Start API

```bash
go run ./apps/api/cmd/server
```

### 2) Start frontend

```bash
cd apps/web
npm install
npm run dev
```

### 3) Run migrations

```bash
./scripts/migrate.sh up
```

## Docker Development Stack

A full local stack is available in deployments/docker:

- postgres
- redis
- api (with hot reload via air)
- web
- migrator profile service

Start stack:

```bash
docker compose -f deployments/docker/docker-compose.yml up --build
```

Run migrations in Docker:

```bash
docker compose -f deployments/docker/docker-compose.yml --profile tools run --rm migrator -command up
```

Detailed Docker usage is documented in [deployments/docker/README.md](deployments/docker/README.md).

## Quality Checks

Run all Go tests:

```bash
go test ./...
```

Run lint script:

```bash
./scripts/lint.sh
```

Frontend checks:

```bash
cd apps/web
npm run lint
npm run typecheck
```

## Current Caveats

These are code-level caveats to keep in mind while developing:

- Frontend auth service currently calls /api/auth/* while backend routes are /api/v1/auth/*.
- Posts handlers read user_id from context while auth middleware currently sets auth_user_id.
- Posts repository methods are not implemented yet.
- Worker and scheduler internals are still scaffolds.

## Related Files

- API route wiring: [apps/api/routes/routes.go](apps/api/routes/routes.go)
- API bootstrap/router: [apps/api/cmd/server/bootstrap/router.go](apps/api/cmd/server/bootstrap/router.go)
- Example env vars: [.env.example](.env.example)
- Migrator runner: [apps/migrator/internal/runner/runner.go](apps/migrator/internal/runner/runner.go)
- Docker dev docs: [deployments/docker/README.md](deployments/docker/README.md)
