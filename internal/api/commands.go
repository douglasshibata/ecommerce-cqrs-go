package api

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jules/ecommerce-cqrs-go/internal/domain"
	"github.com/jules/ecommerce-cqrs-go/internal/eventstore"
)

type Projector interface {
	Project(ctx context.Context, event domain.Event) error
}

type CommandHandler struct {
	store     *eventstore.Store
	projector Projector
}

func NewCommandHandler(store *eventstore.Store, projector Projector) *CommandHandler {
	return &CommandHandler{store: store, projector: projector}
}

func (h *CommandHandler) CreateProduct(ctx context.Context, id, name string, price float64) error {
	event := domain.Event{
		ID:            uuid.New().String(),
		AggregateID:   id,
		AggregateType: "Product",
		Type:          domain.ProductCreatedType,
		Payload: domain.ProductCreatedPayload{
			ID:    id,
			Name:  name,
			Price: price,
		},
		Version:   1,
		CreatedAt: time.Now(),
	}
	if err := h.store.Save(ctx, event); err != nil {
		return err
	}
	return h.projector.Project(ctx, event)
}

func (h *CommandHandler) AdjustStock(ctx context.Context, productID string, quantity int) error {
	events, err := h.store.GetEvents(ctx, productID)
	if err != nil {
		return err
	}

	inventory := &domain.Inventory{}
	for _, e := range events {
		if err := inventory.Apply(e); err != nil {
			return err
		}
	}

	eventType := domain.StockIncreasedType
	if quantity < 0 {
		eventType = domain.StockDecreasedType
		if inventory.Stock+quantity < 0 {
			return domain.ErrInsufficientStock
		}
	}

	event := domain.Event{
		ID:            uuid.New().String(),
		AggregateID:   productID,
		AggregateType: "Inventory",
		Type:          eventType,
		Payload: domain.StockIncreasedPayload{
			ProductID: productID,
			Quantity:  abs(quantity),
		},
		Version:   len(events) + 1,
		CreatedAt: time.Now(),
	}
	if err := h.store.Save(ctx, event); err != nil {
		return err
	}
	return h.projector.Project(ctx, event)
}

func (h *CommandHandler) PlaceOrder(ctx context.Context, orderID, productID string, quantity int) error {
	inventoryEvents, err := h.store.GetEvents(ctx, productID)
	if err != nil {
		return err
	}

	inventory := &domain.Inventory{}
	for _, e := range inventoryEvents {
		if err := inventory.Apply(e); err != nil {
			return err
		}
	}

	if inventory.Stock < quantity {
		return domain.ErrInsufficientStock
	}

	product := &domain.Product{}
	for _, e := range inventoryEvents {
		if e.Type == domain.ProductCreatedType {
			if err := product.Apply(e); err != nil {
				return err
			}
		}
	}

	orderEvent := domain.Event{
		ID:            uuid.New().String(),
		AggregateID:   orderID,
		AggregateType: "Order",
		Type:          domain.OrderPlacedType,
		Payload: domain.OrderPlacedPayload{
			OrderID:    orderID,
			ProductID:  productID,
			Quantity:   quantity,
			TotalPrice: product.Price * float64(quantity),
		},
		Version:   1,
		CreatedAt: time.Now(),
	}

	stockEvent := domain.Event{
		ID:            uuid.New().String(),
		AggregateID:   productID,
		AggregateType: "Inventory",
		Type:          domain.StockDecreasedType,
		Payload: domain.StockDecreasedPayload{
			ProductID: productID,
			Quantity:  quantity,
		},
		Version:   len(inventoryEvents) + 1,
		CreatedAt: time.Now(),
	}

	if err := h.store.Save(ctx, orderEvent, stockEvent); err != nil {
		return err
	}

	// For simplicity, we project individually, but in a real app,
	// this would be asynchronous and robust to failures.
	if err := h.projector.Project(ctx, orderEvent); err != nil {
		return err
	}
	return h.projector.Project(ctx, stockEvent)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
