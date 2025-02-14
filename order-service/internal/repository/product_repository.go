package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/entity"
)

type ProductRepository interface {
	FindAllByIDForUpdate(ctx context.Context, ids []int64) ([]*entity.Product, error)
	FindByID(ctx context.Context, id int64) (*entity.Product, error)
	Save(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	UpdateAllQuantity(ctx context.Context, products []*entity.Product) error
	DeleteByID(ctx context.Context, id int64) error
}

type productRepository struct {
	DB DBTX
}

func NewProductRepository(db DBTX) ProductRepository {
	return &productRepository{
		DB: db,
	}
}

func (r *productRepository) FindAllByIDForUpdate(ctx context.Context, ids []int64) ([]*entity.Product, error) {
	query := `
		SELECT
			id, name, description, price, quantity, created_at, updated_at
		FROM
			products
		WHERE
			id = ANY($1)
		FOR UPDATE
	`

	rows, err := r.DB.QueryContext(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*entity.Product{}
	for rows.Next() {
		product := new(entity.Product)
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price,
			&product.Quantity, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) FindByID(ctx context.Context, id int64) (*entity.Product, error) {
	query := `
		SELECT
			name, description, price, quantity, created_at, updated_at
		FROM
			products
		WHERE
			id = $1
	`

	product := &entity.Product{ID: id}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&product.Name, &product.Description, &product.Price, &product.Quantity,
		&product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return product, nil
}

func (r *productRepository) Save(ctx context.Context, product *entity.Product) error {
	query := `
		INSERT INTO
			products (id, name, description, price, quantity)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			created_at, updated_at
	`

	return r.DB.QueryRowContext(ctx, query, product.ID, product.Name, product.Description, product.Price, product.Quantity).Scan(&product.CreatedAt, &product.UpdatedAt)
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	query := `
		UPDATE
			products
		SET
			name = $1, description = $2, price = $3, quantity = $4, updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $5
		RETURNING
			created_at, updated_at
	`

	return r.DB.QueryRowContext(ctx, query, product.Name, product.Description, product.Price, product.Quantity, product.ID).Scan(&product.CreatedAt, &product.UpdatedAt)
}

func (r *productRepository) UpdateAllQuantity(ctx context.Context, products []*entity.Product) error {
	query := `
		INSERT INTO
			products (id, name, description, price, quantity)
		VALUES
			%s
		ON CONFLICT (id) DO UPDATE
		SET
			quantity = EXCLUDED.quantity,
			updated_at = CURRENT_TIMESTAMP
	`

	params := []string{}
	args := []any{}
	for i, product := range products {
		params = append(params, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
		args = append(args, product.ID, product.Name, product.Description, product.Price, product.Quantity)
	}

	query = fmt.Sprintf(query, strings.Join(params, ","))
	_, err := r.DB.ExecContext(ctx, query, args...)

	return err
}

func (r *productRepository) DeleteByID(ctx context.Context, id int64) error {
	query := `
		DELETE FROM
			products
		WHERE
			id = $1
	`

	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}
