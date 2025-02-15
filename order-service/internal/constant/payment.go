package constant

import "time"

var (
	PaymentNotificationDelays = []time.Duration{4 * time.Hour, 8 * time.Hour, 16 * time.Hour}
)
