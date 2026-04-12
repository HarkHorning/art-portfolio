# Admin Portal

Private web portal for managing art, categories, and orders. Two users (Hark + partner). Zero admin code in the public frontend.

---

## Core Principle

The public site and the admin portal are completely separate applications. They share the same backend database but talk to different API route groups. A visitor to the public site cannot discover, access, or be affected by the admin portal in any way.

---

## Architecture

```
Public Site                        Admin Portal
harkhorning.com                    admin.harkhorning.com (or local only)
        │                                  │
        ▼                                  ▼
  /api/v1/*                         /api/admin/*
  No auth required                  JWT required on every request
        │                                  │
        └──────────────┬────────────────────┘
                       ▼
                  Cloud SQL
```

### Why a Separate App

- Admin code never ships in the public JavaScript bundle
- The admin URL is not linked from the public site — not in HTML, not in JS, not discoverable by crawlers
- If the admin auth is broken or misconfigured, public visitors are completely unaffected
- Can be run locally with no deployment at all

---

## Recommended Deployment Strategy

### Phase 1: Local Only (Start Here)

Run the admin app on your laptop. Point it at the production API. Never deploy it publicly.

```
Your laptop → npm run dev (admin app) → production /api/admin/* → Cloud SQL
```

**Pros:**
- Zero attack surface — nothing is publicly reachable
- No infrastructure to set up
- Works immediately

**Cons:**
- You and your partner need the repo on your machines
- Can't access from phone or other devices

### Phase 2: Private Cloud Run Service (When Needed)

Deploy admin as a second Cloud Run service. The URL is long and unguessable by default (`https://admin-xxxxx-uc.a.run.app`). Don't link it anywhere.

Optionally restrict access to your home/office IP via Cloud Armor — one config rule, takes 10 minutes.

**Recommendation:** Start with Phase 1. Move to Phase 2 when you need access away from your machine.

---

## Authentication

### Approach: Environment Variable Credentials + JWT

No user table in the database. Two hardcoded users stored as environment variables with bcrypt-hashed passwords.

```
ADMIN_USERNAME_1=hark
ADMIN_PASSWORD_HASH_1=$2a$12$...   (bcrypt hash, not the real password)
ADMIN_USERNAME_2=partner
ADMIN_PASSWORD_HASH_2=$2a$12$...
ADMIN_JWT_SECRET=32-byte-random-string
ADMIN_TOKEN_EXPIRY=8h
```

### Login Flow

```
Admin app → POST /api/admin/login { username, password }
         ← { token: "eyJ..." }  (JWT, 8hr expiry)

All subsequent requests →  Authorization: Bearer eyJ...
Backend middleware verifies token on every /api/admin/* route
```

### Why Not a Database User Table

- You have exactly two users. A user table is unnecessary complexity.
- No signup flow = no signup endpoint to attack
- No password reset flow = no reset endpoint to attack
- Rotating credentials = update an env var and redeploy. Simple.

### Why Not Google OAuth

OAuth is a good option long-term but adds external dependency and setup complexity. Start simple.

---

## Backend Changes

### New Route Group

```go
admin := router.Group("/api/admin")
{
    // Login does not require auth middleware
    admin.POST("/login", handle.AdminLogin)

    // Everything else requires a valid JWT
    authorized := admin.Group("/")
    authorized.Use(handle.AdminAuthMiddleware)
    {
        authorized.GET("/art", handle.AdminListArt)
        authorized.POST("/art", handle.AdminCreateArt)
        authorized.PUT("/art/:id", handle.AdminUpdateArt)
        authorized.DELETE("/art/:id", handle.AdminDeleteArt)
        authorized.PUT("/art/:id/order", handle.AdminReorderArt)

        authorized.GET("/categories", handle.AdminListCategories)
        authorized.POST("/categories", handle.AdminCreateCategory)
        authorized.DELETE("/categories/:id", handle.AdminDeleteCategory)
    }
}
```

### JWT Middleware

```go
func (h *Handler) AdminAuthMiddleware(c *gin.Context) {
    token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
    if !validateJWT(token, os.Getenv("ADMIN_JWT_SECRET")) {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    c.Next()
}
```

---

## Admin App Structure

Separate SvelteKit app in `admin/` at the repo root.

```
admin/
├── src/
│   ├── routes/
│   │   ├── +layout.svelte          ← auth guard: redirect to /login if no token
│   │   ├── login/
│   │   │   └── +page.svelte        ← username/password form
│   │   ├── art/
│   │   │   ├── +page.svelte        ← art list with edit/delete
│   │   │   ├── new/
│   │   │   │   └── +page.svelte    ← add new art piece
│   │   │   └── [id]/
│   │   │       └── +page.svelte    ← edit existing piece
│   │   └── categories/
│   │       └── +page.svelte        ← manage categories
│   └── lib/
│       └── api.ts                  ← fetch wrapper that attaches JWT header
├── package.json
└── svelte.config.js
```

### Auth Guard

The layout checks for a stored token and redirects to `/login` if missing or expired. All pages inherit this automatically.

```svelte
<!-- +layout.svelte -->
<script>
    import { browser } from '$app/environment';
    import { goto } from '$app/navigation';

    if (browser) {
        const token = localStorage.getItem('admin_token');
        if (!token) goto('/login');
    }
</script>
```

### API Wrapper

```ts
// lib/api.ts
export async function adminFetch(path: string, options: RequestInit = {}) {
    const token = localStorage.getItem('admin_token');
    return fetch(`/api/admin${path}`, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
            ...options.headers
        }
    });
}
```

---

## What the Portal Can Do

### Art Management
- View all art pieces
- Add a new piece (title, description, medium/category, image URL, display order, portrait flag)
- Edit any field on an existing piece
- Toggle availability (for purchase flow later)
- Delete a piece

### Category Management
- View all categories
- Add a new category (name + slug)
- Delete a category (cascades via foreign key)

### Future: Order Management
- View orders
- Update status (paid → shipped → delivered)
- Add tracking number
- View decrypted shipping address (when purchase flow is built)

---

## What the Portal Does NOT Do

- Image upload/hosting — paste a URL from Cloud Storage
- User management — credentials are in env vars
- Analytics — use Cloud Logging / Google Analytics for that
- Public-facing anything — it is entirely private

---

## Security Checklist

- [ ] Passwords are bcrypt-hashed, never stored in plaintext
- [ ] JWT secret is a random 32-byte string stored in env vars
- [ ] `/api/admin/*` routes are completely separate from `/api/v1/*`
- [ ] Admin app is not linked from the public site
- [ ] Admin app origin is not in the public CORS allowlist
- [ ] Tokens expire after 8 hours
- [ ] (Phase 2) Cloud Armor IP restriction on admin Cloud Run service

---

## Implementation Order

1. Backend JWT auth middleware + `/api/admin/login`
2. Backend admin CRUD endpoints for art
3. Backend admin CRUD endpoints for categories
4. `admin/` SvelteKit app scaffolded
5. Login page
6. Art list + edit forms
7. Category management
8. Test locally against dev backend
9. (Optional) Deploy as private Cloud Run service

---

## Tasks

- [ ] Add `golang-jwt/jwt` dependency to backend
- [ ] Implement `AdminLogin` handler (bcrypt verify + issue JWT)
- [ ] Implement `AdminAuthMiddleware`
- [ ] Add admin art CRUD endpoints
- [ ] Add admin category CRUD endpoints
- [ ] Scaffold `admin/` SvelteKit app
- [ ] Build login page
- [ ] Build art management pages (list, new, edit)
- [ ] Build category management page
- [ ] Add `make admin` command to Makefile (runs admin dev server)
- [ ] Document credential setup in admin README
