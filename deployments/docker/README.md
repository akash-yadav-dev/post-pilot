# Post Pilot Dev Deployment

## Start frontend + backend + databases

```bash
docker compose -f deployments/docker/docker-compose.yml up --build
```

Services:
- Web (Next.js): http://localhost:3000
- API (Go): http://localhost:8080
- PostgreSQL: localhost:5432
- Redis: localhost:6379

## Backend auto-reload

The API container uses `air` with config in `deployments/docker/.air.api.toml`.
Any Go file change under the mounted repo rebuilds and restarts the API automatically.

## Run migrations on demand

Run all pending migrations:

```bash
docker compose -f deployments/docker/docker-compose.yml --profile tools run --rm migrator -command up
```

Check current migration version:

```bash
docker compose -f deployments/docker/docker-compose.yml --profile tools run --rm migrator -command version
```

Run down migration (destructive):

```bash
docker compose -f deployments/docker/docker-compose.yml --profile tools run --rm migrator -command down -yes
```

Generate Swagger endpoint catalog (dev tooling):

```bash
docker compose -f deployments/docker/docker-compose.yml --profile tools run --rm swagger
```

Generated output:
- docs/swagger/swagger.json

## Stop stack

```bash
docker compose -f deployments/docker/docker-compose.yml down
```

Remove volumes too:

```bash
docker compose -f deployments/docker/docker-compose.yml down -v
```
