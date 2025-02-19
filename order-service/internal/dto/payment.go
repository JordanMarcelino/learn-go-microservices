package dto

import "github.com/jordanmarcelino/learn-go-microservices/order-service/internal/entity"

type PaymentResponse struct {
	OrderID int64  `json:"order_id"`
	Status  string `json:"status"`
}

type PayOrderRequest struct {
	CustomerID    int64
	CustomerEmail string
	RequestID     string `json:"request_id" binding:"required"`
	OrderID       int64  `json:"order_id" binding:"required"`
}

func ToPaymentResponse(order *entity.Order) *PaymentResponse {
	return &PaymentResponse{
		OrderID: order.ID,
		Status:  order.Status,
	}
}
