package dto

import "github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"

type PaymentReminderEvent struct {
	OrderID int64  `json:"order_id"`
	UserID  int64  `json:"user_id"`
	Email   string `json:"email"`
	DueDate string `json:"due_date"`
}

func (e PaymentReminderEvent) Key() string {
	return constant.PaymentReminderKey
}
