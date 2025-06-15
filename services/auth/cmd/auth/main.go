package main

import (
	"auth/internal/app/config"
	"auth/internal/app/delivery/rest/login"
	"auth/internal/app/delivery/rest/registration"
	"auth/internal/app/delivery/rest/verification"
	"auth/internal/app/logger"
	userstorage "auth/internal/app/repository/postgres"
	"auth/internal/app/server"
	"auth/service"
	"context"
	"github.com/R1ckNash/Bank/pkg/postgres"
	"github.com/R1ckNash/Bank/pkg/transaction_manager"
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

	// pg conn
	pool, err := postgres.NewConnectionPool(ctx, cfg.DBUrl,
		postgres.WithMaxConnIdleTime(5*time.Minute),
		postgres.WithMaxConnLifeTime(time.Hour),
		postgres.WithMaxConnectionsCount(10),
		postgres.WithMinConnectionsCount(5),
	)
	if err != nil {
		logg.Fatal("db connection error", zap.Error(err))
	}

	// repository
	txManager := transaction_manager.New(pool)
	userRepo := userstorage.New(txManager)

	// services
	authService := service.NewAuthService(service.Deps{
		UserRepository:     userRepo,
		TransactionManager: txManager,
		JwtSecret:          cfg.JWTSecret,
		Logger:             logg,
	})

	// delivery
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.URLFormat,
		middleware.Recoverer,
		middleware.Timeout(5*time.Second),
	)

	r.Post("/registration", registration.New(logg, authService))
	r.Post("/login", login.New(logg, authService))
	r.Get("/user/verify/{user_id}", verification.New(logg, authService))

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
