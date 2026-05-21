# Setup

## Prerequisites

| Tool | Version | Install |
|---|---|---|
| Go | 1.26+ | https://go.dev/dl |
| Docker | latest | https://docs.docker.com/get-docker |
| golang-migrate CLI | latest | `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest` |
| swag CLI | v1.16+ | `go install github.com/swaggo/swag/cmd/swag@latest` |

## First-time setup

**1. Clone and enter the API directory**

```bash
git clone <repo-url>
cd apps/api
```

**2. Copy the example environment file**

```bash
cp .env.example .env
```

Then fill in your values:

```env
DATABASE_URL=postgresql://username:password@localhost:5432/forge_api?sslmode=disable
PORT=8080
APP_VERSION=1.0.0
JWT_SECRET=your-secret-here
```

**3. Start the database**

```bash
docker compose up -d
```

This starts PostgreSQL on port `5432` and pgAdmin on port `5050`.

> pgAdmin login: `admin@example.com` / `password`

**4. Run migrations**

```bash
make database_up
```

**5. Start the API**

```bash
go run ./cmd/api
# or
make run
```

The API will be available at `http://localhost:8080`.

---

## Working with migrations

**Create a new migration**

```bash
make migrate_create name=add_tags_to_tasks
```

This creates two files in `db/migrations/`:
- `NNN_add_tags_to_tasks.up.sql` — forward migration
- `NNN_add_tags_to_tasks.down.sql` — rollback migration

Always fill in both files before applying.

**Apply migrations**

```bash
make database_up
```

**Roll back the last migration**

```bash
make database_down
```

---

## Working with Swagger

After adding or editing annotations in any handler file, regenerate the docs:

```bash
make swagger
```

Then open `http://localhost:8080/swagger/index.html`.

To format your annotation comments:

```bash
make swagger_fmt
```

---

## Running checks

```bash
make audit    # go mod tidy + verify, fmt, vet, test
```
