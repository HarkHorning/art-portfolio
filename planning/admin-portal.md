# Admin Portal — TUI

A terminal UI for managing art, categories, and orders. Built in Go with Bubble Tea. Runs locally, connects directly to the database. No web server, no deployed service, no attack surface.

---

## Why TUI Over Web Portal

| Web Portal | TUI |
|------------|-----|
| Separate SvelteKit app to maintain | One Go binary, same repo |
| JWT auth system to build | Auth = GCP credentials (already required) |
| Needs deployment or local server | `go run ./cmd/admin` or installed binary |
| Admin API endpoints exposed on backend | No backend changes at all |
| CORS, tokens, sessions | None of the above |

The TUI replaces the entire `/api/admin/*` route group, the JWT middleware, and the separate frontend app. The backend stays exactly as it is.

---

## Architecture

```
Your machine
    │
    ▼
portfolio-admin (Bubble Tea TUI)
    │
    ▼
Cloud SQL Proxy  (encrypted tunnel, requires GCP credentials)
    │
    ▼
Cloud SQL  (same database the public site uses)
```

Direct DB connection. The TUI talks to the database, not the API. Full access, no restrictions, no auth layer to build — GCP IAM handles it.

For local development, connects to the Docker MySQL container directly on port 3307.

---

## Repo Structure

Same repo, new binary alongside the existing server:

```
backend/
├── cmd/
│   ├── server/
│   │   └── main.go        ← existing web server
│   └── admin/
│       └── main.go        ← TUI entry point
├── internal/
│   ├── admin/             ← TUI screens and logic
│   │   ├── app.go         ← root model, screen routing
│   │   ├── art.go         ← art list + edit screens
│   │   └── categories.go  ← category management screen
│   ├── repo/              ← shared with server (same DB layer)
│   └── config/            ← shared with server
```

The TUI reuses the existing `repo` and `config` packages. No duplication.

---

## Tech Stack

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** — TUI framework (Elm architecture in Go)
- **[Bubbles](https://github.com/charmbracelet/bubbles)** — pre-built components: tables, text inputs, spinners
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** — styling and layout
- **Cloud SQL Go Connector** or **Cloud SQL Proxy** — for production DB access
- **sqlx** — already a dependency, used as-is

---

## Screens

```
┌─────────────────────┐
│  Main Menu          │
│  > Art              │
│    Categories       │
│    (Orders)         │
│    Quit             │
└─────────────────────┘

┌─────────────────────────────────────────────────┐
│  Art                                  [N]ew [Q]uit │
│ ─────────────────────────────────────────────────│
│  ID  Title              Medium    Order  Available │
│  1   Woman with Flowers Acrylic   1      ✓        │
│  2   Boat on Lake       Oil       2      ✓        │
│  3   Horse Watercolor   Watercolor 3     ✓        │
│                                                   │
│  [Enter] Edit  [D] Delete  [↑↓] Navigate          │
└─────────────────────────────────────────────────┘

┌──────────────────────────┐
│  Edit Art Piece          │
│                          │
│  Title:        [        ]│
│  Description:  [        ]│
│  Display Order:[        ]│
│  Available:    [✓]       │
│                          │
│  [Enter] Save  [Esc] Back│
└──────────────────────────┘
```

---

## Database Connection

### Local Development
Connects to the Docker MySQL container directly:
```
DB_HOST=127.0.0.1
DB_PORT=3307
DB_USER=root
DB_PASSWORD=devpassword
DB_NAME=portfolio
```

### Production (Cloud SQL)
Two options:

**Option A: Cloud SQL Proxy (simpler)**
```bash
# Start proxy in background
cloud-sql-proxy PROJECT:REGION:INSTANCE --port 3306

# Run admin (connects via TCP to local proxy)
go run ./cmd/admin
```

**Option B: Cloud SQL Go Connector (no proxy binary needed)**
```go
import "cloud.google.com/go/cloudsqlconn"
// Authenticates via Application Default Credentials (gcloud auth)
```

Recommendation: start with Option A (simpler), move to B later.

---

## Running It

```bash
# Local dev (Docker DB)
make admin

# Production (requires Cloud SQL Proxy running)
make admin-prod
```

Makefile targets to add:
```makefile
admin:
    powershell -NoProfile -Command "$$env:DB_PORT='3307'; Set-Location backend; go run ./cmd/admin"

admin-build:
    powershell -NoProfile -Command "Set-Location backend; go build -o bin/portfolio-admin ./cmd/admin"
```

---

## What It Can Do

### Art Management
- List all art pieces in a table
- Add new piece (title, description, category, image URL, display order, portrait flag)
- Edit any field
- Toggle availability
- Reorder display order
- Delete piece

### Category Management
- List all categories
- Add new category (name + slug)
- Delete category

### Future: Orders
- List orders with status
- Update status (paid → shipped → delivered)
- Add tracking number
- View decrypted shipping address

---

## Dependencies to Add

```
github.com/charmbracelet/bubbletea
github.com/charmbracelet/bubbles
github.com/charmbracelet/lipgloss
```

Optionally for production Cloud SQL:
```
cloud.google.com/go/cloudsqlconn
```

---

## Tasks

- [ ] Add Bubble Tea dependencies (`go get`)
- [ ] Scaffold `cmd/admin/main.go` entry point
- [ ] Create `internal/admin/app.go` — root model and screen routing
- [ ] Add admin-specific repo methods if needed (e.g. full art list without limit)
- [ ] Build main menu screen
- [ ] Build art list screen (Bubbles table)
- [ ] Build art edit/create form
- [ ] Build categories screen
- [ ] Add `make admin` and `make admin-build` to Makefile
- [ ] Test against local Docker DB
- [ ] Test against production via Cloud SQL Proxy
