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

## Theme

- Supports light and dark themes.
- Toggle in the header; preference is saved in `localStorage` (`theme=light|dark`).
- Default is dark for backward compatibility; if no preference is saved, the UI follows system preference on first load.

## Env

- `API_BASE` (server): The backend base URL used by route handlers (no CORS needed).
- `NEXT_PUBLIC_API_BASE` (client): Optional; direct browser calls to the backend.

## Pages

- `/` Overview
- `/mappings` List + create Repo↔Project mappings
