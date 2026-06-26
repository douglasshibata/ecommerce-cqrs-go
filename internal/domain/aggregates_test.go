package domain

import (
	"testing"
)

func TestInventory_Apply(t *testing.T) {
	inventory := &Inventory{}

	// Test ProductCreated
	createdEvent := Event{
		Type: ProductCreatedType,
		Payload: map[string]interface{}{
			"id":    "p1",
			"name":  "Laptop",
			"price": 1200.0,
		},
	}
	if err := inventory.Apply(createdEvent); err != nil {
		t.Fatalf("Failed to apply ProductCreated: %v", err)
	}
	if inventory.ProductID != "p1" || inventory.Stock != 0 {
		t.Errorf("Expected product p1 with 0 stock, got %s and %d", inventory.ProductID, inventory.Stock)
	}

	// Test StockIncreased
	increasedEvent := Event{
		Type: StockIncreasedType,
		Payload: map[string]interface{}{
			"product_id": "p1",
			"quantity":   10.0, // JSON unmarshal uses float64
		},
	}
	if err := inventory.Apply(increasedEvent); err != nil {
		t.Fatalf("Failed to apply StockIncreased: %v", err)
	}
	if inventory.Stock != 10 {
		t.Errorf("Expected stock 10, got %d", inventory.Stock)
	}

	// Test StockDecreased
	decreasedEvent := Event{
		Type: StockDecreasedType,
		Payload: map[string]interface{}{
			"product_id": "p1",
			"quantity":   3.0,
		},
	}
	if err := inventory.Apply(decreasedEvent); err != nil {
		t.Fatalf("Failed to apply StockDecreased: %v", err)
	}
	if inventory.Stock != 7 {
		t.Errorf("Expected stock 7, got %d", inventory.Stock)
	}

	// Test Insufficient Stock
	insufficientEvent := Event{
		Type: StockDecreasedType,
		Payload: map[string]interface{}{
			"product_id": "p1",
			"quantity":   10.0,
		},
	}
	if err := inventory.Apply(insufficientEvent); err == nil {
		t.Error("Expected error for insufficient stock, got nil")
	}
}
