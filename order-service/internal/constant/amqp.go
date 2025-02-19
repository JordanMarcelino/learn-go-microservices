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

	OrderSuccessExchange = "notifications"
	OrderSuccessKey      = "send.order.success"
	OrderSuccessQueue    = "order-success"
)

const (
	AMQPRetryDelay = 3
	AMQPRetryLimit = 3
)
