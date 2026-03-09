# LogVault ⬡

> Open-source Docker app to browse & download server logs from `/app/logs` via a clean web UI. Supports folder navigation, one-click downloads, and optional session-based auth. Built in Go on a `scratch` base image (~5 MB).

---

## Features

- 📁 Browse nested log directories like a file explorer
- ⬇️ One-click download for any log file
- 🔐 Optional session-based login page (env-driven, no config files)
- 🗂️ Mount multiple services under `/app/logs/<service>`
- 🌐 `BASE_PATH` support for reverse proxy subpath deployments
- 🐳 Final Docker image built on `scratch` — ~5 MB, zero OS overhead
- 🛡️ Path traversal blocked, logs mounted read-only, `HttpOnly` session cookies
- 💚 `/health` endpoint always public for container orchestrators

---

## Architecture

```
golang:alpine  ── build stage (compiles static binary, CGO_ENABLED=0)
     │
     └── scratch ── runtime stage (binary + CA certs only, ~5 MB)
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

All internal links, redirects, breadcrumbs, form actions, and cookie paths are automatically prefixed. Accessing `https://devbot.server247.info/logvault/` will work correctly end-to-end.

---

## Authentication

Authentication is **disabled by default**. Set both env vars to enable the login page:

```yaml
environment:
  AUTH_USER: "admin"
  AUTH_PASSWORD: "your-secret-password"
```

| Behaviour | Detail |
|---|---|
| Auth disabled | All pages accessible without login |
| Auth enabled | Login page shown on first visit |
| Session TTL | 8 hours (in-memory, resets on container restart) |
| `/health` | Always public — never requires auth |

---

## Endpoints

All endpoints are prefixed with `BASE_PATH` when set (e.g. `/logvault/browse/`).

| Endpoint | Description |
|---|---|
| `GET /` | Redirects to `/browse/` |
| `GET /browse/` | Root log directory browser |
| `GET /browse/<path>` | Browse a subdirectory (e.g. `/browse/api`) |
| `GET /download/<path>` | Download a specific log file |
| `GET /health` | JSON health check — always public |
| `GET /login` | Login page (only when auth is enabled) |
| `POST /logout` | Clears session cookie |

---

## Configuration

All configuration is via environment variables — no config files needed.

| Variable | Default | Description |
|---|---|---|
| `BASE_PATH` | _(unset)_ | URL prefix for reverse proxy (e.g. `/logvault`). No trailing slash. |
| `AUTH_USER` | _(unset)_ | Username for login. Auth disabled if blank. |
| `AUTH_PASSWORD` | _(unset)_ | Password for login. Auth disabled if blank. |
| `PORT` | `8080` | Port the server listens on. |

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
|---|---|---|
| Build | `golang:alpine` | ~350 MB |
| **Final** | **`scratch`** | **~5 MB** |

---

## Security Notes

- Compiled with `CGO_ENABLED=0` — fully static, no libc dependency
- Mount logs as `:ro` (read-only) to prevent any writes
- Path traversal (`../../etc/passwd`) is rejected via `filepath.Clean` prefix check
- No shell or OS utilities inside the final image — minimal attack surface
- Session tokens are 16-byte cryptographically random hex strings
- Session cookie is `HttpOnly` and `SameSite=Lax`

---

