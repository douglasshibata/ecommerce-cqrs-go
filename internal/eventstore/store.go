package eventstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/jules/ecommerce-cqrs-go/internal/domain"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Save(ctx context.Context, events ...domain.Event) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, event := range events {
		payload, err := json.Marshal(event.Payload)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx,
			"INSERT INTO events (id, aggregate_id, aggregate_type, event_type, payload, version, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			event.ID, event.AggregateID, event.AggregateType, event.Type, payload, event.Version, event.CreatedAt,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) GetEvents(ctx context.Context, aggregateID string) ([]domain.Event, error) {
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, aggregate_id, aggregate_type, event_type, payload, version, created_at FROM events WHERE aggregate_id = $1 ORDER BY version ASC",
		aggregateID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var e domain.Event
		var payload []byte
		if err := rows.Scan(&e.ID, &e.AggregateID, &e.AggregateType, &e.Type, &payload, &e.Version, &e.CreatedAt); err != nil {
			return nil, err
		}

		var p interface{}
		if err := json.Unmarshal(payload, &p); err != nil {
			return nil, err
		}
		e.Payload = p
		events = append(events, e)
	}
	return events, nil
}
