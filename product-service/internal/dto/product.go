package dto

import (
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/entity"
	"github.com/shopspring/decimal"
)

type ProductResponse struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Quantity    int             `json:"quantity"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string          `json:"name" binding:"required,max=255"`
	Description string          `json:"description" binding:"required"`
	Price       decimal.Decimal `json:"price" binding:"required,dgt=0"`
	Quantity    int             `json:"quantity" binding:"required,min=0"`
}

func ToProductResponse(product *entity.Product) *ProductResponse {
	return &ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CreatedAt:   product.CreatedAt.String(),
		UpdatedAt:   product.UpdatedAt.String(),
	}
}
