# Feature Roadmap

Long-term feature planning for the portfolio project.

---

## Privacy Principles

All features follow these rules:

1. **No PII in database** - If the database leaks, no one gets hurt
2. **Payment processing is external** - Stripe handles all payment data
3. **Minimal data collection** - Only store what's absolutely necessary
4. **Addresses are encrypted** - Shipping addresses encrypted at rest, decrypted only for fulfillment

---

## Feature Overview

| Feature | Priority | Status | Dependencies |
|---------|----------|--------|--------------|
| Navigation Bar | High | **Complete** | None |
| About Me Page | Medium | **Complete** (placeholder) | None |
| Art Details Page | High | **Complete** | None |
| Art Category Filters | High | **Complete** | None |
| Image Protection | High | Not started | None |
| Purchase Flow | Medium | Not started | Art Details, Stripe |
| Admin CLI Tool | Medium | Not started | None |
| Order Tracking | Low | Not started | Purchase Flow |

---

## Feature 1: Art Details Page

**Priority:** High
**Status:** Complete

### Description
Clicking on an art piece opens a dedicated page with full details and purchase option.

### User Flow
```
Art Grid → Click tile → /art/[id] → See details → "Purchase" button
```

### Frontend Routes
```
/art              - Art gallery grid (existing)
/art/[id]         - Art detail page (new)
```

### UI Components
```
ArtDetail.svelte
├── Large image display
├── Title
├── Description
├── Category tags
├── Dimensions
├── Medium (oil, acrylic, digital, etc.)
├── Price
├── Availability status
└── "Purchase" button (if available)
```

### API Endpoints
```
GET /api/v1/art           - List all art (existing)
GET /api/v1/art/:id       - Get single art piece details (new)
```

### Database Changes
Extend `art_tiles` table:
```sql
ALTER TABLE art_tiles ADD COLUMN description TEXT;
ALTER TABLE art_tiles ADD COLUMN dimensions VARCHAR(50);
ALTER TABLE art_tiles ADD COLUMN medium VARCHAR(100);
ALTER TABLE art_tiles ADD COLUMN price_cents INT;
ALTER TABLE art_tiles ADD COLUMN available BOOLEAN DEFAULT true;
ALTER TABLE art_tiles ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
```

### Tasks
- [ ] Create migration for new columns
- [x] Add `GetArtByID` handler
- [x] Create ArtDetail.svelte component
- [x] Add SvelteKit route `/art/[id]`
- [x] Update ArtTile to link to detail page

---

## Feature 2: Art Category Filters

**Priority:** High
**Status:** Complete

### Description
Home page shows a curated list of pieces ordered by `display_order`. Clicking a category filter fetches pieces from the backend — no client-side filtering. The server does the work; only relevant pieces are sent to the frontend.

### Categories
| Name | Slug |
|------|------|
| Oil | oil |
| Acrylic | acrylic |
| Watercolor | watercolor |
| Pencil Drawing | pencil-drawing |
| Mixed | mixed |
| Pastel | pastel |
| Misc | misc |

### User Flow
```
Home page (curated list, no filter active)
    ↓
Click category button (e.g. "Watercolor")
    ↓
GET /api/v1/art?category=watercolor
    ↓
Grid replaces with filtered results
    ↓
Click "All" → back to curated list
```

### API Endpoints
```
GET /api/v1/art                        - Curated home list (ordered by display_order)
GET /api/v1/art?category=watercolor    - Backend-filtered by category slug
GET /api/v1/categories                 - All categories, fetched on page load to build filter UI
```

The handler branches on whether the `category` query param is present. Same endpoint, two SQL queries.

### Database Schema
No changes needed. The schema already supports this:
```
art_tiles ──< art_categories >── categories
```

### Backend Changes
- `sql_resource.go` — add `TilesByCategory(slug string)` query (JOIN through art_categories)
- `sql_resource.go` — add `AllCategories()` query
- `handler.go` — `GetArtTiles` reads optional `?category=` param and calls the right query
- `handler.go` — add `GetCategories` handler
- `router.go` — add `GET /api/v1/categories`
- `art_model.go` — add `Category` model (`id`, `name`, `slug`)

### Frontend Changes
On page load, fetch both in parallel:
- `/api/v1/art` → populate grid
- `/api/v1/categories` → populate filter buttons

On category click, fetch `/api/v1/art?category=slug` and replace grid contents.

### UI Components
```
+page.svelte
├── CategoryFilter.svelte
│   ├── "All" button (active by default)
│   ├── Category buttons (dynamic from /api/v1/categories)
│   └── Active state styling
└── ArtGrid.svelte
    └── Receives tiles as a prop (no internal fetch logic)
```

### Tasks
- [x] Add `TilesByCategory(slug string)` to `sql_resource.go`
- [x] Add `AllCategories()` to `sql_resource.go`
- [x] Update `GetArtTiles` handler to branch on `?category=` param
- [x] Add `GetCategories` handler
- [x] Add `GET /api/v1/categories` route
- [x] Add `Category` model to `art_model.go`
- [x] Refactor `ArtGrid.svelte` to accept tiles as a prop
- [x] Create `FilterSidebar.svelte` (collapsible sidebar)
- [x] Update `+page.svelte` to fetch and manage filter state

---

## Feature 3: Image Protection & Anti-Scraping

**Priority:** High
**Status:** Not started

### Description
Protect artwork from being easily scraped. Show low-res/watermarked images by default. Users must solve CAPTCHA to view high-resolution images.

### Image Tiers

| Tier | Resolution | Watermark | Access |
|------|------------|-----------|--------|
| Thumbnail | 300px | Yes | Public (grid view) |
| Preview | 800px | Yes | Public (detail page) |
| High-res | Full | No | CAPTCHA required |

### User Flow
```
Art Grid (thumbnails)
    ↓
Click → Art Detail Page (preview, watermarked)
    ↓
Click "View Full Resolution"
    ↓
CAPTCHA challenge (Cloudflare Turnstile or hCaptcha)
    ↓
Success → Temporary signed URL (expires in 10 min)
    ↓
View high-res image
```

### Anti-Scraping Measures

#### 1. Image Processing Pipeline
```
Original upload
    ↓
Generate thumbnail (300px, watermarked)
    ↓
Generate preview (800px, watermarked)
    ↓
Store original (never directly accessible)
```

#### 2. Signed URLs for High-Res
High-res images are never served directly. Generate temporary signed URLs:

```go
func generateSignedURL(imageID string, expiry time.Duration) string {
    expires := time.Now().Add(expiry).Unix()
    signature := hmac.Sign(fmt.Sprintf("%s:%d", imageID, expires), secretKey)
    return fmt.Sprintf("/api/v1/images/%s/full?expires=%d&sig=%s", imageID, expires, signature)
}
```

#### 3. CAPTCHA Integration (Cloudflare Turnstile - Free)
```go
// Backend: Verify CAPTCHA token
func verifyCaptcha(token string) bool {
    resp, _ := http.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify",
        url.Values{
            "secret":   {os.Getenv("TURNSTILE_SECRET_KEY")},
            "response": {token},
        })
    var result struct { Success bool }
    json.NewDecoder(resp.Body).Decode(&result)
    return result.Success
}
```

```svelte
<!-- Frontend: Turnstile widget -->
<script src="https://challenges.cloudflare.com/turnstile/v0/api.js" async defer></script>

<div class="cf-turnstile"
     data-sitekey="YOUR_SITE_KEY"
     data-callback="onCaptchaSuccess">
</div>
```

#### 4. Rate Limiting
```go
// Per-IP rate limits
var limiter = rate.NewLimiter(rate.Every(time.Second), 10) // 10 requests/second

func rateLimitMiddleware(c *gin.Context) {
    if !limiter.Allow() {
        c.AbortWithStatus(429)
        return
    }
    c.Next()
}
```

#### 5. Additional Protections
- Disable right-click context menu (minor deterrent)
- CSS overlay on images (prevents simple drag-and-drop)
- Referrer checking (reject hotlinking)
- User-Agent filtering (block known bots)

### API Endpoints
```
GET /api/v1/art/:id/thumbnail    - Public, watermarked
GET /api/v1/art/:id/preview      - Public, watermarked
POST /api/v1/art/:id/request-hires
    Request:  { captcha_token: "xxx" }
    Response: { url: "/api/v1/images/xxx?expires=123&sig=abc", expires_in: 600 }
GET /api/v1/images/:id           - Requires valid signature
```

### Database Changes
```sql
ALTER TABLE art_tiles ADD COLUMN thumbnail_path VARCHAR(255);
ALTER TABLE art_tiles ADD COLUMN preview_path VARCHAR(255);
ALTER TABLE art_tiles ADD COLUMN original_path VARCHAR(255);  -- Never exposed directly
```

### Storage Structure (Cloud Storage)
```
images/
├── thumbnails/
│   └── {id}_thumb.jpg      (public)
├── previews/
│   └── {id}_preview.jpg    (public)
└── originals/
    └── {id}_original.jpg   (private, signed URLs only)
```

### Environment Variables
```
TURNSTILE_SITE_KEY=xxx
TURNSTILE_SECRET_KEY=xxx
IMAGE_SIGNING_KEY=32-byte-random-key
```

### Why Cloudflare Turnstile?
- Free (unlike reCAPTCHA enterprise)
- Privacy-focused (no tracking)
- Often invisible (smart challenge)
- Already using Cloudflare for DNS

### Tasks
- [ ] Set up image processing pipeline (resize, watermark)
- [ ] Create Cloud Storage buckets with proper permissions
- [ ] Implement signed URL generation
- [ ] Integrate Cloudflare Turnstile
- [ ] Add rate limiting middleware
- [ ] Create "View Full Resolution" flow in frontend
- [ ] Add referrer/hotlink protection
- [ ] Test with various scraping tools

### Limitations (Accepted)
- Determined scrapers can still screenshot
- This is about raising the bar, not perfection
- Goal: Make bulk scraping impractical

---

## Feature 4: About Me Page

**Priority:** Medium
**Status:** Complete (placeholder content)

### Description
Static page with bio, artist statement, and contact info.

### Routes
```
/about            - About page
```

### What's Done
- [x] Created `frontend/src/routes/about/+page.svelte`
- [x] Added route to SvelteKit
- [x] Created navigation bar with link to About page
- [ ] Write actual content (currently placeholder)

### Content to Add Later
```
AboutPage.svelte
├── Profile photo
├── Bio paragraph
├── Artist statement
├── Skills/mediums list
├── Contact links (email, social)
└── Optional: Timeline/journey
```

### Files Created
- `frontend/src/routes/about/+page.svelte` - About page
- `frontend/src/lib/components/navbar/Navbar.svelte` - Navigation bar
- Updated `frontend/src/routes/+layout.svelte` - Added navbar to all pages

---

## Feature 5: Purchase Flow

**Priority:** Medium
**Status:** Not started
**Dependencies:** Art Details Page, Stripe Account

### Description
Allow users to purchase art. Payment handled entirely by Stripe - we never see card numbers.

### Privacy Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    WHAT WE STORE                                 │
├─────────────────────────────────────────────────────────────────┤
│  orders table:                                                   │
│    - order_id (our reference)                                   │
│    - stripe_payment_intent_id (for lookup, not PII)             │
│    - art_id                                                      │
│    - amount_cents                                                │
│    - status (pending, paid, shipped, delivered)                 │
│    - shipping_address_encrypted (AES-256, see below)            │
│    - created_at                                                  │
│                                                                  │
│  NO: names, emails, phone numbers, card numbers                 │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    WHAT STRIPE STORES                            │
├─────────────────────────────────────────────────────────────────┤
│  - Card numbers (tokenized)                                     │
│  - Billing address                                               │
│  - Email (for receipts)                                          │
│  - Customer name                                                 │
│                                                                  │
│  We query Stripe when needed, don't duplicate                   │
└─────────────────────────────────────────────────────────────────┘
```

### Address Encryption

Shipping address is the only PII we must store. Encrypt it:

```go
// Encrypt before storing
encrypted := encrypt(address, os.Getenv("ADDRESS_ENCRYPTION_KEY"))
db.Exec("INSERT INTO orders (shipping_address_encrypted) VALUES (?)", encrypted)

// Decrypt only when fulfilling order
decrypted := decrypt(order.ShippingAddressEncrypted, key)
```

If database leaks:
- Attacker sees encrypted blob, useless without key
- Key is in environment variable, not in database

### User Flow
```
Art Detail Page
    ↓
Click "Purchase" ($X)
    ↓
Stripe Checkout (hosted by Stripe)
    ↓
Enter shipping address (in Stripe)
    ↓
Enter payment (in Stripe)
    ↓
Stripe redirects to /order/success?session_id=xxx
    ↓
We create order record, mark art as sold
    ↓
Show confirmation
```

### API Endpoints
```
POST /api/v1/checkout/create-session
  Request:  { art_id: 123 }
  Response: { checkout_url: "https://checkout.stripe.com/..." }

GET /api/v1/checkout/success?session_id=xxx
  - Webhook alternative: verify session, create order

POST /api/v1/webhooks/stripe
  - Stripe sends payment confirmation
  - We create/update order record
```

### Database Schema
```sql
CREATE TABLE orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_number VARCHAR(20) UNIQUE NOT NULL,  -- e.g., "ORD-2024-001"
    art_id INT NOT NULL REFERENCES art_tiles(id),
    stripe_payment_intent_id VARCHAR(100),
    stripe_checkout_session_id VARCHAR(100),
    amount_cents INT NOT NULL,
    status ENUM('pending', 'paid', 'shipped', 'delivered', 'cancelled') DEFAULT 'pending',
    shipping_address_encrypted TEXT,  -- AES-256 encrypted JSON
    tracking_number VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Stripe Integration

```go
// Create checkout session
params := &stripe.CheckoutSessionParams{
    PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
    LineItems: []*stripe.CheckoutSessionLineItemParams{{
        PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
            Currency: stripe.String("usd"),
            ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
                Name: stripe.String(art.Title),
            },
            UnitAmount: stripe.Int64(int64(art.PriceCents)),
        },
        Quantity: stripe.Int64(1),
    }},
    Mode: stripe.String("payment"),
    SuccessURL: stripe.String("https://yourdomain.com/order/success?session_id={CHECKOUT_SESSION_ID}"),
    CancelURL: stripe.String("https://yourdomain.com/art/" + artID),
    ShippingAddressCollection: &stripe.CheckoutSessionShippingAddressCollectionParams{
        AllowedCountries: stripe.StringSlice([]string{"US"}),  // Expand as needed
    },
}
session, _ := session.New(params)
```

### Environment Variables
```
STRIPE_SECRET_KEY=sk_live_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx
ADDRESS_ENCRYPTION_KEY=32-byte-random-key
```

### Tasks
- [ ] Create Stripe account
- [ ] Add stripe-go dependency
- [ ] Create orders migration
- [ ] Implement address encryption/decryption
- [ ] Create checkout session endpoint
- [ ] Create Stripe webhook handler
- [ ] Create success page
- [ ] Mark art as unavailable after purchase
- [ ] Test with Stripe test mode

---

## Feature 6: Admin Portal

**Priority:** Medium
**Status:** Not started

### Description
A private web portal for managing art, orders, and content — usable by Hark and one partner. Zero admin code lives in the public frontend, keeping the public attack surface minimal.

### Architecture: Separate Service

```
Public Site                    Admin Portal
harkhorning.com                admin.harkhorning.com (or local only)
     │                                │
     ▼                                ▼
Public Backend API             Protected Backend API
/api/v1/*  (no auth)          /api/admin/* (JWT required)
     │                                │
     └──────────────┬─────────────────┘
                    ▼
               Cloud SQL
```

The admin portal is a **completely separate SvelteKit app** in `admin/` at the repo root. It talks to protected `/api/admin/*` routes on the same backend. The public frontend has no admin code, no admin routes, and no auth logic whatsoever.

### Why Separate Service

| Concern | Approach |
|---------|----------|
| Attack surface | Admin code never ships in public bundle |
| Discoverability | Admin URL is not linked from public site |
| Auth failure | Only affects admin, not public visitors |
| Deployment | Can run locally only — no public exposure needed |

### Authentication

Two admin users (you and your partner). No database user table needed — credentials stored as environment variables:

```
ADMIN_USERNAME_1=hark
ADMIN_PASSWORD_HASH_1=bcrypt_hash
ADMIN_USERNAME_2=partner_name
ADMIN_PASSWORD_HASH_2=bcrypt_hash
```

Login returns a short-lived JWT. All `/api/admin/*` routes verify the JWT via middleware. Tokens expire after 8 hours.

### Deployment Options

**Option A: Local only (lowest risk)**
Run the admin app locally, point it at the production API. Never deployed publicly. Only accessible from your machine.

**Option B: Private Cloud Run service**
Deploy to a separate Cloud Run URL. Not linked anywhere public. Optionally restrict to specific IP addresses via Cloud Armor.

Start with Option A. Promote to B when you need remote access.

### Admin Capabilities

**Art Management**
- Add new art pieces (title, description, category, image URL, display order)
- Edit existing pieces
- Toggle availability
- Reorder display

**Category Management**
- Add / rename / remove categories
- Assign categories to pieces

**Order Management** (after purchase flow is built)
- View orders
- Mark as shipped with tracking number
- View decrypted shipping address

### Backend Changes

New route group with JWT middleware:
```go
admin := router.Group("/api/admin")
admin.Use(handle.AdminAuthMiddleware)
{
    admin.POST("/login", handle.AdminLogin)
    admin.GET("/art", handle.AdminListArt)
    admin.POST("/art", handle.AdminCreateArt)
    admin.PUT("/art/:id", handle.AdminUpdateArt)
    admin.DELETE("/art/:id", handle.AdminDeleteArt)
    admin.GET("/categories", handle.AdminListCategories)
    admin.POST("/categories", handle.AdminCreateCategory)
    admin.DELETE("/categories/:id", handle.AdminDeleteCategory)
}
```

### Frontend Structure

```
admin/                        ← separate SvelteKit app
├── src/
│   ├── routes/
│   │   ├── +layout.svelte   ← checks auth, redirects to /login if missing
│   │   ├── login/
│   │   │   └── +page.svelte
│   │   ├── art/
│   │   │   ├── +page.svelte       ← art list
│   │   │   └── [id]/+page.svelte  ← edit form
│   │   └── categories/
│   │       └── +page.svelte
├── package.json
└── svelte.config.js
```

### Environment Variables

```
ADMIN_USERNAME_1=xxx
ADMIN_PASSWORD_HASH_1=xxx  (bcrypt)
ADMIN_USERNAME_2=xxx
ADMIN_PASSWORD_HASH_2=xxx  (bcrypt)
ADMIN_JWT_SECRET=32-byte-random-key
ADMIN_TOKEN_EXPIRY=8h
```

### Tasks
- [ ] Add JWT middleware to backend
- [ ] Add `/api/admin/login` endpoint
- [ ] Add admin art CRUD endpoints
- [ ] Add admin category endpoints
- [ ] Create `admin/` SvelteKit app
- [ ] Build login page
- [ ] Build art management UI
- [ ] Build category management UI
- [ ] Add layout auth guard (redirect to login if no token)
- [ ] Test locally against dev backend

### Security Model

```
┌────────────────────────────────────────────────────���────────────┐
│                    ATTACK SURFACE: NEAR ZERO                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  To use admin CLI, attacker needs:                              │
│    1. Physical access to your machine (or SSH)                  │
│    2. GCP credentials with Cloud SQL access                     │
│    3. Knowledge that the CLI exists                             │
│                                                                  │
│  No web endpoint to attack                                       │
│  No API to brute force                                          │
│  No session tokens to steal                                      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Architecture

```
Your laptop
    ↓
portfolio-admin CLI (Go binary)
    ↓
Cloud SQL Proxy (encrypted tunnel)
    ↓
Cloud SQL (requires IAM auth)
```

### CLI Commands

```bash
# Art management
portfolio-admin art list
portfolio-admin art add --title "Sunset" --price 500 --category landscapes
portfolio-admin art update --id 5 --available false
portfolio-admin art delete --id 5

# Order management
portfolio-admin orders list
portfolio-admin orders list --status paid
portfolio-admin orders view --id 123
portfolio-admin orders update --id 123 --status shipped --tracking "1Z999..."
portfolio-admin orders decrypt-address --id 123  # Shows shipping address

# Category management
portfolio-admin categories list
portfolio-admin categories add --name "Landscapes" --slug landscapes
portfolio-admin categories delete --slug landscapes

# User activity logs
portfolio-admin logs list --limit 100
portfolio-admin logs list --date 2024-03-15
portfolio-admin logs search --path "/api/v1/art"
portfolio-admin logs stats  # Request counts, popular pages, etc.

# Database utilities
portfolio-admin db migrate
portfolio-admin db backup
portfolio-admin db stats
```

### User Activity Logs

**What we log (non-sensitive):**
```sql
CREATE TABLE access_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    method VARCHAR(10),           -- GET, POST, etc.
    path VARCHAR(255),            -- /api/v1/art/5
    status_code INT,              -- 200, 404, etc.
    response_time_ms INT,
    ip_hash VARCHAR(64),          -- SHA256 of IP (not the actual IP)
    user_agent_hash VARCHAR(64),  -- SHA256 of UA (not actual UA)
    country VARCHAR(2),           -- From Cloudflare header, if available
    referer_domain VARCHAR(100),  -- Just domain, not full URL
    art_id INT,                   -- If request was for specific art
    INDEX idx_timestamp (timestamp),
    INDEX idx_path (path)
);
```

**What we DON'T log:**
- Full IP addresses (hashed only - can detect same visitor, can't identify them)
- Full user agents (hashed only)
- Request bodies
- Any PII

**Why hash IPs?**
- Can still detect: "Same visitor viewed 5 pieces"
- Can still detect: "Unusual traffic from one source"
- Cannot: "This is John Smith from 123 Main St"

### Implementation

#### Directory Structure
```
cmd/
├── server/
│   └── main.go          # Existing web server
└── admin/
    └── main.go          # Admin CLI

internal/
├── admin/
│   ├── art.go           # Art CRUD commands
│   ├── orders.go        # Order commands
│   ├── logs.go          # Log viewing
│   └── db.go            # DB utilities
```

#### CLI Framework
Use [Cobra](https://github.com/spf13/cobra) for CLI structure:

```go
// cmd/admin/main.go
package main

import (
    "github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "portfolio-admin",
        Short: "Admin CLI for portfolio management",
    }

    rootCmd.AddCommand(artCmd)
    rootCmd.AddCommand(ordersCmd)
    rootCmd.AddCommand(logsCmd)
    rootCmd.AddCommand(dbCmd)

    rootCmd.Execute()
}
```

#### Connection via Cloud SQL Proxy
```go
// Requires Cloud SQL Proxy running locally
// Or use Cloud SQL Go connector

import "cloud.google.com/go/cloudsqlconn"

func connectDB() (*sql.DB, error) {
    d, err := cloudsqlconn.NewDialer(context.Background())
    if err != nil {
        return nil, err
    }

    mysql.RegisterDialContext("cloudsql",
        func(ctx context.Context, addr string) (net.Conn, error) {
            return d.Dial(ctx, "project:region:instance")
        })

    return sql.Open("mysql", "user:password@cloudsql(project:region:instance)/dbname")
}
```

### Authentication Flow

```
1. User runs: portfolio-admin orders list

2. CLI checks for GCP credentials:
   - Application Default Credentials (gcloud auth)
   - Service account key file
   - Fails if neither available

3. CLI connects via Cloud SQL Proxy/Connector
   - GCP validates IAM permissions
   - Connection is encrypted

4. Command executes against database

5. Results displayed in terminal
```

### Building & Distribution

```makefile
# Add to Makefile
admin-build:
	go build -o bin/portfolio-admin ./cmd/admin

admin-install:
	go install ./cmd/admin
```

Only you have the binary. Not distributed, not in any registry.

### Environment Setup (One-time)

```bash
# 1. Install Cloud SQL Proxy
gcloud components install cloud_sql_proxy

# 2. Authenticate
gcloud auth application-default login

# 3. Start proxy (in background)
cloud_sql_proxy -instances=PROJECT:REGION:INSTANCE=tcp:3306 &

# 4. Use CLI
portfolio-admin art list
```

Or use the Go Cloud SQL Connector (no proxy needed):
```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
portfolio-admin art list
```

### Tasks
- [ ] Set up Cobra CLI framework
- [ ] Implement Cloud SQL connection (via connector)
- [ ] Add art CRUD commands
- [ ] Add order management commands
- [ ] Create access_logs table
- [ ] Add logging middleware to API
- [ ] Implement log viewing commands
- [ ] Add database backup command
- [ ] Build and test locally

### Why Not a Web Admin Panel?

| Web Admin | CLI Tool |
|-----------|----------|
| Needs authentication system | Uses GCP IAM |
| Exposed to internet | Not exposed |
| Can be brute forced | No endpoint to attack |
| Session management | No sessions |
| CSRF, XSS risks | No browser = no web vulns |
| Needs HTTPS setup | Cloud SQL Proxy encrypts |

---

## Database Leak Scenario

If someone gets full database access, they see:

| Table | What they get | Risk |
|-------|---------------|------|
| art_tiles | Art metadata, prices | None - public info |
| categories | Category names | None - public info |
| orders | Order IDs, amounts, encrypted addresses | Low - can't decrypt without key |
| access_logs | Hashed IPs, paths, timestamps | None - can't identify anyone |

**They cannot obtain:**
- Customer names
- Email addresses
- Phone numbers
- Payment info (in Stripe)
- Decrypted addresses (key not in DB)
- Real IP addresses (only hashes)
- Real user agents (only hashes)

---

## Implementation Order

### Phase 1: Core Features (No Payment)
1. Art Details Page
2. Category Filters
3. Image Protection & Anti-Scraping
4. About Me Page

### Phase 2: Admin & Infrastructure
5. Admin CLI Tool (manage art, view logs)
6. Access logging middleware

### Phase 3: E-commerce
7. Stripe Integration
8. Purchase Flow
9. Order Tracking (via Admin CLI)

---

## Tech Stack Additions

| Feature | Frontend | Backend | External |
|---------|----------|---------|----------|
| Art Details | SvelteKit routing | New endpoint | - |
| Categories | Filter component | Query params | - |
| Image Protection | CAPTCHA widget, signed URL handling | Image processing, rate limiting, signed URLs | Cloudflare Turnstile |
| About | Static page | - | - |
| Admin CLI | None (CLI only) | Cobra CLI, Cloud SQL Connector | GCP IAM |
| Purchases | Redirect to Stripe | Stripe SDK | Stripe |
| Orders | None (via CLI) | CLI commands | - |

---

## Migration Plan

```sql
-- Migration 001: Art details columns
ALTER TABLE art_tiles ADD COLUMN description TEXT;
ALTER TABLE art_tiles ADD COLUMN dimensions VARCHAR(50);
ALTER TABLE art_tiles ADD COLUMN medium VARCHAR(100);
ALTER TABLE art_tiles ADD COLUMN price_cents INT;
ALTER TABLE art_tiles ADD COLUMN available BOOLEAN DEFAULT true;

-- Migration 002: Categories
CREATE TABLE categories (...);
CREATE TABLE art_categories (...);

-- Migration 003: Image paths (for protection tiers)
ALTER TABLE art_tiles ADD COLUMN thumbnail_path VARCHAR(255);
ALTER TABLE art_tiles ADD COLUMN preview_path VARCHAR(255);
ALTER TABLE art_tiles ADD COLUMN original_path VARCHAR(255);

-- Migration 004: Access logs (privacy-preserving)
CREATE TABLE access_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    method VARCHAR(10),
    path VARCHAR(255),
    status_code INT,
    response_time_ms INT,
    ip_hash VARCHAR(64),          -- SHA256, not actual IP
    user_agent_hash VARCHAR(64),  -- SHA256, not actual UA
    country VARCHAR(2),
    referer_domain VARCHAR(100),
    art_id INT,
    INDEX idx_timestamp (timestamp),
    INDEX idx_path (path)
);

-- Migration 005: Orders
CREATE TABLE orders (...);
```

---

## Open Questions

1. **Shipping scope** - US only initially? International later?
2. **Prints vs originals** - Same purchase flow? Different pricing?
3. **Inventory** - One original per piece? Or editions?
4. **Refunds** - Handle via Stripe dashboard manually?
5. **Notifications** - Email when order ships? (Would need email service)

---

## Future Ideas (Not Planned)

- Newsletter signup
- Commission request form
- Customer accounts (for order history)
- Print-on-demand integration
- Gallery exhibition calendar
