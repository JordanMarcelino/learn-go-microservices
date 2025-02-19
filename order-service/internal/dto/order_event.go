package dto

import (
	"fmt"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/entity"
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

func ToOrderCreatedEvent(order *entity.Order) *OrderCreatedEvent {
	items := []*OrderItemEvent{}
	for _, item := range order.Items {
		items = append(items, &OrderItemEvent{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return &OrderCreatedEvent{
		Id:    order.ID,
		Items: items,
	}
}

type OrderCancelledEvent struct {
	Id    int64             `json:"id"`
	Items []*OrderItemEvent `json:"items"`
}

func (e *OrderCancelledEvent) ID() string {
	return fmt.Sprintf("%d", e.Id)
}

func ToOrderCancelledEvent(order *entity.Order) *OrderCancelledEvent {
	items := []*OrderItemEvent{}
	for _, item := range order.Items {
		items = append(items, &OrderItemEvent{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return &OrderCancelledEvent{
		Id:    order.ID,
		Items: items,
	}
}

type AutoCancelEvent struct {
	OrderID int64  `json:"order_id"`
	UserID  int64  `json:"user_id"`
	Email   string `json:"email"`
}

func (e *AutoCancelEvent) Key() string {
	return constant.AutoCancelKey
}

type CancelNotificationEvent struct {
	OrderID int64  `json:"order_id"`
	UserID  int64  `json:"user_id"`
	Email   string `json:"email"`
}

func (e *CancelNotificationEvent) Key() string {
	return constant.CancelNotificationKey
}

type OrderSuccessEvent struct {
	OrderID int64  `json:"order_id"`
	UserID  int64  `json:"user_id"`
	Email   string `json:"email"`
}

func (e *OrderSuccessEvent) Key() string {
	return constant.OrderSuccessKey
}
