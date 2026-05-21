# forge — Python scripts

Async scripts for seeding and testing the Forge API. Built with `httpx`, `Faker`, and `rich`.

---

## Requirements

- Python 3.14+
- [`uv`](https://github.com/astral-sh/uv) — package manager

---

## Setup

```bash
# install dependencies
uv sync
```

---

## Scripts

### `main.py` — seeder

Seeds the API with fake users and tasks. Each user is registered, then a configurable number of tasks are created under that account.

```bash
# defaults: 10 users, 5 tasks each
uv run main.py

# custom counts
uv run main.py --users 50 --tasks-per-user 10

# point at a different host
API_URL=http://localhost:9090 uv run main.py
```

**Options**

| Flag | Default | Description |
|---|---|---|
| `--users` | `10` | Number of users to create |
| `--tasks-per-user` | `5` | Tasks to create per user |

**Environment variables**

| Variable | Default | Description |
|---|---|---|
| `API_URL` | `http://localhost:8080` | Base URL of the Forge API |

The seeder checks the `/health` endpoint before running — if the API is unreachable it exits early with a clear message.

---

### `test_create_task.py` — quick task test

One-shot script for manually verifying that task creation works with a known token. Useful when debugging auth or the task endpoint.

```bash
uv run test_create_task.py
```

Update the `token` variable inside the file with a valid JWT before running.

---

## Dependencies

| Package | Purpose |
|---|---|
| `httpx` | Async HTTP client |
| `faker` | Fake user/task data generation |
| `rich` | Progress bars and console output |
| `pydantic` | Data validation |
