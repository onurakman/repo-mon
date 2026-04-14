# Repo Monitor

Desktop application for monitoring local git repositories. Tracks uncommitted changes, unpushed/unpulled commits, stash state, branch status, merge conflicts, and remote accessibility.

## Tech Stack

- **Backend:** Go, Wails v2, GORM + SQLite
- **Frontend:** Vue 3, TypeScript, Tailwind CSS 4, PrimeVue 4, Pinia, VueUse
- **Git:** [go-git](https://github.com/go-git/go-git) (pure Go, no git CLI dependency)

## Prerequisites

- Go 1.24+
- Node.js 18+
- npm 9+
- GTK 3 and WebKit2GTK 4.1 dev libraries (Linux)

### Linux (Ubuntu/Debian)

```bash
sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.1-dev
```

### Install Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## Getting Started

```bash
# Install dependencies
make install-deps

# Run in dev mode
make dev

# Build production binary
make build
```

The built binary is at `build/bin/repo-mon`.

## Usage

1. **Add Repository** -- Click "+" in the sidebar, browse to a git repo folder
2. **Dashboard** -- View all repos as cards (grid) or rows (list), filter by tags
3. **Tags** -- Create colored tags to categorize repositories
4. **Settings** -- Switch theme (Carbon/Slate Blue/Purple), toggle dark/light mode, adjust poll interval

## Features

- Monitor uncommitted changes, staged files, untracked files
- Track unpushed/unpulled commits per remote
- Detect merge conflicts
- Count stash entries
- Remote accessibility check with graceful timeout (10s)
- Per-repo and global configurable poll interval
- Grid/list view toggle
- Tag-based filtering
- 3 color themes x dark/light mode

## Project Structure

```
repo-mon/
├── app.go                  # Wails App struct + all bindings
├── main.go                 # Entry point
├── internal/
│   ├── models/             # GORM models (Repository, Tag, UserSettings)
│   ├── database/           # SQLite init + auto-migration
│   ├── git/                # Git operations via go-git + tests
│   ├── monitor/            # Polling scheduler + status computation
│   └── service/            # Business logic (CRUD)
├── frontend/
│   └── src/
│       ├── views/          # Dashboard, AddRepo, Tags, Settings
│       ├── components/     # Sidebar, RepoCard, StatusBadge, etc.
│       ├── stores/         # Pinia stores (repo, tag, settings)
│       ├── composables/    # usePolling
│       ├── themes/         # CSS variable themes
│       └── router/         # Vue Router config
└── build/bin/              # Built binary
```

## Testing

```bash
make test        # All Go tests
make test-git    # Git wrapper tests only
```

## Data Storage

- SQLite database stored at `~/.config/repo-mon/repo-mon.db`
- Stores: repositories, tags, user settings
- Repo status is computed at runtime (not persisted)
