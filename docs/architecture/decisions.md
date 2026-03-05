# Architecture Decision Records

## ADR-001: Modular Monolith over Microservices

**Status:** Accepted
**Date:** 2026-03-05

**Context:** retrosnack is an early-stage project with a small team. True microservices add operational
complexity (service discovery, distributed tracing, inter-service auth) without benefit at this scale.

**Decision:** One Go binary with domain-separated packages (`internal/auth`, `internal/catalog`, etc.)
that mirror microservice boundaries. Single Cloud Run deployment.

**Consequences:** Simple deployment and debugging. Can be extracted to true microservices later by
promoting each `internal/` package to its own repo and adding gRPC/HTTP transport.

---

## ADR-002: SvelteKit PWA on Cloudflare Pages

**Status:** Accepted
**Date:** 2026-03-05

**Context:** Mobile-first shopping experience. Customers browse on phones; installability improves
retention. Cloudflare Pages is free with global CDN and no egress fees.

**Decision:** SvelteKit with `vite-plugin-pwa`, deployed to Cloudflare Pages via `adapter-cloudflare`.

**Consequences:** Zero hosting cost. App is installable on Android/iOS home screens. Service worker
enables offline browsing of cached pages.

---

## ADR-003: Neon PostgreSQL

**Status:** Accepted
**Date:** 2026-03-05

**Context:** Need a managed PostgreSQL with a free tier that works well with Fly.io.

**Decision:** Neon serverless PostgreSQL. Free tier: 0.5 GB, auto-suspend when idle.

**Consequences:** Cold-start latency on first connection after idle period. Acceptable for low-traffic
MVP. Uses pgx/v5 connection pool with `pgxpool` to amortize connection cost.

---

## ADR-004: Cloudflare R2 for Media

**Status:** Accepted
**Date:** 2026-03-05

**Context:** Product images need reliable object storage. S3 charges egress; Cloudflare R2 does not.

**Decision:** Cloudflare R2 with S3-compatible API. 10 GB free storage, no egress fees.

**Consequences:** Images served directly from R2 public URL or via Cloudflare CDN. Go `media` module
uses `aws-sdk-go-v2` with a custom R2 endpoint.

---

## ADR-005: Stripe Checkout for Payments

**Status:** Accepted
**Date:** 2026-03-05

**Context:** Need a secure, hosted payment flow. Custom card forms require PCI compliance overhead.

**Decision:** Stripe Checkout (redirect-based). Webhook at `POST /api/webhooks/stripe` fulfills orders.
Stripe signing secret validates every webhook event.

**Consequences:** No PCI scope on our servers. Order fulfillment is event-driven, not synchronous.
$100 Stripe credit available for initial transaction fees.

---

## ADR-006: Instagram oEmbed (No API Auth)

**Status:** Accepted
**Date:** 2026-03-05

**Context:** retrosnack's Instagram (@retrosnack.shop) is the primary product discovery channel.
Each product should link to its Instagram post. Instagram's Graph API requires app review for
advanced access; oEmbed works for public posts with no auth.

**Decision:** Store `instagram_post_url` per product. Render post embed using Instagram oEmbed API
on product pages. Cache `embed_html` in `instagram_links` table. Fall back to direct link if
oEmbed is unavailable.

**Consequences:** Works immediately, no API approval needed. If Instagram deprecates oEmbed for
unauthenticated use, we fall back gracefully to a link.

---

## ADR-007: sqlc for Type-Safe SQL

**Status:** Accepted
**Date:** 2026-03-05

**Context:** ORMs add magic and performance overhead. Raw SQL is verbose and error-prone.

**Decision:** sqlc generates type-safe Go code from SQL queries. goose manages migrations. pgx/v5
is the driver.

**Consequences:** SQL is the source of truth. Schema changes require migration + sqlc regeneration.
No runtime reflection overhead.

---

## ADR-008: Google Cloud Run over Fly.io for API Hosting

**Status:** Accepted
**Date:** 2026-03-05

**Context:** Fly.io was the initial hosting choice but its free tier was removed — all apps now
require a paid plan. A truly free alternative is needed for an early-stage project with low traffic.

**Decision:** Google Cloud Run. Free tier: 2 million requests/month, 360,000 GB-seconds of compute,
180,000 vCPU-seconds — enough for substantial early traffic at zero cost. Scales to zero when idle.
Managed HTTPS, no nginx sidecar required. Secrets stored in Google Secret Manager and injected as
env vars. GitHub Actions deploys via Workload Identity Federation (no long-lived service account keys).

**Consequences:**
- Dockerfile production stage is simpler (just the Go binary, no nginx).
- Cold starts are fast for a statically compiled Go binary (~200 ms).
- Neon PostgreSQL connection must tolerate brief reconnection on cold start — `pgxpool` handles this.
- GCP project setup required: Artifact Registry repo, Secret Manager secrets, WIF configuration.
