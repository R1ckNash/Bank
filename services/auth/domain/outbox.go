package domain

import (
	"github.com/google/uuid"
	"time"
)

type OutboxMessage struct {
	ID            int64      `db:"id"`
	AggregateType string     `db:"aggregate_type"`
	AggregateID   uuid.UUID  `db:"aggregate_id"`
	Type          string     `db:"type"`
	Payload       string     `db:"payload"`
	CreatedAt     time.Time  `db:"created_at"`
	SentAt        *time.Time `db:"sent_at"`
}
