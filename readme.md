# LogVault ⬡

> Open-source Docker app to browse, view, tail, and download server logs from `/app/logs` via a clean web UI. Supports folder navigation, live tail streaming, full file viewer with keyword search, and optional session-based auth. Built in Go on a `scratch` base image (~9.21 MB).

---

## Features

- 📁 Browse nested log directories like a file explorer
- 👁️ View full file content with line numbers and keyword search
- ⬇️ One-click download for any log file
- 📡 Live tail via Server-Sent Events (SSE) with real-time keyword filter
- 🔍 Keyword search with highlight, prev/next navigation, and filter mode
- 🔃 Sortable table columns — click Name, Size, or Modified to sort (asc/desc)
- 🔐 Optional session-based login page (env-driven, no config files)
- 🔒 Login lockout after repeated failed attempts (`MAX_LOGIN_ATTEMPTS` / `LOGIN_LOCKOUT_HOURS`)
- 🗂️ Mount multiple services under `/app/logs/<service>`
- 🌐 `BASE_PATH` support for reverse proxy subpath deployments
- 🖼️ Custom logo/favicon across every page via `LOGO_URL` or a mounted `LOGO_PATH` file
- 🐳 Final Docker image built on `scratch` — ~9.21 MB, zero OS overhead
- 🛡️ Path traversal blocked, logs mounted read-only, `HttpOnly` session cookies
- 💚 `/health` endpoint always public for container orchestrators

---

## Architecture

```
golang:alpine  ── build stage (compiles static binary, CGO_ENABLED=0)
     │
     └── scratch ── runtime stage (binary + CA certs only, ~9.21 MB)
```

---

## Quick Start

### 1. Clone & configure

Edit `docker-compose.yml` to mount your service log directories:

```yaml
volumes:
  ## Format: /host/path/to/logs:/app/logs/<service-name>:ro
  - /app/api/logs:/app/logs/api:ro
  - /app/identity/logs:/app/logs/identity:ro
```

Each mounted path becomes a browsable folder in the UI.

### 2. Run

```bash
docker compose up -d
```

Open [http://localhost:8080](http://localhost:8080)

---

## File Viewer

Click **≡ View** on any log file to open the full file viewer:

- Line numbers with scrollable content
- **Search** input (or press `Ctrl+F`) — highlights all keyword matches in amber
- `↑` / `↓` buttons or `Enter` / `Shift+Enter` to jump between matches
- **Filter** toggle — hides non-matching lines entirely
- Files over 5 MB show only the last 5 MB with a warning banner

---

## Live Tail

Click **⊞ Tail** to open the real-time log tail:

- Streams new lines via Server-Sent Events as they are written
- **Filter** input in the toolbar — matching lines are highlighted, non-matching lines are dimmed
- **↺ Clear** button to reset the stream and start fresh
- Auto-reconnects on disconnect

---

## Reverse Proxy / Subpath Deployment

LogVault supports being served under a URL subpath via the `BASE_PATH` env var.

### Nginx config

```nginx
location /logvault/ {
    proxy_pass         http://logvault:8080/logvault/;
    proxy_set_header   Host              $host;
    proxy_set_header   X-Real-IP         $remote_addr;
    proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
    proxy_set_header   X-Forwarded-Proto $scheme;
}
```

### docker-compose env

```yaml
environment:
  BASE_PATH: "/logvault"
```

All internal links, redirects, breadcrumbs, form actions, and cookie paths are automatically prefixed.

---

## Authentication

Authentication is **disabled by default**. Set both env vars to enable the login page:

```yaml
environment:
  AUTH_USER: "admin"
  AUTH_PASSWORD: "your-secret-password"
```

| Behaviour | Detail |
| --- | --- |
| Auth disabled | All pages accessible without login |
| Auth enabled | Login page shown on first visit |
| Session TTL | 8 hours by default, configurable via `SESSION_TTL_HOURS` (in-memory, resets on container restart) |
| `/health` | Always public — never requires auth |

---

## Endpoints

All endpoints are prefixed with `BASE_PATH` when set (e.g. `/logvault/browse/`).

| Endpoint | Description |
| --- | --- |
| `GET /` | Redirects to `/browse/` |
| `GET /browse/` | Root log directory browser |
| `GET /browse/<path>` | Browse a subdirectory |
| `GET /view/<path>` | View full file content with search |
| `GET /tail/<path>` | Live tail viewer page |
| `GET /tail-stream/<path>` | SSE stream endpoint (used by tail page) |
| `GET /download/<path>` | Download a log file |
| `GET /health` | JSON health check — always public |
| `GET /logo` | Serves the `LOGO_PATH` file — always public |
| `GET /login` | Login page (only when auth is enabled) |
| `POST /logout` | Clears session cookie |

---

## Configuration

All configuration is via environment variables — no config files needed.

| Variable | Default | Description |
| --- | --- | --- |
| `BASE_PATH` | _(unset)_ | URL prefix for reverse proxy (e.g. `/logvault`). No trailing slash. |
| `AUTH_USER` | _(unset)_ | Username for login. Auth disabled if blank. |
| `AUTH_PASSWORD` | _(unset)_ | Password for login. Auth disabled if blank. |
| `PORT` | `8080` | Port the server listens on. |
| `LOGO_URL` | _(unset)_ | Image URL or `data:` URI to use as the logo instead of the default icon. Shown on the login page, navbar, tail, and viewer pages. Ignored if `LOGO_PATH` is also set. |
| `LOGO_PATH` | _(unset)_ | Path to a logo file mounted into the container (e.g. `/app/logo.png`). Served publicly at `/logo`; takes priority over `LOGO_URL`. |
| `SESSION_TTL_HOURS` | `8` | Session lifetime in hours before requiring re-login. |
| `MAX_LOGIN_ATTEMPTS` | `5` | Failed login attempts (global, not per-IP/user) before login is blocked entirely. |
| `LOGIN_LOCKOUT_HOURS` | `24` | How long login stays blocked after hitting `MAX_LOGIN_ATTEMPTS`. Resets on restart. |

---

## Volume Convention

Mount each service's log directory under `/app/logs/<name>`:

```yaml
volumes:
  - /app/api/logs:/app/logs/api:ro
  - /app/identity/logs:/app/logs/identity:ro
  - /var/log/nginx:/app/logs/nginx:ro
```

This makes each service appear as a top-level folder in the browser.

---

## Image Size

| Stage | Base | Approx. Size |
| --- | --- | --- |
| Build | `golang:alpine` | ~350 MB |
| **Final** | **`scratch`** | **~9.21 MB** |

---

## Security Notes

- Compiled with `CGO_ENABLED=0` — fully static, no libc dependency
- Mount logs as `:ro` (read-only) to prevent any writes
- Path traversal (`../../etc/passwd`) is rejected via `filepath.Clean` prefix check
- No shell or OS utilities inside the final image — minimal attack surface
- Session tokens are 16-byte cryptographically random hex strings
- Session cookie is `HttpOnly` and `SameSite=Lax`

---
