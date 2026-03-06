# retrosnack

A Progressive Web App (PWA) ecommerce store for [@retrosnack.shop](https://instagram.com/retrosnack.shop) — a women's thrift clothing store sourcing secondhand brand clothing, accessories, and shoes in good condition for less than retail. Community members can also sell through the store.

---

## Architecture

retrosnack is built as a **modular monolith** — a single Go binary with clean domain-separated packages that mirror microservice boundaries. This gives the simplicity and zero-cost of a single deployment today, with a clear extraction path to true microservices later.

```mermaid
graph TD
    User["User (Browser / Mobile)"]

    subgraph Cloudflare
        CF_CDN["Cloudflare CDN"]
        R2["Cloudflare R2\nObject Storage"]
        Pages["Cloudflare Pages\nSvelteKit PWA"]
    end

    subgraph "Google Cloud"
        CloudRun["Cloud Run\nGo API"]
    end

    subgraph "Go API Modules"
        Auth["auth"]
        Catalog["catalog"]
        Inventory["inventory"]
        Orders["orders"]
        Payments["payments"]
        Instagram["instagram"]
        Media["media"]
    end

    subgraph "Data"
        Neon["Neon\nPostgreSQL"]
    end

    subgraph "External"
        Stripe["Stripe\nCheckout"]
        IGEmbed["Instagram\noEmbed API"]
    end

    subgraph "CI/CD"
        GHA["GitHub Actions"]
    end

    User --> CF_CDN
    CF_CDN --> Pages
    Pages --> CloudRun
    CloudRun --> Auth & Catalog & Inventory & Orders & Payments & Instagram & Media
    Auth & Catalog & Inventory & Orders --> Neon
    Payments --> Stripe
    Instagram --> IGEmbed
    Media --> R2
    GHA --> Pages
    GHA --> CloudRun
```

### Domain Modules

| Module | Responsibility |
|---|---|
| `auth` | JWT-based authentication |
| `catalog` | Products, categories, sizing |
| `inventory` | Stock per item/variant (most items are one-of-a-kind) |
| `orders` | Order lifecycle |
| `payments` | Stripe Checkout integration |
| `instagram` | oEmbed link management per product |
| `media` | Image upload/serve via Cloudflare R2 |

---

## Tech Stack

| Layer | Technology | Rationale | Cost |
|---|---|---|---|
| Frontend | SvelteKit + `vite-plugin-pwa` | Fast, lightweight, installable PWA | Free |
| Frontend hosting | Cloudflare Pages | Global CDN, zero egress fees | Free |
| Backend | Go + chi router | Simple, fast, single binary | Free |
| Backend hosting | Google Cloud Run | 2M req/month free, scales to zero, managed TLS | Free |
| Database | Neon PostgreSQL | Serverless, 0.5 GB free tier | Free |
| Object storage | Cloudflare R2 | 10 GB free, no egress fees | Free |
| Payments | Stripe Checkout | Hosted checkout, webhook fulfillment | $100 credit |
| ORM / queries | sqlc + goose + pgx/v5 | Type-safe SQL, zero runtime overhead | Free |
| CI/CD | GitHub Actions | Lint, test, build, deploy | Free |

---

## Repository Structure

```
retrosnack/
├── apps/
│   └── frontend/                  # SvelteKit PWA
│       ├── src/
│       │   ├── lib/               # Shared utilities and components
│       │   ├── routes/            # SvelteKit file-based routing
│       │   └── app.html
│       ├── static/
│       │   ├── manifest.json      # PWA manifest
│       │   └── icons/             # App icons for installation
│       ├── package.json
│       ├── svelte.config.js
│       └── vite.config.js         # vite-plugin-pwa configuration
│
├── services/
│   └── api/                       # Go modular monolith
│       ├── cmd/
│       │   └── server/
│       │       └── main.go        # Entry point
│       ├── internal/
│       │   ├── auth/              # JWT authentication
│       │   ├── catalog/           # Products, categories, sizing
│       │   ├── inventory/         # Stock tracking
│       │   ├── orders/            # Order lifecycle
│       │   ├── payments/          # Stripe Checkout
│       │   ├── instagram/         # oEmbed link management
│       │   └── media/             # Image upload via R2
│       ├── db/
│       │   ├── migrations/        # goose SQL migrations
│       │   └── queries/           # sqlc SQL query definitions
│       ├── pkg/
│       │   ├── config/            # Environment configuration
│       │   ├── middleware/        # HTTP middleware (auth, logging, CORS)
│       │   └── httputil/          # Shared HTTP helpers
│       ├── go.mod
│       └── go.sum
│
├── infrastructure/
│   ├── nginx/
│   │   └── nginx.conf             # Reverse proxy config
│   ├── docker/
│   │   └── Dockerfile             # Multi-stage Go build
│   └── cloudrun/
│       └── service.yaml           # Cloud Run service definition (Knative)
│
├── docs/
│   └── architecture/
│       └── decisions.md           # Architecture decision records (ADRs)
│
├── .github/
│   └── workflows/
│       ├── ci.yml                 # Lint, test, build on PRs
│       └── deploy.yml             # Deploy to Cloud Run + Cloudflare Pages on main
│
├── docker-compose.yml             # Local development stack (Postgres, API, frontend)
├── sqlc.yaml                      # sqlc code generation config
└── README.md
```

---

## Development Workflow

### Prerequisites

- [Go](https://go.dev/dl/) 1.23+
- [Node.js](https://nodejs.org/) 20+ and [pnpm](https://pnpm.io/)
- [Docker](https://www.docker.com/) and Docker Compose
- [sqlc](https://sqlc.dev/) — `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- [goose](https://github.com/pressly/goose) — `go install github.com/pressly/goose/v3/cmd/goose@latest`

### Local Development

A `Makefile` is provided so you don't need to remember individual commands or pass secrets in the terminal. All env vars are loaded from `.env` automatically.

1. **Clone the repository**

   ```bash
   git clone https://github.com/retrosnack-clothing/retrosnack.git
   cd retrosnack
   ```

2. **Copy environment variables**

   ```bash
   cp .env.example .env
   ```

   The defaults in `.env.example` point at the local Docker Postgres — no changes needed to get started. Fill in Stripe and R2 keys if you need to test payments or image uploads.

3. **Install dependencies and start the database**

   ```bash
   make install       # install Go and frontend deps
   make db            # start local PostgreSQL via Docker
   ```

4. **Run migrations and generate types**

   ```bash
   make migrate       # apply SQL migrations
   make sqlc          # generate type-safe Go code
   ```

5. **Start the services** (two terminals)

   ```bash
   make api           # terminal 1 — Go API on localhost:8080
   make frontend      # terminal 2 — SvelteKit on localhost:5173
   ```

   Or start everything at once via Docker Compose:

   ```bash
   make dev
   ```

Run `make help` to see all available commands.

---

## Deployment Overview

### Cloudflare Pages (Frontend)

- Push to `main` triggers GitHub Actions to build and deploy the SvelteKit app to Cloudflare Pages.
- Global CDN distribution, zero egress fees, automatic HTTPS.

### Google Cloud Run (Backend API)

- Go binary built via multi-stage Docker build in `infrastructure/docker/Dockerfile`.
- Cloud Run handles HTTPS termination natively — no nginx sidecar needed.
- Service definition: `infrastructure/cloudrun/service.yaml`.
- Secrets managed via Google Secret Manager and referenced as env vars in the service YAML.
- Free tier: 2 million requests/month, 360,000 GB-seconds compute — sufficient for early traffic.
- Scales to zero when idle; cold starts are fast (~200 ms) for a static Go binary.
- Deploy via GitHub Actions using Workload Identity Federation (no long-lived service account keys).

### Neon PostgreSQL

- Serverless PostgreSQL, connects from Cloud Run via `DATABASE_URL`.
- Free tier: 0.5 GB storage, auto-suspend when idle.

### Cloudflare R2

- Product images uploaded via Go `media` module using the S3-compatible API.
- Served directly from R2 public bucket URL or via Cloudflare CDN.

### Stripe

- Payments handled via Stripe Checkout (hosted, redirect-based).
- Fulfillment triggered by webhook at `POST /api/webhooks/stripe`.
- Stripe signing secret validated on every webhook event.

---

## Environment Variables

| Variable | Description | Example |
|---|---|---|
| `DATABASE_URL` | Neon PostgreSQL connection string | `postgres://user:pass@host/db?sslmode=require` |
| `JWT_SECRET` | Secret for signing JWT tokens | `random-32-byte-string` |
| `STRIPE_SECRET_KEY` | Stripe API secret key | `sk_live_...` |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook signing secret | `whsec_...` |
| `R2_ACCOUNT_ID` | Cloudflare account ID | `abc123...` |
| `R2_ACCESS_KEY_ID` | R2 S3-compatible access key | `...` |
| `R2_SECRET_ACCESS_KEY` | R2 S3-compatible secret key | `...` |
| `R2_BUCKET_NAME` | R2 bucket name for product images | `retrosnack-media` |
| `R2_PUBLIC_URL` | Public base URL for serving images | `https://media.retrosnack.shop` |
| `PORT` | HTTP port for the Go API | `8080` |
| `ENV` | Environment (`development`/`production`) | `production` |

**GitHub Actions secrets (Cloud Run deployment):**

| Secret | Description |
|---|---|
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | Workload Identity Federation provider resource name |
| `GCP_SERVICE_ACCOUNT` | GCP service account email used for deployment |
| `GCP_PROJECT_ID` | Google Cloud project ID |
| `GCP_REGION` | Cloud Run region (e.g. `us-central1`) |

Copy `.env.example` to `.env` for local development. Never commit `.env` to version control.
