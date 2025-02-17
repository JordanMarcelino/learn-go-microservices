package constant

const (
	PaymentReminderExchange = "notifications"
	PaymentReminderKey      = "send.payment.reminder"
	PaymentReminderQueue    = "reminder"

	AutoCancelExchange = "orders"
	AutoCancelKey      = "order.auto.cancel"
	AutoCancelQueue    = "order-auto-cancel"
)

const (
	AMQPRetryDelay = 3
	AMQPRetryLimit = 3
)
