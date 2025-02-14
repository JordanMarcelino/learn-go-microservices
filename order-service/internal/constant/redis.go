package constant

import "time"

const (
	CreateOrderTTL           = 7 * time.Second
	CreateOrderRetryInterval = 1 * time.Second
	CreateOrderRetryLimit    = 3
)
