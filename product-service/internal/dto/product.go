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

type UpdateProductRequest struct {
	ID          int64
	Name        string          `json:"name" binding:"required,max=255"`
	Description string          `json:"description" binding:"required"`
	Price       decimal.Decimal `json:"price" binding:"required,dgt=0"`
	Quantity    int             `json:"quantity" binding:"required,min=0"`
}

type GetProductRequest struct {
	ID int64
}

type DeleteProductRequest struct {
	ID int64
}

type SearchProductRequest struct {
	Name        string `form:"name"`
	Description string `form:"description"`
	Limit       int64  `form:"limit" binding:"gte=1,lte=25"`
	Page        int64  `form:"page" binding:"gte=1"`
}

func ToProductResponses(products []*entity.Product) []*ProductResponse {
	res := []*ProductResponse{}
	for _, product := range products {
		res = append(res, ToProductResponse(product))
	}
	return res
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
