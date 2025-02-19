package dto

import "github.com/shopspring/decimal"

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

type GetOrderRequest struct {
	OrderID int64
	UserID  int64
	Email   string
}
