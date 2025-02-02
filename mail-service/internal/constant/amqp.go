package constant

const (
	SendVerificationExchange = "notifications"
	SendVerificationKey      = "send.verification"
	SendVerificationQueue    = "verification"

	AccountVerifiedExchange = "notifications"
	AccountVerifiedKey      = "account.verified"
	AccountVerifiedQueue    = "verified"
)

const (
	AMQPRetryDelay = 3
	AMQPRetryLimit = 3
)
