# Development Planning

Future improvements for the portfolio project, organized by priority.

---

## Local Development

**To run locally (no cloud services needed):**
```bash
cd deployment/docker
docker compose up --build
```

- Frontend: http://localhost:3000
- Backend: http://localhost:8080
- MySQL: localhost:3307

**Note:** Images won't load locally because they're stored in cloud storage. This will be addressed in a future update with either local image storage or placeholder images.

---

## Infrastructure Decisions

### Cloud Provider: Google Cloud Only
All production infrastructure uses GCP:
- **Cloud Run** - cheap, serverless containers (~$0-5/mo)
- **GKE Autopilot** - full Kubernetes for interviews (~$30+/mo)
- **Cloud SQL** - shared MySQL database (~$7-10/mo)
- **Artifact Registry** - container images

### Database Architecture
```
Local Development     →  Local MySQL (Docker container)
Cloud Run (prod)      →  Cloud SQL (shared)
GKE Kubernetes (prod) →  Cloud SQL (shared)
```

All production deployments share the same Cloud SQL instance. Data changes are visible across all environments.

---

## High Priority

### 1. Database Migrations
**Status:** Structure created (migration files ready, not yet integrated into app)
**Why:** Current `InitSchema`/`SeedData` runs on every startup. Can't modify schema without losing data.

**Implementation:**
```
deployment/
└── migrations/
    ├── 000001_initial_schema.up.sql
    ├── 000001_initial_schema.down.sql
    ├── 000002_seed_art_data.up.sql
    └── 000002_seed_art_data.down.sql
```

**Tool:** [golang-migrate](https://github.com/golang-migrate/migrate)

**Changes needed:**
1. Create migrations directory
2. Move schema from `schema.go` to SQL files
3. Remove `InitSchema`/`SeedData` from app startup
4. Add `make migrate` and `make migrate-local` commands

---

### 2. Tests
**Status:** Not started

**Backend tests:**
```
backend/
├── internal/
│   ├── api/
│   │   └── handler_test.go
│   └── repo/
│       └── mysql_test.go
```

**Frontend tests:**
- Vitest for unit tests
- Playwright for E2E tests

**CI integration:**
```yaml
- name: Run tests
  run: go test ./...
```

---

### 3. API Versioning
**Status:** Complete

Changed `/api/` to `/api/v1/`:
- `backend/internal/api/router.go`
- `frontend/nginx.conf`
- `frontend/src/lib/components/artGrid/ArtGrid.svelte`

---

### 4. Cloud SQL Connection Handling
**Status:** Complete
**Why:** Cloud Run uses Unix socket (`/cloudsql/...`), GKE uses TCP. App needs to handle both.

**Implementation:**
```go
func buildDSN(cfg Config) string {
    if strings.HasPrefix(cfg.Host, "/cloudsql/") {
        // Cloud Run - Unix socket
        return fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true",
            cfg.User, cfg.Password, cfg.Host, cfg.Database)
    }
    // Standard TCP connection
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}
```

---

## Medium Priority

### 5. Structured Logging
**Status:** Complete

Replace `fmt.Println` with Go's `log/slog`:

```go
slog.Error("database connection failed", "host", cfg.Host, "error", err)
```

Benefits:
- JSON output for Cloud Logging
- Log levels
- Structured fields for filtering

---

### 6. Health Check Endpoints
**Status:** Complete

```go
GET /health  // Is server running?
GET /ready   // Is server ready? (pings database)
```

Cloud Run and GKE use these for traffic routing.

---

### 7. CI/CD for GCP
**Status:** Not started

Replace Azure GitHub Actions with GCP:

```yaml
- name: Auth to Google Cloud
  uses: google-github-actions/auth@v2
  with:
    credentials_json: ${{ secrets.GCP_SA_KEY }}

- name: Deploy to Cloud Run
  run: gcloud run deploy ...
```

---

### 8. Frontend Error Handling
**Status:** Complete

ArtGrid.svelte now handles:
- Loading state
- API errors
- Empty results
```

---

## Lower Priority

### 9. Environment Configuration
**Status:** Complete

Single config struct with environment detection:

```go
type Config struct {
    Environment string // "local", "cloudrun", "k8s"
    Database    DatabaseConfig
    Server      ServerConfig
}

func LoadConfig() Config {
    env := os.Getenv("ENVIRONMENT")
    // Load appropriate defaults based on environment
}
```

---

### 10. API Documentation
**Status:** Not started

Add OpenAPI/Swagger:
```
backend/
└── api/
    └── openapi.yaml
```

---

### 11. Rate Limiting
**Status:** Not started

```go
router.Use(ratelimit.New(100, time.Minute))
```

---

### 12. Admin Portal
**Status:** Not started — see `feature_roadmap.md` Feature 6 for full plan

---

### 13. Page Titles & Meta Tags
**Status:** Not started

Add `<svelte:head>` to each page:
```svelte
<svelte:head>
    <title>Hark Horning — Portfolio</title>
    <meta name="description" content="..." />
</svelte:head>
```
Affects SEO and browser tab labels. Each page needs its own title.

Pages to update: `/`, `/about`, `/art/[id]`

---

### 14. Custom Error Page
**Status:** Not started

SvelteKit uses `+error.svelte` for unmatched routes and runtime errors.
Currently shows a raw default. Should match site style.

```
frontend/src/routes/+error.svelte
```

---

### 15. Open Graph Meta Tags
**Status:** Not started

Makes links look good when shared on social media:
```svelte
<meta property="og:title" content="Hark Horning" />
<meta property="og:image" content="..." />
<meta property="og:description" content="..." />
```

---

### 16. Art Detail Back Button
**Status:** Not started

Currently "← Back" goes to `/` which loses filter state.
Should use `history.back()` so the user returns to their filtered view.

```svelte
<button onclick={() => history.back()}>← Back</button>
```

---

## Completed Items

| Item | Date | Notes |
|------|------|-------|
| API Versioning | 2026-04-03 | `/api/` → `/api/v1/` |
| Remove exposed credentials | 2026-04-03 | Removed terraform.tfvars from git history |
| Kubernetes secrets | 2026-04-03 | Moved passwords to K8s Secrets |
| Navigation bar | 2026-04-05 | Created minimal navbar component |
| About page (placeholder) | 2026-04-05 | Route created, needs content |
| Frontend error handling | 2026-04-07 | Added loading/error states to ArtGrid |
| Health check /ready endpoint | 2026-04-07 | Pings database to verify connectivity |
| Migrations structure | 2026-04-08 | Created folder, initial schema SQL, README |
| Makefile | 2026-04-08 | Local dev, cloudrun, and db-only commands |
| Art Details Page | 2026-04-11 | `/art/[id]` route, backend endpoint, clickable tiles |
| Category Filters | 2026-04-11 | Collapsible sidebar, backend filtering, `/api/v1/categories` |
| Cloudrun env.template | 2026-04-08 | Template for GCP deployment config |
| .gitignore updates | 2026-04-08 | Added cloudrun, cloudflare, credentials entries |
| Cloud SQL connection handling | 2026-04-08 | Unix socket support for Cloud Run |
| Structured logging | 2026-04-08 | Replaced log with slog, JSON for cloud |
| Environment configuration | 2026-04-08 | Unified config with auto-detection |

---

## Migration Path: Azure → GCP

### Files to Update
- [ ] `deployment/terraform/*.tf` - Rewrite for GCP
- [ ] `deployment/kubernetes/*.yaml` - Update image URLs to Artifact Registry
- [ ] `.github/workflows/deploy.yaml` - Replace Azure auth with GCP auth
- [x] Backend config - Support Cloud SQL Unix socket

### Files to Delete
- [ ] Any Azure-specific configs

### New Files to Create
- [ ] `deployment/cloudrun/deploy.sh`
- [ ] `deployment/cloudrun/.env` (from template, gitignored)
- [ ] `Makefile`
- [ ] `deployment/migrations/*.sql`

---

## Cost Projections

| Setup | Monthly Cost |
|-------|--------------|
| Local only | $0 |
| Cloud Run + Cloud SQL | ~$10-15 |
| GKE + Cloud SQL | ~$40-50 |
| Full setup (both options available) | ~$10-15 (GKE off by default) |
