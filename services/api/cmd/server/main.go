package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aidantrabs/kenko"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/retrosnack-clothing/retrosnack/internal/auth"
	"github.com/retrosnack-clothing/retrosnack/internal/catalog"
	"github.com/retrosnack-clothing/retrosnack/internal/instagram"
	"github.com/retrosnack-clothing/retrosnack/internal/inventory"
	"github.com/retrosnack-clothing/retrosnack/internal/media"
	"github.com/retrosnack-clothing/retrosnack/internal/orders"
	"github.com/retrosnack-clothing/retrosnack/internal/payments"
	"github.com/retrosnack-clothing/retrosnack/pkg/config"
	"github.com/retrosnack-clothing/retrosnack/pkg/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("parse database url: %w", err)
	}
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute
	poolCfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("database ping: %w", err)
	}
	logger.Info("database connected")

	// wire domain modules
	authRepo := auth.NewRepository(pool)
	authSvc := auth.NewService(authRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authSvc)

	catalogRepo := catalog.NewRepository(pool)
	catalogSvc := catalog.NewService(catalogRepo)
	catalogHandler := catalog.NewHandler(catalogSvc, cfg.JWTSecret)

	inventoryRepo := inventory.NewRepository(pool)
	inventorySvc := inventory.NewService(inventoryRepo)
	inventoryHandler := inventory.NewHandler(inventorySvc)

	ordersRepo := orders.NewRepository(pool)
	ordersSvc := orders.NewService(ordersRepo, inventorySvc)
	ordersHandler := orders.NewHandler(ordersSvc, cfg.JWTSecret)

	paymentsSvc := payments.NewService(ordersSvc, cfg.SquareAccessToken, cfg.SquareLocationID, cfg.SquareWebhookSigKey, cfg.SquareWebhookNotifURL)
	paymentsHandler := payments.NewHandler(paymentsSvc, cfg.SquareApplicationID, cfg.SquareLocationID)

	instagramRepo := instagram.NewRepository(pool)
	instagramSvc := instagram.NewService(instagramRepo)
	instagramHandler := instagram.NewHandler(instagramSvc, cfg.JWTSecret)

	mediaRepo := media.NewRepository(pool)
	mediaSvc := media.NewService(cfg, mediaRepo)
	mediaHandler := media.NewHandler(mediaSvc, cfg.JWTSecret)

	// health monitoring
	k, err := kenko.New(
		kenko.WithTarget("square", "https://connect.squareup.com/v2/locations"),
		kenko.WithInterval(30*time.Second),
		kenko.WithLogger(logger),
	)
	if err != nil {
		return fmt.Errorf("init health monitor: %w", err)
	}
	go k.Run(context.Background())

	// router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.CORS(cfg.Env))
	r.Use(middleware.SecureHeaders)
	r.Use(middleware.MaxBodySize(10 << 20)) // 10 MB

	checker := k.Checker()
	healthHandler := kenko.HandleHealth(checker)
	r.Get("/health", healthHandler)
	r.Head("/health", healthHandler)
	r.Get("/ready", kenko.HandleReady(checker))
	r.Get("/status", kenko.HandleStatus(checker))

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
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("server starting", "addr", addr, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server stopped", "error", err)
		}
	}()

	<-done
	logger.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("forced shutdown: %w", err)
	}

	logger.Info("server stopped gracefully")
	return nil
}
