# Development Planning

---

## Local Development

**Fully containerized:**
```
make local
```

**Native backend (faster iteration):**
```
make dev
```
Then in a second terminal: `cd frontend && npm run dev`

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Admin portal: http://localhost:8080/admin (admin / devpassword)
- MySQL: localhost:3307

**Images:** Served from `gs://hark-portfolio-images` (GCS, public bucket). Seed data uses GCS URLs.

---

## Infrastructure

### Cloud Provider: Google Cloud
- **Cloud Run** — serverless containers (~$0-5/mo)
- **Cloud SQL** — MySQL 8.0, `db-f1-micro`, 10GB SSD (~$7-10/mo)
- **Artifact Registry** — Docker images
- **GCS** — Image storage (`hark-portfolio-images`)

### Database Architecture
```
Local dev  →  MySQL in Docker (port 3307)
Cloud Run  →  Cloud SQL via Unix socket (/cloudsql/...)
```

### Deployment
```
PowerShell (Windows):
  .\deployment\cloudrun\setup.ps1   # one-time infrastructure setup
  .\deployment\cloudrun\deploy.ps1  # build, push, deploy

Make:
  make cloudrun                     # calls deploy.ps1
```

---

## Completed Items

| Item | Notes |
|------|-------|
| API versioning `/api/v1/` | All routes versioned |
| Remove exposed credentials | Removed from git history |
| Navigation bar | Minimal, links to Work / Prints / About |
| About page | Content written |
| Art details page `/art/[id]` | Full detail with categories, size, price |
| Art category filters | Sidebar, backend-filtered |
| Prints page `/prints` | Size + price filters, same layout as work page |
| Print detail page `/prints/[id]` | Stock status, size, price |
| Size + price filters (artwork) | Work page filter sidebar extended |
| Size + price fields on art_tiles | Nullable, migration 3 |
| Quantity in stock on prints | Migration 3 |
| 4-column responsive grid | 4→3→2→1 at breakpoints |
| Page titles | All routes have `<svelte:head><title>` |
| Custom 404 / error page | `+error.svelte` |
| Back button | `history.back()` preserves filter state |
| Rate limiting | Per-IP, 10 req/s burst 20 |
| Health / ready endpoints | `/health`, `/ready` |
| Structured logging | `log/slog`, JSON in cloud |
| Environment config | Auto-detects local / cloudrun |
| Database migrations | golang-migrate, embedded SQL, 6 migrations |
| GCS image storage | Public bucket, seed data uses GCS URLs |
| Normalized images table | `images` table replaces `url_low`/`url_high` columns |
| Normalized prints schema | `art_tile_id` FK, no duplicate title/url columns |
| Orders table | Full shipping address, phone, Stripe ID placeholder |
| Soft delete | `archived_at` on art_tiles and prints |
| Admin portal | HTMX + Go templates at `/admin` |
| Admin auth | Session cookie, bcrypt, brute-force protected |
| Admin — art CRUD | Create, edit, archive, category assignment |
| Admin — image upload | GCS upload, magic byte validation, HTMX swap |
| Admin — prints CRUD | Create, edit, archive, stock management |
| Admin — categories CRUD | Add, delete |
| GCP infrastructure scripts | `setup.ps1`, `deploy.ps1` (PowerShell) |
| Cloud Run deploy script | Builds images, pushes to Artifact Registry, deploys both services |
| Cloud Run deployment | Live — frontend + backend on Cloud Run, Cloud SQL via Unix socket |
| nginx API proxy | `proxy_ssl_server_name on`, `Host: $proxy_host` required for Cloud Run HTTPS upstream |
| Footer | Auto-updating year |

---

## Remaining Work

### High Priority
| Item | Notes |
|------|-------|
| CI/CD pipeline | GitHub Actions — auto-deploy on push to main |

### Medium Priority
| Item | Notes |
|------|-------|
| Purchase flow (Stripe) | Checkout session, webhook, mark sold |
| Image protection | Low/high res tiers, CAPTCHA for high-res (Cloudflare Turnstile) |
| Open Graph meta tags | Social share previews |
| Admin — orders UI | View/update order status (schema done, no UI yet) |

### Lower Priority
| Item | Notes |
|------|-------|
| Tests | Backend handler tests, repo integration tests |
| Logging TUI | Bubble Tea viewer for access logs (future) |
| API documentation | OpenAPI/Swagger |
| Custom domain | Cloudflare DNS → Cloud Run URLs |

---

## Cost Projections

| Setup | Monthly |
|-------|---------|
| Local only | $0 |
| Cloud Run + Cloud SQL | ~$10-15 |
