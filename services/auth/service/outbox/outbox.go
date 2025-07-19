package outbox

import (
	"auth/domain"
	"auth/internal/kafka"
	"auth/service/auth"
	"context"
	"go.uber.org/zap"
	"time"
)

type Worker struct {
	Repo      auth.OutboxRepository
	Producer  *kafka.Producer
	Logger    *zap.Logger
	Topic     string
	BatchSize int
	Interval  time.Duration
}

func New(repo auth.OutboxRepository, producer *kafka.Producer, logger *zap.Logger, topic string, batchSize int, interval time.Duration) *Worker {
	return &Worker{
		Repo:      repo,
		Producer:  producer,
		Logger:    logger,
		Topic:     topic,
		BatchSize: batchSize,
		Interval:  interval,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	w.Logger.Info("Outbox worker started", zap.Duration("interval", w.Interval))

	for {
		select {
		case <-ticker.C:
			w.processBatch(ctx)
		case <-ctx.Done():
			w.Logger.Info("Outbox worker stopped")
			return
		}
	}
}

func (w *Worker) processBatch(ctx context.Context) {
	messages, err := w.Repo.GetUnsentMessages(ctx, w.BatchSize)
	if err != nil {
		w.Logger.Error("failed to get unsent messages from outbox", zap.Error(err))
		return
	}
	if len(messages) == 0 {
		return
	}

	for _, msg := range messages {
		err = w.sendToKafka(msg)
		if err != nil {
			w.Logger.Error("failed to send outbox event to kafka", zap.Error(err), zap.Int64("outbox_id", msg.ID))
			continue
		}
		if err := w.Repo.MarkAsSent(ctx, msg.ID); err != nil {
			w.Logger.Error("failed to mark outbox event as sent", zap.Error(err), zap.Int64("outbox_id", msg.ID))
		}
	}
}

func (w *Worker) sendToKafka(msg *domain.OutboxMessage) error {
	return w.Producer.SendMessage(w.Topic, msg.AggregateID.String(), []byte(msg.Payload))
}
