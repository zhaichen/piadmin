# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**piadmin** — A lightweight Raspberry Pi management web service. Go backend + Vue 3 frontend compiled into a single binary (~13MB). Designed to run on Raspberry Pi with minimal memory footprint (<30MB RSS).

## Build & Run Commands

```bash
# Development
make dev-backend          # Run Go backend (port 8080)
make dev-frontend         # Run Vite dev server (port 5173, proxies API to 8080)

# Build
make all                  # Build frontend + native binary
make build-arm64          # Cross-compile for Raspberry Pi 64-bit
make build-arm7           # Cross-compile for Raspberry Pi 32-bit

# Test
make test                 # Run Go tests with race detector
```

## Architecture

- **Backend**: Go with chi router, gopsutil for system metrics, gorilla/websocket for real-time push, creack/pty for terminal
- **Frontend**: Vue 3 + TypeScript + Tailwind CSS v4 + ECharts + xterm.js, embedded via `go:embed`
- **Auth**: Single password + JWT token, suitable for LAN use
- **Config**: YAML file + environment variables (`PIADMIN_ADDR`, `PIADMIN_PASSWORD`, `PIADMIN_AUTH_ENABLED`, `PIADMIN_MONITOR_INTERVAL`)

### Key backend packages

- `internal/monitor/` — System metrics collection (CPU/memory/disk/network/temperature) with pub-sub via `Collector`
- `internal/api/` — REST + WebSocket handlers. `router.go` registers all routes and SPA fallback
- `internal/auth/` — JWT auth with middleware (Bearer header + WebSocket query param)
- `internal/config/` — Config loading: YAML → env var → defaults
- `internal/process/` — Process listing and kill via gopsutil
- `internal/services/` — systemd service management (Linux only)
- `internal/network/` — Network interface information
- `internal/terminal/` — WebSocket-based PTY terminal (creack/pty)
- `internal/filemanager/` — File browse/upload/download/delete with path safety
- `internal/gpio/` — GPIO control via sysfs (Raspberry Pi Linux only)
- `web/embed.go` — Embeds frontend dist into Go binary

### API

- `POST /api/auth/login` — Login with password, returns JWT
- `GET /api/system/snapshot` — Current system metrics
- `GET /api/ws/monitor?token=xxx` — WebSocket real-time system metrics
- `GET /api/processes` — Process list
- `DELETE /api/processes?pid=&force=` — Kill process
- `GET /api/services` — systemd services list
- `GET/POST /api/services/{name}` — Service status/action
- `GET /api/network/interfaces` — Network interfaces
- `GET /api/ws/terminal?token=xxx` — WebSocket terminal
- `GET/POST/DELETE /api/files` — File operations
- `GET/POST /api/gpio/*` — GPIO control

### Frontend structure

- `views/Layout.vue` — Sidebar navigation layout (wraps all authenticated pages)
- `views/Dashboard.vue` — Real-time system monitoring with WebSocket
- `composables/useWebSocket.ts` — Auto-reconnecting WebSocket hook
- `api/client.ts` — HTTP/WS client with JWT auth
