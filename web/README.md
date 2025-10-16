# Plane Integration UI

Next.js 15 + Tailwind CSS admin UI for managing Plane ↔ CNB mappings.

## Dev

1) Install deps

```
npm ci
```

2) Run

```
API_BASE=http://localhost:8080 npm run dev
```

Open http://localhost:3000

## Env

- `API_BASE` (server): The backend base URL used by route handlers (no CORS needed).
- `NEXT_PUBLIC_API_BASE` (client): Optional; direct browser calls to the backend.

## Pages

- `/` Overview
- `/mappings` List + create Repo↔Project mappings

