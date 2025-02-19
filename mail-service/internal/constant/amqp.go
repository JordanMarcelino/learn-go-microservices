package constant

const (
	SendVerificationExchange = "notifications"
	SendVerificationKey      = "send.verification"
	SendVerificationQueue    = "verification"

	AccountVerifiedExchange = "notifications"
	AccountVerifiedKey      = "account.verified"
	AccountVerifiedQueue    = "verified"

	PaymentReminderExchange = "notifications"
	PaymentReminderKey      = "send.payment.reminder"
	PaymentReminderQueue    = "reminder"

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
