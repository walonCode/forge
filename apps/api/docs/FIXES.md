# Forge API — Fixes, Hardening & Tests

This document records the correctness bugs, security hardening, and test coverage
added to the API. Everything below has been built (`go build ./...`), vetted
(`go vet ./...`), unit-tested (`go test ./...`), and — where it has runtime
behaviour — verified end-to-end against a live server + PostgreSQL.

---

## 1. Correctness bugs

### 1.1 `tasks` repository — column/query bugs
`internals/modules/tasks/repository.go`

- `Scan` argument order did not match the `SELECT` column order, so reads either
  populated the wrong fields or failed outright.
- `DeleteTask` / `GetTask` / `UpdateTask` filtered on a non-existent `taskId`
  column instead of `id`.
- `UpdateTask` never refreshed `updated_at`.
- `GetTasks` did not check `rows.Err()` and registered `defer rows.Close()` after
  the scan loop (leak on early return).

**Fix:** explicit column lists that match the scan order, `WHERE ... id = $n`,
`updated_at = now()`, `rows.Err()` check, and `defer rows.Close()` right after the
query.

### 1.2 `DeleteTask` — swapped arguments
`internals/modules/tasks/handler.go`

The handler called `service.DeleteTask(ctx, taskId, userId)` against a
`(ctx, userId, taskId)` signature, so the query became
`WHERE userId = <taskId> AND id = <userId>` — it matched nothing and returned
`204` while deleting nothing.

**Fix:** call `service.DeleteTask(ctx, userId, taskId)`.
**Verified:** delete now returns `204` and the row is actually removed
(DB row count → 0).

### 1.3 `UpdateTask` — body decoded as a bare bool
`internals/modules/tasks/handler.go`

The handler decoded a raw JSON boolean into `var isCompleted bool`, but the
documented body (and `UpdateTaskRequest`) is an object `{"isCompleted": true}`,
so the documented request always failed with `400`.

**Fix:** decode `UpdateTaskRequest` and use `request.IsCompleted`.
**Verified:** `PATCH /task/{id}` with `{"isCompleted":true}` returns `200` and
persists `is_completed = true`.

### 1.4 `GetTasks` double pointer
Repository/service returned `*[]Task` and the handler wrapped it again as `&task`
(`**[]Task`). Changed the whole path to return `[]Task`.

---

## 2. Data integrity

### 2.1 Unique usernames + signup race
`db/migrations/003_add_unique_username.up.sql`, `internals/modules/auth`

`users.username` had no `UNIQUE` constraint, and signup used a
check-then-insert (`FindUserByUsername` then `CreateUser`) — a TOCTOU race that
let two concurrent signups create duplicate usernames, which in turn made login
(`QueryRow` → first row) ambiguous.

**Fix:**
- New migration adds `CREATE UNIQUE INDEX users_username_unique ON users (username)`.
- Signup dropped the pre-check and now relies on the unique index, mapping a
  unique-constraint violation to `409 Conflict`
  (`database.IsUniqueViolation`, SQLSTATE `23505`).
- `users` profile update maps the same violation to `409` as a safety net.

**Verified:** a second signup with the same username returns `409`.

---

## 3. Security

### 3.1 JWT secret is required and injected
`pkg/configs`, `pkg/utils/jwt.go`, `internals/modules/auth`, `internals/server`

- Removed the hardcoded fallback secret that was committed in source — if
  `JWT_SECRET` had been unset in production the API would have signed tokens with
  a publicly known key (token forgery / auth bypass).
- `configs.Load()` now returns an error when `JWT_SECRET` or `DATABASE_URL` is
  missing (fail-fast at startup).
- The secret is loaded once in config and injected via dependency injection:
  `CreateToken(id, secret, type, ttl)`, `VerifyToken(token, secret)`,
  `AuthMiddleware(secret)`. The package-level `var secretKey` is gone.

### 3.2 JWT algorithm confusion guard
`VerifyToken` rejects any token not signed with HMAC (blocks `alg=none` and
RS/HS confusion attacks).

### 3.3 Access vs refresh tokens
`pkg/utils/jwt.go`, `internals/modules/auth`

Previously both tokens were identical (same claims and TTL), so a "refresh"
token bought nothing.

- Tokens now carry a `typ` claim (`access` / `refresh`) and distinct TTLs
  (access 24h, refresh 7d).
- `AuthMiddleware` rejects a refresh token used as an access token.
- New `POST /auth/refresh` exchanges a valid refresh token for a fresh
  access/refresh pair; presenting an access token there is rejected.

**Verified:** refresh-as-access → `401`; `/auth/refresh` with a refresh token →
`200` + new pair; `/auth/refresh` with an access token → `401`.

### 3.4 Rate limiting on auth endpoints
`internals/server/route.go`

`/auth/*` is wrapped in a per-client-IP limiter (`httprate`, 20 req/min) to blunt
credential brute-forcing. The key is the client IP resolved by the existing
`RealIP` middleware.

### 3.5 Request body size cap
`internals/middleware/body_limit.go`

A global `MaxBodyBytes` middleware caps request bodies at 1 MB via
`http.MaxBytesReader`, so handlers can't buffer unbounded input.

### 3.6 Secrets no longer logged
`internals/server/server.go`

Startup previously logged the entire config struct (JWT secret + DB URL with
password). It now logs only the app version.

### 3.7 CORS
`AllowCredentials` set to `false` — the API authenticates with the
`Authorization` header, not cookies, and credentialed mode is invalid with a
wildcard origin.

---

## 4. Input validation

`pkg/utils/validate.go` and the handlers

The structs carried `validate:"..."` tags, but nothing ran them, so e.g. signup
accepted a 1-character password. A shared validator (`go-playground/validator`)
is now invoked in the signup, login, task-create, profile-update, and
password-change handlers.

**Verified:** signup with a short password → `400`.

---

## 5. Startup / config robustness

- `pkg/database/database.go` — `Connect` now `Ping`s the database (5s timeout) so
  a bad `DATABASE_URL` fails fast instead of erroring on the first query.
  **Verified:** the server refused to start when PostgreSQL was down.
- `internals/server/server.go` — the HTTP server now binds `cfg.Port` (was a
  hardcoded `:8080`); the Swagger URL uses the port too.
- `pkg/configs/configs.go` — typed `*Config` (was `any`), fixed the
  `DATABSE_URL` → `DatabaseURL` typo, default port `8080`.
- `internals/modules/health` — `Health()` reads the version from config instead
  of a hardcoded `"1.0.0"`.

---

## 6. Consistency / smaller fixes

- `auth.CreateUser` no longer ignores the id-generation error.
- bcrypt cost extracted to a `bcryptCost` constant.
- `auth` repository local query variables renamed `sql` → `query`
  (stopped shadowing the `database/sql` package).
- `auth.FindUserByUsername` parameter renamed `email` → `username`.
- `AuthResponse.Data` is now a pointer so it is omitted (not emitted as empty
  tokens) on logout.
- `middleware` context key type renamed `contextkey` → `contextKey`.
- `AuthMiddlware` → `AuthMiddleware`; the fixed `logout` handler replaced the
  `nil` route that would have panicked.

---

## 7. Tests

New tests (all passing) — `make test`:

| Package | File(s) | Covers |
|---|---|---|
| `pkg/utils` | `env_test.go`, `jwt_test.go` | env parsing; token round-trip, type claim, wrong-secret, expired, tampered, empty-secret |
| `pkg/configs` | `configs_test.go` | required-var validation + defaults (100%) |
| `auth` | `middleware_test.go`, `service_test.go`, `handler_test.go` | middleware (valid/missing/malformed/bad/refresh-rejected), password hashing, login/signup/logout status codes |
| `tasks` | `service_test.go`, `handler_test.go` | field mapping, **delete arg-order guard**, object-body decode, create validation |
| `users` | `service_test.go`, `handler_test.go` | profile/password/delete logic, HTTP status mapping, no-password-leak |

A `make test` / `make test_coverage` target was added.

---

## 8. How to verify locally

```bash
cd apps/api
make run           # docker compose up + go run   (needs JWT_SECRET & DATABASE_URL in .env)
make database_up   # apply migrations, including 003_add_unique_username
make test          # run the suite
make swagger       # regenerate docs (includes /auth/refresh and the users module)
```

Swagger UI: `http://localhost:8080/swagger/index.html`.
