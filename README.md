# retrosnack clothing

online storefront for [@retrosnack.shop](https://instagram.com/retrosnack.shop) - curated secondhand women's clothing, accessories, and shoes.

## stack

| layer | tech | hosting |
|-------|------|---------|
| frontend | SvelteKit 5, Tailwind CSS v4, PWA | Cloudflare Pages |
| backend | Go 1.26, Chi, pgx/v5, sqlc | Render |
| database | PostgreSQL | Neon |
| payments | Square Web Payments SDK | - |
| media | S3-compatible object storage | Cloudflare R2 |
| ci/cd | GitHub Actions | - |

## quickstart

```sh
cp services/api/.env.example services/api/.env
# fill in values (DATABASE_URL and JWT_SECRET are required, rest is optional)

make install    # go + node deps
make db         # start local postgres
make migrate    # run migrations
make api        # api on :8080
make frontend   # frontend on :5173
```

or all at once: `make dev`

## project layout

```
retrosnack/
├── apps/
│   └── frontend/                SvelteKit PWA
│       ├── src/
│       │   ├── lib/             shared utilities, api client, stores, components
│       │   ├── routes/          file-based routing
│       │   └── app.css          global styles (Tailwind v4)
│       ├── static/              PWA manifest, icons
│       ├── svelte.config.js
│       └── vite.config.js
│
├── services/
│   └── api/                     Go API server
│       ├── cmd/server/          entrypoint
│       ├── internal/            domain modules
│       │   ├── auth/            JWT authentication
│       │   ├── catalog/         products, categories, variants
│       │   ├── inventory/       stock tracking
│       │   ├── orders/          order lifecycle
│       │   ├── payments/        Square integration
│       │   ├── instagram/       oEmbed links
│       │   └── media/           image upload via R2
│       ├── pkg/                 shared packages (config, middleware, httputil)
│       └── db/                  migrations + sql queries
│
├── infrastructure/
│   ├── docker/Dockerfile        multi-stage Go build
│   └── nginx/nginx.conf         reverse proxy config
│
├── docs/                        project documentation
├── docker-compose.yml           local dev stack
├── render.yaml                  Render service blueprint
├── sqlc.yaml                    sqlc config
└── Makefile                     dev commands
```

## make targets

run `make help` or see the [Makefile](Makefile).

| command | what it does |
|---------|--------------|
| `make dev` | full stack via docker compose |
| `make api` | run api locally |
| `make frontend` | run frontend locally |
| `make db` | start local postgres |
| `make migrate` | apply migrations |
| `make test` | run go tests |
| `make typecheck` | run svelte-check |

## docs

- [services and infrastructure](docs/services.md) - every external service, free tier limits, and links
- [architecture](docs/architecture.md) - system design, data flow, and key decisions
- [environment variables](docs/env.md) - all config vars for backend, frontend, and ci/cd
