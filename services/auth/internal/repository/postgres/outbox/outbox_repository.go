package outbox

import (
	"auth/domain"
	"auth/internal/repository/postgres/user"
	"context"
	"errors"
	pkgerrors "github.com/R1ckNash/Bank/pkg/errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type OutboxRepository struct {
	driver user.QueryEngineProvider
}

func New(driver user.QueryEngineProvider) *OutboxRepository {
	return &OutboxRepository{driver: driver}
}

func (s *OutboxRepository) AddMessage(ctx context.Context, msg *domain.OutboxMessage) error {
	const api = "postgres.AddMessage"

	query := `INSERT INTO auth_outbox (aggregate_type, aggregate_id, type, payload, created_at)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := s.driver.GetQueryEngine(ctx).QueryRow(ctx, query, msg.AggregateType, msg.AggregateID, msg.Type, msg.Payload, msg.CreatedAt).Scan(&msg.ID)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			return pkgerrors.Wrap(api, domain.ErrAlreadyExists)
		}
		return pkgerrors.Wrap(api, err)
	}
	return nil
}

func (s *OutboxRepository) GetUnsentMessages(ctx context.Context, limit int) ([]*domain.OutboxMessage, error) {
	const api = "postgres.GetUnsentMessages"

	query := `SELECT id, aggregate_type, aggregate_id, type, payload, created_at, sent_at
			  FROM auth_outbox
			  WHERE sent_at IS NULL
			  ORDER BY created_at
			  LIMIT $1`
	rows, err := s.driver.GetQueryEngine(ctx).Query(ctx, query, limit)
	if err != nil {
		return nil, pkgerrors.Wrap(api, err)
	}
	defer rows.Close()

	result := make([]*domain.OutboxMessage, 0)
	for rows.Next() {
		var msg domain.OutboxMessage
		err = rows.Scan(
			&msg.ID,
			&msg.AggregateType,
			&msg.AggregateID,
			&msg.Type,
			&msg.Payload,
			&msg.CreatedAt,
			&msg.SentAt,
		)
		if err != nil {
			return nil, pkgerrors.Wrap(api, err)
		}
		result = append(result, &msg)
	}
	return result, nil
}

func (s *OutboxRepository) MarkAsSent(ctx context.Context, id int64) error {
	const api = "postgres.MarkAsSent"

	now := time.Now()
	query := `UPDATE auth_outbox SET sent_at=$1 WHERE id=$2`
	_, err := s.driver.GetQueryEngine(ctx).Exec(ctx, query, now, id)
	if err != nil {
		return pkgerrors.Wrap(api, err)
	}
	return nil
}
