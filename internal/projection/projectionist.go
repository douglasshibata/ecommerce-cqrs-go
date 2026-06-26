package projection

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jules/ecommerce-cqrs-go/internal/domain"
)

type Projectionist struct {
	db *sql.DB
}

func NewProjectionist(db *sql.DB) *Projectionist {
	return &Projectionist{db: db}
}

func (p *Projectionist) Project(ctx context.Context, event domain.Event) error {
	switch event.Type {
	case domain.ProductCreatedType:
		var payload domain.ProductCreatedPayload
		data, _ := json.Marshal(event.Payload)
		json.Unmarshal(data, &payload)

		tx, err := p.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		_, err = tx.ExecContext(ctx, "INSERT INTO products (id, name, price) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING", payload.ID, payload.Name, payload.Price)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO inventory (product_id, stock) VALUES ($1, 0) ON CONFLICT (product_id) DO NOTHING", payload.ID)
		if err != nil {
			return err
		}

		return tx.Commit()

	case domain.StockIncreasedType:
		var payload domain.StockIncreasedPayload
		data, _ := json.Marshal(event.Payload)
		json.Unmarshal(data, &payload)

		tx, err := p.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		_, err = tx.ExecContext(ctx, "UPDATE inventory SET stock = stock + $1, updated_at = NOW() WHERE product_id = $2", payload.Quantity, payload.ProductID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO inventory_history (product_id, change, reason) VALUES ($1, $2, $3)", payload.ProductID, payload.Quantity, "Stock Increase")
		if err != nil {
			return err
		}

		return tx.Commit()

	case domain.StockDecreasedType:
		var payload domain.StockDecreasedPayload
		data, _ := json.Marshal(event.Payload)
		json.Unmarshal(data, &payload)

		tx, err := p.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		_, err = tx.ExecContext(ctx, "UPDATE inventory SET stock = stock - $1, updated_at = NOW() WHERE product_id = $2", payload.Quantity, payload.ProductID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO inventory_history (product_id, change, reason) VALUES ($1, $2, $3)", payload.ProductID, -payload.Quantity, "Stock Decrease")
		if err != nil {
			return err
		}

		return tx.Commit()

	case domain.OrderPlacedType:
		var payload domain.OrderPlacedPayload
		data, _ := json.Marshal(event.Payload)
		json.Unmarshal(data, &payload)

		_, err := p.db.ExecContext(ctx, "INSERT INTO orders (id, product_id, quantity, total_price, status) VALUES ($1, $2, $3, $4, $5)",
			payload.OrderID, payload.ProductID, payload.Quantity, payload.TotalPrice, "Placed")
		return err

	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}
