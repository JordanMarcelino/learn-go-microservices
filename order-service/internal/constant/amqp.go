package constant

const (
	PaymentReminderExchange = "notifications"
	PaymentReminderKey      = "send.payment.reminder"
	PaymentReminderQueue    = "reminder"

	AutoCancelExchange = "orders"
	AutoCancelKey      = "order.auto.cancel"
	AutoCancelQueue    = "order-auto-cancel"

	CancelNotificationExchange = "notifications"
	CancelNotificationKey      = "send.cancel.notification"
	CancelNotificationQueue    = "cancel-notification"
)

const (
	AMQPRetryDelay = 3
	AMQPRetryLimit = 3
)
