package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
)

type Aggregate interface {
	Apply(event Event) error
}

type Product struct {
	ID    string
	Name  string
	Price float64
}

func (p *Product) Apply(event Event) error {
	switch event.Type {
	case ProductCreatedType:
		payload := event.Payload.(map[string]interface{})
		p.ID = payload["id"].(string)
		p.Name = payload["name"].(string)
		p.Price = payload["price"].(float64)
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
	return nil
}

type Inventory struct {
	ProductID string
	Stock     int
}

func (i *Inventory) Apply(event Event) error {
	switch event.Type {
	case ProductCreatedType:
		payload := event.Payload.(map[string]interface{})
		i.ProductID = payload["id"].(string)
		i.Stock = 0
	case StockIncreasedType:
		payload := event.Payload.(map[string]interface{})
		i.Stock += int(payload["quantity"].(float64))
	case StockDecreasedType:
		payload := event.Payload.(map[string]interface{})
		quantity := int(payload["quantity"].(float64))
		if i.Stock < quantity {
			return errors.New("insufficient stock")
		}
		i.Stock -= quantity
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
	return nil
}

type Order struct {
	ID        string
	ProductID string
	Quantity  int
	Status    string
}

func (o *Order) Apply(event Event) error {
	switch event.Type {
	case OrderPlacedType:
		payload := event.Payload.(map[string]interface{})
		o.ID = payload["order_id"].(string)
		o.ProductID = payload["product_id"].(string)
		o.Quantity = int(payload["quantity"].(float64))
		o.Status = "Placed"
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
	return nil
}
