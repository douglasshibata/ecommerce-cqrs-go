package api

import (
	"context"
	"database/sql"
	"time"
)

type QueryHandler struct {
	db *sql.DB
}

func NewQueryHandler(db *sql.DB) *QueryHandler {
	return &QueryHandler{db: db}
}

type ProductRead struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (h *QueryHandler) ListProducts(ctx context.Context) ([]ProductRead, error) {
	rows, err := h.db.QueryContext(ctx, "SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []ProductRead
	for rows.Next() {
		var p ProductRead
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

type InventoryHistoryRead struct {
	ProductID string    `json:"product_id"`
	Change    int       `json:"change"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *QueryHandler) GetInventoryHistory(ctx context.Context, productID string) ([]InventoryHistoryRead, error) {
	rows, err := h.db.QueryContext(ctx, "SELECT product_id, change, reason, created_at FROM inventory_history WHERE product_id = $1 ORDER BY created_at DESC", productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []InventoryHistoryRead
	for rows.Next() {
		var i InventoryHistoryRead
		if err := rows.Scan(&i.ProductID, &i.Change, &i.Reason, &i.CreatedAt); err != nil {
			return nil, err
		}
		history = append(history, i)
	}
	return history, nil
}

type OrderRead struct {
	ID         string    `json:"id"`
	ProductID  string    `json:"product_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (h *QueryHandler) ListOrders(ctx context.Context) ([]OrderRead, error) {
	rows, err := h.db.QueryContext(ctx, "SELECT id, product_id, quantity, total_price, status, created_at FROM orders ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []OrderRead
	for rows.Next() {
		var o OrderRead
		if err := rows.Scan(&o.ID, &o.ProductID, &o.Quantity, &o.TotalPrice, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
