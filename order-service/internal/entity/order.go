package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID          int64
	RequestID   string
	CustomerID  int64
	TotalAmount decimal.Decimal
	Description string
	Status      string
	Items       []*OrderItem
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderItem struct {
	ID        int64
	OrderID   int64
	ProductID int64
	Quantity  int
	Price     decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
}
