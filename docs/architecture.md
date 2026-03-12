# architecture

## system overview

```mermaid
graph TD
    User["browser / mobile"]

    subgraph "Cloudflare"
        Pages["Cloudflare Pages - SvelteKit PWA"]
        R2["Cloudflare R2 - images"]
        CDN["Cloudflare CDN - DNS + SSL"]
    end

    subgraph "Render"
        API["Go API"]
    end

    subgraph "external"
        Neon["Neon PostgreSQL"]
        Square["Square payments"]
        IG["Instagram oEmbed API"]
    end

    User --> CDN --> Pages
    Pages -- "fetch /api/*" --> API
    API -- "queries" --> Neon
    API -- "charge cards" --> Square
    API -- "fetch embeds" --> IG
    API -- "upload images" --> R2
    Pages -- "load images" --> R2
    Square -- "webhooks" --> API
```

## request flow

```mermaid
sequenceDiagram
    participant B as browser
    participant F as SvelteKit
    participant A as Go API
    participant DB as Neon PostgreSQL
    participant SQ as Square

    B->>F: visit /shop
    F->>A: GET /api/products
    A->>DB: SELECT products
    DB-->>A: rows
    A-->>F: JSON response
    F-->>B: render product grid

    B->>F: add to cart
    Note over F: stored in localStorage

    B->>F: checkout
    F->>A: GET /api/payments/config
    A-->>F: application_id, location_id
    F->>SQ: tokenize card (Web Payments SDK)
    SQ-->>F: card token (source_id)
    F->>A: POST /api/orders
    A->>DB: INSERT order + reserve inventory
    DB-->>A: order
    A-->>F: order_id
    F->>A: POST /api/payments/process
    A->>SQ: Payments.Create(source_id, amount)
    SQ-->>A: payment result
    A->>DB: UPDATE order status = paid, deduct inventory
    A-->>F: payment_id, status
    F-->>B: redirect to confirmation
```

## backend

the api is a single Go binary with domain-separated packages under `internal/`.

### domain modules

| module | responsibility | depends on |
|--------|---------------|------------|
| `auth` | JWT registration and login | - |
| `catalog` | products, categories, variants | - |
| `inventory` | stock tracking per variant (reserve, deduct, release) | - |
| `orders` | order lifecycle management | `inventory` |
| `payments` | Square card charges and webhook processing | `orders` |
| `instagram` | oEmbed link management per product | - |
| `media` | image upload and serving via R2 | - |

### module dependencies

```mermaid
graph LR
    payments --> orders
    orders --> inventory
    catalog --> |images| media
```

all other modules are independent.

### order state machine

```mermaid
stateDiagram-v2
    [*] --> pending: order created
    pending --> paid: payment completed
    pending --> cancelled: admin cancels
    paid --> shipped: admin ships
    shipped --> delivered: admin confirms delivery
    cancelled --> [*]
    delivered --> [*]
```

inventory is reserved on order creation and deducted on payment. if an order is cancelled, reserved inventory is released.

## layering

each module follows the same pattern:

```
handler.go      HTTP layer - routing, input validation, response formatting
    |
service.go      business logic - depends on repository interface
    |
repository.go   data access - SQL queries via pgx
    |
model.go        domain types
```

dependencies flow inward: handler -> service -> repository. services depend on interfaces, not concrete types. all wiring happens in `cmd/server/main.go` via constructor injection.

### middleware chain

requests pass through middleware in this order:

```
request
  -> RequestID        generate unique id, add to context + response header
  -> Logger           structured slog (method, path, status, duration, ip, request_id)
  -> Recoverer        catch panics, return 500
  -> RealIP           extract client ip from X-Forwarded-For
  -> CORS             environment-aware (wildcard in dev, allowlist in prod)
  -> SecureHeaders    X-Content-Type-Options, X-Frame-Options, Referrer-Policy
  -> MaxBodySize      10 MB limit
  -> RequireJSON      415 on non-json POST/PUT/PATCH (except multipart)
  -> handler
```

auth-protected routes add:

```
  -> Auth             validate JWT Bearer token
  -> RequireRole      check role claim (admin, seller)
```

rate-limited routes (auth, payments) add:

```
  -> RateLimit        per-ip token bucket with Retry-After header
```

### shared packages

| package | purpose |
|---------|---------|
| `pkg/config` | loads env vars, validates required ones, fails fast on startup |
| `pkg/middleware` | request ID, CORS, auth, rate limiting, logging, security headers, content-type validation |
| `pkg/httputil` | JSON response helpers, error formatting (5xx errors are sanitized, logged server-side) |

## frontend

SvelteKit 5 with Svelte 5 runes (`$state`, `$derived`, `$effect`, `$props`). deployed as a PWA to Cloudflare Pages via `adapter-cloudflare`.

### key files

| file | purpose |
|------|---------|
| `src/lib/api.ts` | typed API client for all backend endpoints |
| `src/lib/stores/cart.svelte.ts` | cart state with localStorage persistence |
| `src/lib/stores/toast.svelte.ts` | toast notification store |
| `src/lib/components/` | reusable UI (ProductCard, Navbar, Footer, Toast, etc.) |
| `src/routes/` | pages (home, shop, product detail, cart, checkout, confirmation, about) |
| `src/app.css` | Tailwind v4 with custom theme tokens |

### client-side features

- search, category filter, and price sort on shop page (all via `$derived.by()`)
- "you might also like" recommendations on product detail pages
- toast notifications on add-to-cart
- optimistic cart updates (localStorage, no server round-trip)
- PWA installable on mobile home screens

## database schema

```mermaid
erDiagram
    users {
        uuid id PK
        string email UK
        string password_hash
        string role
    }
    categories {
        uuid id PK
        string name
        string slug UK
        uuid parent_id FK
    }
    products {
        uuid id PK
        string title
        text description
        uuid category_id FK
        string brand
        string condition
        int price_cents
        uuid seller_id FK
        string instagram_post_url
    }
    variants {
        uuid id PK
        uuid product_id FK
        string size
        string color
        string sku UK
    }
    inventory {
        uuid id PK
        uuid variant_id UK
        int quantity
        int reserved
    }
    product_images {
        uuid id PK
        uuid product_id FK
        string url
        int position
    }
    orders {
        uuid id PK
        uuid user_id FK
        string status
        int total_cents
        string checkout_session_id
    }
    order_items {
        uuid id PK
        uuid order_id FK
        uuid variant_id FK
        int quantity
        int price_cents
    }
    instagram_links {
        uuid id PK
        uuid product_id UK
        string post_url
        text embed_html
    }

    users ||--o{ products : sells
    users ||--o{ orders : places
    categories ||--o{ products : contains
    products ||--o{ variants : has
    products ||--o{ product_images : has
    products ||--|| instagram_links : links
    variants ||--|| inventory : tracks
    orders ||--o{ order_items : contains
    variants ||--o{ order_items : references
```

## key decisions

| decision | rationale |
|----------|-----------|
| sqlc over ORM | SQL is the source of truth, type-safe Go code generated from queries, zero runtime overhead |
| Square over Stripe | unifies in-person and online sales under one provider, single dashboard for reporting |
| Neon over self-hosted postgres | managed, serverless, free tier with automated backups |
| Cloudflare R2 over S3 | no egress fees, S3-compatible API, same ecosystem as Pages |
| Render over Cloud Run/Fly.io | simplest free deployment - connect repo and go, no GCP/Fly complexity |
| SvelteKit over Next.js | smaller bundle, simpler mental model, built-in PWA support via vite plugin |
| client-side filtering over API | small catalog (<100 items), avoids extra endpoints, instant UX |
| embedded Square card form over redirect | user never leaves the site, better conversion |
