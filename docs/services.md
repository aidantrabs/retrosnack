# services and infrastructure

every external service used by retrosnack, what it does, and its free tier limits.

## hosting

### Cloudflare Pages (frontend)

- https://pages.cloudflare.com
- deploys the SvelteKit PWA
- global CDN, automatic HTTPS, zero egress fees
- **free tier:** unlimited sites, 500 builds/month, 100k requests/day
- deployed via GitHub Actions on push to `main`
- dashboard: https://dash.cloudflare.com

### Render (backend api)

- https://render.com
- runs the Go API as a Docker container
- managed HTTPS termination - no nginx sidecar needed in production
- **free tier:** 750 hours/month, auto-deploy from GitHub
- sleeps after 15 min idle; kept warm via UptimeRobot
- service config: [`render.yaml`](../render.yaml)
- dashboard: https://dashboard.render.com

## database

### Neon PostgreSQL

- https://neon.tech
- serverless managed postgres, connects from Render via `DATABASE_URL`
- **free tier:** 0.5 GB storage, 1 project, auto-suspend after 5 min idle
- cold-start latency on first query after suspend (~1-2s), mitigated by connection pooling (pgxpool)
- automated point-in-time backups included
- dashboard: https://console.neon.tech

## object storage

### Cloudflare R2

- https://developers.cloudflare.com/r2
- stores product images via S3-compatible API
- Go `media` module uses `aws-sdk-go-v2` with custom R2 endpoint
- images served directly from R2 public bucket URL
- **free tier:** 10 GB storage, 10 million reads/month, no egress fees
- managed in Cloudflare dashboard under R2 section

## payments

### Square

- https://squareup.com/ca/en
- handles all online payments via Web Payments SDK (embedded card form) and Payments API (server-side charge)
- webhook at `POST /api/webhooks/square` for payment event processing
- HMAC signature verified on every webhook
- same provider used for in-person sales - unified reporting
- **cost:** 2.9% + 30c per online transaction
- sandbox environment available for testing
- developer dashboard: https://developer.squareup.com/apps
- sandbox test cards: https://developer.squareup.com/docs/devtools/sandbox/payments

## dns and cdn

### Cloudflare (dns)

- https://cloudflare.com
- manages DNS for the domain
- provides CDN caching, DDoS protection, and SSL certificates
- **free tier:** unlimited DNS queries, basic CDN, universal SSL

## monitoring

### UptimeRobot

- https://uptimerobot.com
- pings `GET /health` every 5 minutes
- keeps the Render free tier from sleeping during active hours
- alerts on downtime via email
- **free tier:** 50 monitors, 5-minute intervals
- dashboard: https://dashboard.uptimerobot.com

## ci/cd

### GitHub Actions

- https://github.com/features/actions
- runs on every PR and push to `main`
- **api:** lint (golangci-lint), test, build - see [ci-api.yml](../.github/workflows/ci-api.yml)
- **frontend:** format check (prettier), type check (svelte-check), build - see [ci-frontend.yml](../.github/workflows/ci-frontend.yml)
- **security:** weekly govulncheck + pnpm audit - see [security.yml](../.github/workflows/security.yml)
- **deploy:** frontend to Cloudflare Pages on push to main - see [deploy-frontend.yml](../.github/workflows/deploy-frontend.yml)
