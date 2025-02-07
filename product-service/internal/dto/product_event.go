package dto

import (
	"fmt"

	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/entity"
	"github.com/shopspring/decimal"
)

type ProductCreatedEvent struct {
	Id          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Quantity    int             `json:"quantity"`
}

func (e ProductCreatedEvent) ID() string {
	return fmt.Sprintf("%d", e.Id)
}

func ToProductCreatedEvent(product *entity.Product) *ProductCreatedEvent {
	return &ProductCreatedEvent{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
	}
}
