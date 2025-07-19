package main

import (
	"auth/internal/config"
	"auth/internal/delivery/rest/login"
	"auth/internal/delivery/rest/registration"
	"auth/internal/delivery/rest/verification"
	"auth/internal/kafka"
	"auth/internal/logger"
	"auth/internal/repository/postgres/outbox"
	"auth/internal/repository/postgres/user"
	"auth/internal/server"
	"auth/service/auth"
	outboxworker "auth/service/outbox"
	"context"
	"github.com/R1ckNash/Bank/pkg/postgres"
	"github.com/R1ckNash/Bank/pkg/transaction_manager"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// config
	cfg := config.MustLoad()

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

	// kafka producer
	kafkaProducer, err := kafka.NewProducer([]string{"kafka:9092"}, logg)
	if err != nil {
		logg.Fatal("cannot create kafka producer", zap.Error(err))
	}
	defer kafkaProducer.Close()

	// repository
	txManager := transaction_manager.New(pool)
	userRepo := user.New(txManager)
	outboxRepo := outbox.New(txManager)

	outboxWorker := outboxworker.New(outboxRepo, kafkaProducer, logg, "auth-events", 10, 2*time.Second)
	go outboxWorker.Run(ctx)

	// services
	authService := auth.NewAuthService(auth.Deps{
		UserRepository:     userRepo,
		OutboxRepository:   outboxRepo,
		TransactionManager: txManager,
		Producer:           kafkaProducer,
		JwtSecret:          cfg.JWTSecret,
		Logger:             logg,
	})

	// delivery
	r := chi.NewRouter()
	r.Use(
		middleware.Recoverer,
		middleware.Logger,
		middleware.Heartbeat("/ping"),
		middleware.RequestID,
		middleware.URLFormat,
		middleware.Timeout(5*time.Second),
	)

	r.Route("/user", func(r chi.Router) {
		r.Post("/registration", registration.New(logg, authService))
		r.Post("/login", login.New(logg, authService))
		r.Get("/verify/{user_id}", verification.New(logg, authService))
	})

	application := server.New(logg, r, cfg.Port)

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
