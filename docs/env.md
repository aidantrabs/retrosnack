# environment variables

## backend (services/api)

see `services/api/.env.example` for a copyable template.

### required

| variable | description |
|----------|-------------|
| `DATABASE_URL` | postgres connection string |
| `JWT_SECRET` | secret for signing JWT tokens |

### optional

| variable | default | description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP port for the api |
| `ENV` | `development` | `development` or `production` - controls CORS, logging |
| `SQUARE_ACCESS_TOKEN` | - | Square API access token |
| `SQUARE_APPLICATION_ID` | - | Square application ID (used by frontend payment config) |
| `SQUARE_LOCATION_ID` | - | Square location ID |
| `SQUARE_ENVIRONMENT` | `sandbox` | `sandbox` or `production` |
| `SQUARE_WEBHOOK_SIG_KEY` | - | HMAC key for verifying Square webhook signatures |
| `SQUARE_WEBHOOK_NOTIF_URL` | - | full URL Square sends webhooks to |
| `R2_ACCOUNT_ID` | - | Cloudflare account ID |
| `R2_ACCESS_KEY_ID` | - | R2 S3-compatible access key |
| `R2_SECRET_ACCESS_KEY` | - | R2 S3-compatible secret key |
| `R2_BUCKET_NAME` | - | R2 bucket name |
| `R2_PUBLIC_URL` | - | public URL for serving images |

## frontend (apps/frontend)

see `apps/frontend/.env.example`.

| variable | description |
|----------|-------------|
| `PUBLIC_API_URL` | backend API URL (e.g. `http://localhost:8080` for dev) |

## github actions secrets

set these in the repo's Settings > Secrets and variables > Actions.

| secret | used by | description |
|--------|---------|-------------|
| `CLOUDFLARE_API_TOKEN` | deploy-frontend | Cloudflare API token with Pages edit permission |
| `CLOUDFLARE_ACCOUNT_ID` | deploy-frontend | Cloudflare account ID |
| `PUBLIC_API_URL` | deploy-frontend | production API URL for the frontend build |
