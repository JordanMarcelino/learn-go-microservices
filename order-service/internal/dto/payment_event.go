package dto

import "github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"

type PaymentReminderEvent struct {
	OrderID int64  `json:"order_id"`
	Email   string `json:"email"`
}

func (e PaymentReminderEvent) Key() string {
	return constant.PaymentReminderKey
}
