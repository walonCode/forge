# forge

A monorepo that combines tools — Go API, Next.js dashboard, Bash/Python scripts, AWS infrastructure, and a Rust CLI — to make a project work seamlessly end to end.

> **Status**: Layer 1 (Go API) in progress · Layer 3 (UI) complete · Layers 2, 4, 5, 6 coming.

---

## What this is

Forge is a monorepo that wires together every layer of a real backend system. Each tool is chosen deliberately — the goal is to understand how they fit together, not just that they work individually:

| Layer | What | Status |
|---|---|---|
| 1 | Go API — net/http, sqlc, JWT, Prometheus, Chi | In progress |
| 2 | GitHub Actions CI/CD | Coming soon |
| 3 | Next.js dashboard | Done |
| 4 | Bash + Python automation scripts | In progress |
| 5 | AWS — ECS, RDS, S3, CDK, Terraform | Coming soon |
| 6 | Rust CLI (`blast`) — load tester | Coming soon |

---

## Stack

**Backend**
- Go 1.22+ with Chi router
- PostgreSQL 16 via `pgx`
- `sqlc` for type-safe queries
- `golang-migrate` for migrations
- `golang-jwt` for auth
- Prometheus metrics

**Frontend**
- Next.js 15 (App Router), TypeScript strict
- Tailwind CSS v4, TanStack Query v5
- API types generated from OpenAPI spec via `openapi-typescript`

**Infrastructure** *(coming soon)*
- AWS ECS Fargate + RDS + S3 + CloudFront
- AWS CDK (TypeScript) + Terraform
- GitHub Actions CI/CD with OIDC federation

**CLI** *(coming soon)*
- Rust — `blast` traffic generator and load tester

**Scripts**
- Python — async seeder with `httpx`, `Faker`, and `rich` progress output *(available)*
- Bash — setup, migrate, seed, healthcheck *(coming soon)*

---

## Repository structure

```
forge/
├── apps/
│   ├── api/          # Go backend
│   │   ├── cmd/api/  # Entry point
│   │   ├── db/       # Migrations + sqlc queries
│   │   ├── docs/     # Generated Swagger docs
│   │   └── internals/
│   │       ├── middleware/
│   │       ├── modules/  # auth · health · tasks
│   │       └── server/
│   └── web/          # Next.js dashboard
│       ├── app/      # health · metrics · data · auth pages
│       └── lib/      # API client · types · providers
├── cli/              # Rust blast CLI (coming soon)
├── infra/
│   ├── cdk/          # AWS CDK TypeScript stacks (coming soon)
│   └── terraform/    # Terraform HCL (coming soon)
├── scripts/
│   └── python/       # Async seeder — users + tasks via httpx
└── .github/
    └── workflows/    # CI/CD pipelines (coming soon)
```

---

## Quick start

**Requirements**: Go 1.22+, Docker, Bun

```bash
# 1. Clone
git clone https://github.com/walonCode/forge.git
cd forge

# 2. Start the database
docker compose up -d

# 3. Run the API
cd apps/api
cp .env.example .env   # fill in your values
make database_up       # run migrations
make run               # starts on :8080

# 4. Run the UI (separate terminal)
cd apps/web
bun install
bun dev                # starts on :3000
```

See [`apps/api/setup.md`](./apps/api/setup.md) and [`apps/web/setup.md`](./apps/web/setup.md) for full setup details.

---

## API

| Route | Description |
|---|---|
| `GET /health` | Liveness check |
| `GET /health/ready` | Readiness check |
| `POST /auth/signup` | Register |
| `POST /auth/login` | Login |
| `POST /task` | Create task (auth required) |
| `GET /tasks` | List tasks (auth required) |
| `DELETE /task/{id}` | Delete task (auth required) |
| `GET /swagger/*` | Swagger UI |
| `GET /metrics` | Prometheus metrics |

---

## Scripts

### Python seeder

Seeds the API with fake users and tasks using `httpx` and `Faker`.

**Requirements**: Python 3.14+, [`uv`](https://github.com/astral-sh/uv)

```bash
cd scripts/python
uv sync

# seed 10 users with 5 tasks each (defaults)
uv run main.py

# custom counts
uv run main.py --users 50 --tasks-per-user 10

# point at a different API
API_URL=http://localhost:9090 uv run main.py
```

---

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).

## Security

See [SECURITY.md](./SECURITY.md).

## License

[MIT](./LICENSE)
