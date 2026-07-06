# Contributing

## Heads up

Forge is a personal passion project. It is built for learning and is not maintained as a collaborative open-source project.

Pull requests and issues are welcome but may not be reviewed, responded to, or merged. If you want to build on this, fork it and make it your own — that is the most honest way to use it.

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

Modules are tested without a live database by implementing the module's
`Repository` interface with an in-memory fake:

- **Service layer** (`service_test.go`): construct the service with a fake
  repository and assert business logic — validation, error mapping, hashing.
- **Handler layer** (`handler_test.go`): drive the handler with `httptest` and a
  fake-backed service; assert status codes, response shape, and that no secret or
  password leaks.
- Prefer table-driven tests. Repository integration tests against a real Postgres
  instance are welcome where SQL correctness matters.

```bash
make test            # go test ./...
make test_coverage   # go test -cover ./...
make audit           # tidy + verify + fmt + vet + test
```

### Linting

```bash
go vet ./...
gofmt -l .   # should print nothing
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

Use [Conventional Commits](https://www.conventionalcommits.org/), scoped to the
area you touched. One focused change per commit.

```
feat(tasks): add task deletion endpoint
fix(auth): handle empty username in login handler
refactor(auth): extract userId in middleware
build(deps): add validator dependency
docs: update api setup guide
chore: run go mod tidy
```

---

## Continuous integration

Pushes to `main` trigger the [`test`](.github/workflows/test.yml) workflow, which
builds, vets, and tests the Go API and type-checks the web app. Make sure the
checklist below passes locally before pushing to `main`.

## Pull request checklist

Before opening a PR:

- [ ] `make test` passes locally (`go build`, `go vet`, `go test`)
- [ ] `gofmt -l .` prints nothing
- [ ] New handlers have Swagger annotations and `make swagger` has been run
- [ ] New migrations have both `.up.sql` and `.down.sql`
- [ ] API changes are synced to the UI (`openapi.json` + `bun run gen:api`)
- [ ] `bunx tsc --noEmit` passes if you touched the UI
- [ ] No TODO comments left in committed code (use GitHub Issues instead)
