package dto

import "github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"

type PaymentReminderEvent struct {
	Email   string `json:"email"`
	OrderID int64  `json:"order_id"`
}

func (e PaymentReminderEvent) Key() string {
	return constant.PaymentReminderKey
}
