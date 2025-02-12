package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       decimal.Decimal
	Quantity    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
