# Setup

## Prerequisites

| Tool | Install |
|---|---|
| Node.js 20+ | https://nodejs.org |
| Bun | `curl -fsSL https://bun.sh/install \| bash` |

The Forge API must be running on port 8080. See `apps/api/setup.md` to get it started.

## First-time setup

**1. Install dependencies**

```bash
cd apps/web
bun install
```

**2. Configure environment**

Create `.env.local`:

```env
API_URL=http://localhost:8080
```

`API_URL` is a **server-only** variable (no `NEXT_PUBLIC_` prefix) — it is read by
Server Components and Server Actions and never shipped to the browser.

**3. Start the dev server**

```bash
bun dev
```

Open [http://localhost:3000](http://localhost:3000).

---

## Regenerating API types

The file `lib/openapi.d.ts` is auto-generated from `openapi.json`. If the API changes, update it:

```bash
# 1. Copy the latest swagger spec from the API
cp ../api/docs/swagger.json openapi.json

# 2. Regenerate types
bun run gen:api
```

Do not edit `lib/openapi.d.ts` by hand.

---

## Running checks

```bash
bun lint      # Biome lint
bun format    # Biome format (writes in place)
bun build     # Type-check + production build
```
