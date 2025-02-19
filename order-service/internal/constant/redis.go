package constant

import "time"

const (
	CreateOrderTTL           = 7 * time.Second
	CreateOrderRetryInterval = 1 * time.Second
	CreateOrderRetryLimit    = 3

	PayOrderTTL           = 10 * time.Second
	PayOrderRetryInterval = 2 * time.Second
	PayOrderRetryLimit    = 3
)
