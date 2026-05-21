# Forge API

A task management REST API built with Go, Chi, PostgreSQL, and Swagger.

## Stack

- **Language**: Go 1.26
- **Router**: [go-chi/chi](https://github.com/go-chi/chi)
- **Database**: PostgreSQL (via `pgx`)
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Auth**: JWT (`golang-jwt/jwt`)
- **Docs**: [swaggo/swag](https://github.com/swaggo/swag)
- **Metrics**: Prometheus

## Quick start

See [setup.md](./setup.md) for full prerequisites and first-time setup.

```bash
make run        # start Docker services + API
```

Swagger UI is available at `http://localhost:8080/swagger/index.html` once the server is running.

## Available commands

| Command | Description |
|---|---|
| `make run` | Start Docker services and run the API |
| `make build` | Build binaries for the current and cross-compile targets |
| `make audit` | Format, vet, and test the codebase |
| `make swagger` | Regenerate Swagger docs from annotations |
| `make swagger_fmt` | Format Swagger comment annotations |
| `make database_up` | Apply all pending migrations |
| `make database_down` | Roll back the last migration |
| `make migrate_create name=<name>` | Create a new migration file pair |

## Project structure

```
apps/api/
├── cmd/api/          # Entry point
├── db/migrations/    # SQL migration files
├── docs/             # Generated Swagger docs (do not edit)
├── internals/
│   ├── middleware/   # Logger, CORS, correlation ID
│   ├── modules/      # Feature modules (auth, tasks, health)
│   └── server/       # Server setup and routing
└── pkg/
    ├── configs/      # Environment config
    ├── database/     # DB connection
    └── utils/        # JWT, response helpers, env
```

## API endpoints

### Health
| Method | Path | Description |
|---|---|---|
| GET | `/health` | Liveness check |
| GET | `/health/ready` | Readiness check |

### Auth
| Method | Path | Description |
|---|---|---|
| POST | `/auth/signup` | Register a new user |
| POST | `/auth/login` | Login and get tokens |

### Tasks (requires `Authorization: Bearer <token>`)
| Method | Path | Description |
|---|---|---|
| POST | `/task` | Create a task |
| GET | `/tasks` | Get all user tasks |
| DELETE | `/task/{id}` | Delete a task |
