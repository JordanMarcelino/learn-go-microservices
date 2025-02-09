package repository

import (
	"context"
	"fmt"

	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/pageutils"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/entity"
)

type ProductRepository interface {
	Search(ctx context.Context, req *dto.SearchProductRequest) ([]*entity.Product, int64, error)
	Save(ctx context.Context, product *entity.Product) error
}

type productRepository struct {
	DB DBTX
}

func NewProductRepository(db DBTX) ProductRepository {
	return &productRepository{
		DB: db,
	}
}

func (r *productRepository) Search(ctx context.Context, req *dto.SearchProductRequest) ([]*entity.Product, int64, error) {
	query := `
		SELECT
			id, name, description, price, quantity, created_at, updated_at, COUNT(*) OVER(PARTITION BY 1)
		FROM
			products
		WHERE
			($1 = '' OR name ILIKE $1)
			AND ($2 = '' OR description ILIKE $2)
		LIMIT $3 OFFSET $4
	`

	name := fmt.Sprintf("%%%s%%", req.Name)
	description := fmt.Sprintf("%%%s%%", req.Description)
	rows, err := r.DB.QueryContext(ctx, query, name, description, req.Limit, pageutils.GetOffset(req.Page, req.Limit))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	total := int64(0)
	products := []*entity.Product{}
	for rows.Next() {
		product := new(entity.Product)
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Quantity, &product.CreatedAt,
			&product.UpdatedAt, &total); err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Save(ctx context.Context, product *entity.Product) error {
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
