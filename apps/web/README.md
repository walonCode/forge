# Forge Web

A dashboard UI for the Forge API — built with Next.js 15, TypeScript, Tailwind CSS, and TanStack Query.

## Stack

- **Framework**: Next.js 16 (App Router)
- **Language**: TypeScript (strict)
- **Styling**: Tailwind CSS v4
- **Data fetching**: TanStack Query v5
- **API types**: generated from `openapi.json` via openapi-typescript
- **Linter/Formatter**: Biome
- **Package manager**: Bun

## Quick start

See [setup.md](./setup.md) for first-time setup.

```bash
bun dev
```

Open [http://localhost:3000](http://localhost:3000).

## Pages

| Route | Description |
|---|---|
| `/` | Health — polls `GET /health` every 5s, shows status dot, uptime, DB, version |
| `/metrics` | Metrics — fetches `GET /metrics`, displays GC latency and Go runtime stats |
| `/data` | Data — fetches `GET /tasks` (requires auth), shows last 20 records in a table |
| `/auth` | Auth — login / signup forms, stores JWT in memory, verifies token |

## Scripts

| Command | Description |
|---|---|
| `bun dev` | Start dev server on port 3000 |
| `bun build` | Production build |
| `bun lint` | Run Biome checks |
| `bun format` | Format with Biome |
| `bun run gen:api` | Regenerate API types from `openapi.json` |

## Project structure

```
apps/web/
├── app/
│   ├── layout.tsx       # Root layout with nav and providers
│   ├── page.tsx         # / — health page
│   ├── metrics/
│   │   └── page.tsx
│   ├── data/
│   │   └── page.tsx
│   └── auth/
│       └── page.tsx
├── lib/
│   ├── api.ts           # Typed fetch wrappers
│   ├── openapi.d.ts     # Auto-generated from openapi.json (do not edit)
│   ├── parse-metrics.ts # Prometheus text format parser
│   └── providers.tsx    # QueryClientProvider + TokenContext
└── openapi.json         # Copy of the API swagger spec
```
