package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MobinaToorani/retrosnack/internal/auth"
	"github.com/MobinaToorani/retrosnack/internal/catalog"
	"github.com/MobinaToorani/retrosnack/internal/instagram"
	"github.com/MobinaToorani/retrosnack/internal/inventory"
	"github.com/MobinaToorani/retrosnack/internal/media"
	"github.com/MobinaToorani/retrosnack/internal/orders"
	"github.com/MobinaToorani/retrosnack/internal/payments"
	"github.com/MobinaToorani/retrosnack/pkg/config"
	"github.com/MobinaToorani/retrosnack/pkg/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("database ping failed", "error", err)
		os.Exit(1)
	}
	logger.Info("database connected")

	// Wire domain modules
	authRepo := auth.NewRepository(pool)
	authSvc := auth.NewService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authSvc)

	catalogRepo := catalog.NewRepository(pool)
	catalogSvc := catalog.NewService(catalogRepo)
	catalogHandler := catalog.NewHandler(catalogSvc)

	inventoryRepo := inventory.NewRepository(pool)
	inventorySvc := inventory.NewService(inventoryRepo)
	inventoryHandler := inventory.NewHandler(inventorySvc)

	ordersRepo := orders.NewRepository(pool)
	ordersSvc := orders.NewService(ordersRepo, inventorySvc)
	ordersHandler := orders.NewHandler(ordersSvc)

	paymentsSvc := payments.NewService(ordersSvc, cfg.StripeSecretKey, cfg.StripeWebhookSecret)
	paymentsHandler := payments.NewHandler(paymentsSvc)

	instagramRepo := instagram.NewRepository(pool)
	instagramSvc := instagram.NewService(instagramRepo)
	instagramHandler := instagram.NewHandler(instagramSvc)

	mediaSvc := media.NewService(cfg)
	mediaHandler := media.NewHandler(mediaSvc)

	// Router
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.CORS(cfg.Env))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api", func(r chi.Router) {
		authHandler.Register(r)
		catalogHandler.Register(r)
		inventoryHandler.Register(r)
		ordersHandler.Register(r)
		paymentsHandler.Register(r)
		instagramHandler.Register(r)
		mediaHandler.Register(r)
	})

	addr := ":" + cfg.Port
	logger.Info("server starting", "addr", addr, "env", cfg.Env)

	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
