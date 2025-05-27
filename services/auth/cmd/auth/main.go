package main

import (
	"Bank/pkg/postgres"
	"Bank/pkg/transaction_manager"
	"auth/internal/app/config"
	"auth/internal/app/delivery/http/registration"
	"auth/internal/app/logger"
	"auth/internal/app/repository/user_storage"
	"auth/internal/app/server"
	"auth/internal/app/services/auth"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// config
	cfg := config.LoadConfig()

	// logger init
	logg, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("error creating logger: %v", err)
	}
	defer logg.Sync()

	// repository
	pool, err := postgres.NewConnectionPool(ctx, cfg.DBUrl,
		postgres.WithMaxConnIdleTime(5*time.Minute),
		postgres.WithMaxConnLifeTime(time.Hour),
		postgres.WithMaxConnectionsCount(10),
		postgres.WithMinConnectionsCount(5),
	)
	if err != nil {
		logg.Fatal("db connection error", zap.Error(err))
	}

	txManager := transaction_manager.New(pool)
	storage := user_storage.New(txManager)

	// services
	authService := auth.NewAuthService(auth.Deps{ // Dependency injection
		UserStorage:        storage,
		TransactionManager: txManager,
	})

	// delivery
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.URLFormat,
		middleware.Recoverer,
	)

	r.Post("/url", registration.New(logg, authService))

	// app init
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		logg.Fatal("error parsing port", zap.Error(err))
	}

	application := server.New(logg, r, port)

	// graceful shutdown
	go func() {
		if err := application.Run(); err != nil {
			logg.Fatal("server error", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := application.Stop(ctxWithTimeout); err != nil {
		logg.Fatal("shutdown error", zap.Error(err))
	}

	logg.Info("server stopped")

}
