package repository

import (
	"context"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/entity"
)

type ProductRepository interface {
	Save(ctx context.Context, product *entity.Product) error
}

type productRepositoryImpl struct {
	DB DBTX
}

func NewProductRepository(db DBTX) ProductRepository {
	return &productRepositoryImpl{
		DB: db,
	}
}

func (r *productRepositoryImpl) Save(ctx context.Context, product *entity.Product) error {
	query := `
		INSERT INTO
			products (name, description, price, quantity)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			id, created_at, updated_at
	`

	return r.DB.QueryRowContext(ctx, query, product.Name, product.Description, product.Price, product.Quantity).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}
