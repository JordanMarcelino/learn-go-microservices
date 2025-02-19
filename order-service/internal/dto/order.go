package dto

import (
	"time"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/entity"
	"github.com/shopspring/decimal"
)

type OrderResponse struct {
	ID          int64                `json:"id"`
	TotalAmount decimal.Decimal      `json:"total_amount"`
	Description string               `json:"description"`
	Status      string               `json:"status"`
	Items       []*OrderItemResponse `json:"items"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}

type OrderItemResponse struct {
	ID        int64           `json:"id"`
	ProductID int64           `json:"product_id"`
	Quantity  int             `json:"quantity"`
	Price     decimal.Decimal `json:"price"`
}

type CreateOrderRequest struct {
	CustomerID    int64
	CustomerEmail string
	RequestID     string              `json:"request_id" binding:"required"`
	Description   string              `json:"description"`
	Items         []*OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type OrderItemRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

type GetOrderRequest struct {
	OrderID int64
}

type CancelOrderRequest struct {
	CustomerID    int64
	CustomerEmail string
	RequestID     string `json:"request_id" binding:"required"`
	OrderID       int64  `json:"order_id" binding:"required"`
}

type SearchOrderRequest struct {
	StartDate time.Time `form:"start-date" time_format:"02-01-2006" binding:"required"`
	EndDate   time.Time `form:"end-date" time_format:"02-01-2006" binding:"required,gtefield=StartDate"`
	Status    string    `form:"status" binding:"omitempty,oneof=PENDING SUCCESS CANCELLED"`
	Limit     int64     `form:"limit" binding:"gte=1,lte=25"`
	Page      int64     `form:"page" binding:"gte=1"`
}

func ToOrderResponses(orders []*entity.Order) []*OrderResponse {
	res := []*OrderResponse{}
	for _, order := range orders {
		res = append(res, ToOrderResponse(order))
	}
	return res

}

func ToOrderResponse(order *entity.Order) *OrderResponse {
	items := []*OrderItemResponse{}
	for _, item := range order.Items {
		items = append(items, &OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	return &OrderResponse{
		ID:          order.ID,
		TotalAmount: order.TotalAmount,
		Description: order.Description,
		Status:      order.Status,
		Items:       items,
		CreatedAt:   order.CreatedAt.String(),
		UpdatedAt:   order.UpdatedAt.String(),
	}
}
