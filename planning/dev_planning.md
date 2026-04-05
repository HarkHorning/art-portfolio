# Development Planning

Future improvements for the portfolio project, organized by priority.

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
**Status:** Not started
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
**Status:** Not started
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
**Status:** Not started

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
**Status:** Partial (has `/health`, needs `/ready`)

```go
GET /health  // Is server running? (already exists)
GET /ready   // Is server ready? (DB connected, migrations run)
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
**Status:** Not started

```typescript
try {
  const res = await fetch(`/api/v1/art`);
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  tiles = await res.json();
} catch (e) {
  error = "Failed to load art";
}
```

---

## Lower Priority

### 9. Environment Configuration
**Status:** Not started

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

### 12. Admin Interface
**Status:** Not started

For adding/editing art without redeploying:
- Simple admin routes with auth
- Or separate admin frontend

---

## Completed Items

| Item | Date | Notes |
|------|------|-------|
| API Versioning | 2026-04-03 | `/api/` → `/api/v1/` |
| Remove exposed credentials | 2026-04-03 | Removed terraform.tfvars from git history |
| Kubernetes secrets | 2026-04-03 | Moved passwords to K8s Secrets |

---

## Migration Path: Azure → GCP

### Files to Update
- [ ] `deployment/terraform/*.tf` - Rewrite for GCP
- [ ] `deployment/kubernetes/*.yaml` - Update image URLs to Artifact Registry
- [ ] `.github/workflows/deploy.yaml` - Replace Azure auth with GCP auth
- [ ] Backend config - Support Cloud SQL Unix socket

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
