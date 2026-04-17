# LogVault — Claude Code Context

## Project Overview

LogVault is a lightweight, self-hosted Docker application for browsing, downloading, and live-tailing log files via a clean web UI. Written entirely in pure Go stdlib with no external dependencies. Designed for containerized environments behind a reverse proxy.

## Repository Layout

```
logvault/
├── app/
│   ├── main.go          # Entire backend + embedded HTML/CSS/JS (~857 lines)
│   └── go.mod           # Go 1.26, zero external deps
├── Dockerfile           # Multi-stage: golang:1.26-alpine → scratch (~5 MB image)
├── docker-compose.yml   # Reference deployment with volume/env examples
├── readme.md            # User-facing docs (features, config, reverse proxy)
├── CLAUDE.md            # This file
└── .gitignore
```

`app_bk/` is a local backup folder — it is gitignored and should never be modified or deployed.

## Architecture

**Single-file Go backend** (`app/main.go`). All HTML/CSS/JS templates are embedded as Go string constants — no separate frontend build step.

### URL Routes

| Route | Handler | Auth |
|---|---|---|
| `/` | redirect → `/browse/` | yes |
| `/browse/<path>` | directory listing | yes |
| `/download/<path>` | file streaming | yes |
| `/tail/<path>` | tail viewer page | yes |
| `/tail-stream/<path>` | SSE stream | yes |
| `/login` | auth page | no |
| `/logout` | clear session | no |
| `/health` | health check | no (always public) |

BASE_PATH is prepended to all routes via the `p()` helper.

### Key Sections in main.go

| Lines | Responsibility |
|---|---|
| 1–38 | Imports, global config vars, constants |
| 39–77 | In-memory session store (mutex, 8h TTL, crypto/rand tokens) |
| 79–87 | `requireAuth` middleware |
| 100–410 | Embedded HTML templates (`loginTmpl`, `browserTmpl`, `tailTmpl`) |
| 410–465 | Template helpers: `formatSize`, `buildCrumbs`, `parentURL` |
| 466–502 | `listDir` — directory walker with sort |
| 536–620 | `/login` and `/logout` handlers |
| 622–657 | `/browse/` handler |
| 659–690 | `/download/` handler |
| 692–789 | `/tail/` page + `/tail-stream/` SSE handler |
| 790–822 | `/health` handler |
| 823–857 | `main()` — env config, route registration |

### Security Constraints

- All file access must remain within `/app/logs/` — enforced via `filepath.Clean` prefix check in every handler.
- Session tokens: 16-byte crypto/rand hex. Sessions are in-memory only (reset on restart).
- Cookies: HttpOnly, SameSite=Lax.
- Container mounts logs as read-only (`:ro`).

## Configuration (Environment Variables)

| Variable | Default | Effect |
|---|---|---|
| `PORT` | `8080` | Listen port |
| `BASE_PATH` | `""` | URL prefix for reverse proxy (e.g. `/logvault`). No trailing slash. |
| `AUTH_USER` | `""` | Login username. Leave blank to disable auth entirely. |
| `AUTH_PASSWORD` | `""` | Login password. Leave blank to disable auth entirely. |

## Build & Run

```bash
# Build and start with Docker Compose
docker compose up -d --build

# View logs
docker compose logs -f logvault

# Rebuild after code change
docker compose up -d --build
```

Local dev (no Docker):
```bash
cd app
go run main.go
# Logs must exist at /app/logs/ or adjust the hardcoded path
```

## Development Workflow

### When Adding a Feature

1. All UI changes go inside the template strings in `main.go` (look for `loginTmpl`, `browserTmpl`, `tailTmpl`).
2. New HTTP handlers follow the pattern: validate path → check `filepath.Clean` prefix → serve.
3. Register new routes in `main()` using `http.HandleFunc(p("/route"), requireAuth(handler))` or without `requireAuth` for public endpoints.
4. `p(path)` must wrap every internal URL to respect BASE_PATH.
5. After changes: `docker compose up -d --build` and verify in browser.

### When Debugging

- Check `docker compose logs logvault` for Go panics or HTTP errors.
- SSE tail streams: verify `X-Accel-Buffering: no` header is set (needed for Nginx).
- Auth issues: sessions reset on container restart — clear cookies and re-login.
- Path issues: double-check `filepath.Clean(path)` prefix guard in the relevant handler.

### Commit Convention

Commit after every meaningful fix or feature addition:
```bash
git add <files>
git commit -m "<type>: <short description>"
# types: feat, fix, refactor, docs, chore
```

Always compare `git diff HEAD~1` after committing to verify the change matches intent.

## Constraints & Gotchas

- **No external Go dependencies** — do not add any. All functionality must use stdlib.
- **Single binary** — the app compiles to a single static binary; keep it that way.
- **No persistent state** — sessions are in-memory. Do not introduce a database or file-based session store without discussing it first.
- **BASE_PATH must prefix all URLs** — every internal link or redirect must use `p("/path")`.
- **app_bk/ is not the source** — always edit `app/main.go`, never `app_bk/`.
- The scratch base image has no shell, no package manager, no OS utilities — the binary must be fully static (`CGO_ENABLED=0`).
