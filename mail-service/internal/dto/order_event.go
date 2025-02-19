package dto

import "github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/constant"

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
