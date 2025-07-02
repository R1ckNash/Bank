package main

import (
	"account/internal/config"
	"account/internal/http-server/handlers/get_handler"
	"account/internal/http-server/handlers/post_handler"
	http_server "account/internal/http-server/server"
	"account/internal/kafka"
	httpdelivery "account/internal/middleware"
	"account/internal/repository/account_storage"
	"account/internal/services/account"
	"context"
	"fmt"
	"github.com/R1ckNash/Bank/pkg/middleware/auth"
	"github.com/R1ckNash/Bank/pkg/postgres"
	"github.com/R1ckNash/Bank/pkg/transaction_manager"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {

	parent := context.Background()
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting account-service with", slog.String("env", cfg.Env))

	// repository
	pool, err := postgres.NewConnectionPool(parent, cfg.DBUrl,
		postgres.WithMaxConnIdleTime(5*time.Minute),
		postgres.WithMaxConnLifeTime(time.Hour),
		postgres.WithMaxConnectionsCount(10),
		postgres.WithMinConnectionsCount(5),
	)
	if err != nil {
		panic("db connection error")
	}

	// kafka
	producer, err := kafka.NewProducer([]string{cfg.Kafka.Host}, log)
	if err != nil {
		panic("could not create kafka producer")
	}
	defer producer.Close()

	txManager := transaction_manager.New(pool)
	storage := account_storage.New(txManager)

	accountService := account.NewAccountService(account.Deps{
		AccountStorage:     storage,
		TransactionManager: txManager,
		Logger:             log,
		EventProducer:      producer,
	})

	router := chi.NewRouter()

	router.Use(
		middleware.Heartbeat("/ping"),
		httpdelivery.PrometheusMiddleware,
		middleware.RequestID,
		middleware.Recoverer,
		middleware.URLFormat,
		auth.AuthMiddleware(cfg.JWTSecret),
	)

	router.Route("/account", func(r chi.Router) {
		r.Get("/{accountId}", get_handler.New(log, accountService))
		r.Post("/create", post_handler.New(log, accountService, cfg))
	})

	router.Handle("/metrics", promhttp.Handler())

	server := http_server.New(log, router, cfg.Port)

	ctx, stop := signal.NotifyContext(parent,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	run(ctx, server, log)
}

func run(ctx context.Context, server *http_server.Server, log *slog.Logger) error {
	// Start HTTP server in a goroutine
	go server.MustRun()

	// Wait until we receive a shutdown signal
	<-ctx.Done()

	log.Info("server: shutting down server gracefully")

	// Create a context with a 20-second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Attempt a graceful shutdown
	if err := server.Stop(shutdownCtx); err != nil {
		return fmt.Errorf("server: shutdown: %w", err)
	}

	log.Info("server: shutdown")
	return nil
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
