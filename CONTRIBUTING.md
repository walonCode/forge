# Contributing

## Before you start

- Check open issues before opening a new one
- For anything non-trivial, open an issue first and describe what you want to change
- One focused change per PR — no bundled refactors

---

## Setup

Follow [`apps/api/setup.md`](./apps/api/setup.md) and [`apps/web/setup.md`](./apps/web/setup.md) to get the full stack running locally before making any changes.

---

## Go API

### Module pattern

Every feature lives in its own module under `internals/modules/`. Each module contains:

```
internals/modules/<name>/
├── domain.go      # types (request, response, DB structs)
├── repository.go  # database queries only — no business logic
├── service.go     # business logic — calls repository
├── handler.go     # HTTP — decodes request, calls service, encodes response
└── module.go      # wires repo → service → handler, exposes Register(r)
```

Write them in that order. `domain.go` first, `handler.go` last.

### Migrations first

Always write the migration before any Go code. Run it. Confirm the table exists in the database. Then write the queries.

```bash
make migrate_create name=<migration_name>  # creates up + down files
make database_up                           # apply
make database_down                         # rollback one
```

### Swagger annotations

Every new handler needs swaggo annotations. After adding them:

```bash
make swagger       # regenerates docs/
make swagger_fmt   # formats annotations
```

### Tests

- **Service layer**: unit tests in `service_test.go`, no database needed
- **Repository layer**: integration tests in `repository_test.go` against a real Postgres instance
- Use transaction rollback pattern — each test wraps in a `tx`, defers `tx.Rollback()`
- Target: 75%+ coverage on service and repository layers

```bash
make audit   # fmt + vet + test
```

### Linting

```bash
go vet ./...
golangci-lint run ./...   # zero warnings before committing
```

### Return 404 not 403

When a user tries to access another user's resource, return `404`. A `403` confirms the resource exists — `404` reveals nothing.

---

## Next.js UI

```bash
cd apps/web
bun lint      # Biome
bun build     # type-check + build
```

If you change the API (new endpoints, modified responses), regenerate the API types:

```bash
cp apps/api/docs/swagger.json apps/web/openapi.json
cd apps/web && bun run gen:api
```

The UI scope is four pages only. Do not add features beyond what is in the four-page spec.

---

## Commit style

Follow the existing commit style in this repo:

```
[feats]: add task deletion endpoint
[updates]: update auth middleware to extract userId
[fix]: handle empty username in login handler
[docs]: update api setup guide
[chore]: run go mod tidy
```

---

## Pull request checklist

Before opening a PR:

- [ ] `go vet ./...` passes
- [ ] `golangci-lint run ./...` passes with zero errors
- [ ] `go test ./...` passes locally
- [ ] New handlers have Swagger annotations and `make swagger` has been run
- [ ] New migrations have both `.up.sql` and `.down.sql`
- [ ] `bun build` passes if you touched the UI
- [ ] No TODO comments left in committed code (use GitHub Issues instead)
