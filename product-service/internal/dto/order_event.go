package dto

import (
	"fmt"
)

type OrderCreatedEvent struct {
	Id    int64             `json:"id"`
	Items []*OrderItemEvent `json:"items"`
}

type OrderItemEvent struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func (e *OrderCreatedEvent) ID() string {
	return fmt.Sprintf("%d", e.Id)
}

type OrderCancelledEvent struct {
	Id    int64             `json:"id"`
	Items []*OrderItemEvent `json:"items"`
}

func (e *OrderCancelledEvent) ID() string {
	return fmt.Sprintf("%d", e.Id)
}
