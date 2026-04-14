# AGENTS.md

Instructions for AI agents working on this project.

## Project Overview

Repo Monitor is a Wails v2 desktop app (Go backend + Vue 3 frontend) that monitors local git repositories.

## Build & Run

```bash
# Required: Go 1.24+, Node 18+, Wails CLI, GTK3 + WebKit2GTK 4.1 (Linux)
# Go 1.26.2 is at ~/go/go1.26.2/bin — add to PATH before running wails commands

export PATH="$HOME/go/go1.26.2/bin:$HOME/go/bin:$PATH"

make install-deps   # Install Go + npm dependencies
make dev            # Dev mode with hot reload
make build          # Production binary → build/bin/repo-mon
make test           # Run Go tests
```

## Key Constraints

- **webkit2_41 build tag required:** Always pass `-tags webkit2_41` to wails commands (Makefile handles this)
- **CWD must be project root** when running `wails build` or `wails dev`
- **frontend/.npmrc** has `legacy-peer-deps=true` — needed for dependency resolution
- **No git actions in the app** — monitoring only, no pull/push/commit
- **Remote errors are swallowed** — VPN/unreachable remotes show warning, app continues
- **SQLite via CGO** — requires C compiler (gcc) for `go-sqlite3`

## Architecture

- **Monolith backend:** All logic in one Go process, Wails bindings expose to frontend
- **Polling:** Each repo gets its own `time.Ticker` goroutine in `internal/monitor/scheduler.go`
- **Git ops:** Native `git` CLI via `os/exec`, 10s timeout on remote operations (`internal/git/git.go`)
- **DB:** SQLite via GORM, auto-migrated at startup (`internal/database/database.go`)
- **State:** `RepoStatus` is runtime-only (not persisted), stored in `sync.Map`

## Frontend

- Wails generates TypeScript bindings at `frontend/wailsjs/go/` — field names are **camelCase** from Go json tags (e.g., `repo.name`, not `repo.Name`), except `ID` which comes from `gorm.Model`
- Stores import directly from `wailsjs/go/models` types
- Themes use CSS custom properties (`--color-*`), set via classes on `<html>` (`theme-neutral-carbon dark`)
- PrimeIcons for icons (`pi pi-*` classes)

## Testing

- Git wrapper tests in `internal/git/git_test.go` — use temp repos, no mocks
- No frontend tests currently

## File Conventions

- Go: standard Go project layout with `internal/`
- Frontend: Vue 3 Composition API + `<script setup lang="ts">` everywhere
- Styles: Tailwind utility classes + CSS variables for theming, no scoped styles
