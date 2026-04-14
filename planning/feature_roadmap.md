# Feature Roadmap

---

## Status Overview

| Feature | Priority | Status |
|---------|----------|--------|
| Art gallery + filters | High | **Complete** |
| Prints page + filters | High | **Complete** |
| Art / print detail pages | High | **Complete** |
| About page | Medium | **Complete** |
| Admin portal (HTMX) | High | **Complete** |
| Cloud Run deployment | High | **Complete** |
| CI/CD pipeline | High | Not started |
| Purchase flow (Stripe) | Medium | Not started |
| Image protection + CAPTCHA | Medium | Not started |
| Orders admin UI | Medium | Not started (schema done) |
| Open Graph meta tags | Low | Not started |
| Logging TUI | Low | Not started |
| Tests | Low | Not started |

---

## Admin Portal

**Status:** Complete
**Location:** `/admin` — served directly from the Go backend

### What it does
- Art management: create, edit, archive, assign categories
- Image upload: drag-and-drop to GCS, high/low variant, magic byte validation
- Prints management: create, edit, archive, stock quantity
- Category management: add, delete
- All behind session cookie auth (bcrypt password, HTTP-only cookie)

### What's not yet in the UI (schema is ready)
- Orders view / status updates — orders table exists, no admin page yet
- Bulk stock update

---

## Purchase Flow (Stripe)

**Status:** Not started
**Dependencies:** Stripe account

### Plan
- Stripe Checkout (hosted) — we never touch card numbers
- Webhook `POST /webhook/stripe` creates order row on successful payment
- `orders` table already exists with all needed fields
- Mark print `sold = true` when `quantity_in_stock` hits 0

### What needs building
- `POST /api/v1/checkout/create-session` — creates Stripe checkout session for a print
- `POST /webhook/stripe` — receives payment confirmation, creates order
- Frontend "Buy" button on print detail page
- Order confirmation page `/order/success`
- Admin orders UI (view status, mark shipped)

### Environment variables needed
```
STRIPE_SECRET_KEY=sk_live_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx
```

---

## Image Protection

**Status:** Not started

### Plan
- Low-res images: public (current state — used for grid and display)
- High-res images: require CAPTCHA (Cloudflare Turnstile — free)
- High-res access returns a short-lived signed URL (10 min expiry)
- Original files stored in private GCS path, never directly accessible

### Image tiers
| Tier | Variant | Access |
|------|---------|--------|
| Display | `low` | Public (current) |
| Original | `high` | CAPTCHA + signed URL |

The `images` table already supports this — `high` and `low` variants per piece. The admin upload already separates them.

### What needs building
- Signed URL generation in Go (HMAC, 10min expiry)
- Cloudflare Turnstile verification endpoint
- "View full resolution" button on art detail page
- GCS bucket ACL: `high` path set to private

### Environment variables needed
```
TURNSTILE_SECRET_KEY=xxx
IMAGE_SIGNING_KEY=32-byte-random
```

---

## CI/CD Pipeline

**Status:** Not started

### Plan
GitHub Actions — on push to `main`:
1. Run `go build ./...`
2. Run tests (when they exist)
3. Build and push Docker images to Artifact Registry
4. Deploy to Cloud Run

### What needs building
- `.github/workflows/deploy.yml`
- GCP service account key stored as GitHub secret (`GCP_SA_KEY`)
- Workload Identity Federation (preferred over key files)

---

## Logging TUI

**Status:** Future — not planned yet

A Bubble Tea terminal viewer for access logs. Not a management tool — just for you to browse traffic patterns, see what's popular, spot anomalies.

What it would show:
- Request counts per endpoint
- Popular art pieces
- Recent 404s
- Traffic over time

Requires adding an `access_logs` table and logging middleware first.

---

## Open Questions

1. **Shipping scope** — US only? International?
2. **Refunds** — Handle via Stripe dashboard manually?
3. **Order notifications** — Email when order ships? (Needs email service — SendGrid, Resend, etc.)
4. **Custom domain** — When to point Cloudflare DNS at Cloud Run URLs?
5. **Print editions** — Limited runs tracked by `quantity_in_stock`, or open editions?
