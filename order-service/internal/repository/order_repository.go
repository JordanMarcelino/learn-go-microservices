package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/entity"
)

type OrderRepository interface {
	FindByRequestID(ctx context.Context, requestID string) (*entity.Order, error)
	Save(ctx context.Context, order *entity.Order) error
}

type orderRepositoryImpl struct {
	DB DBTX
}

func NewOrderRepository(db DBTX) *orderRepositoryImpl {
	return &orderRepositoryImpl{
		DB: db,
	}
}

func (r *orderRepositoryImpl) FindByRequestID(ctx context.Context, requestID string) (*entity.Order, error) {
	query := `
		SELECT
			id, customer_id, total_amount, description, status, created_at, updated_at
		FROM
			orders
		WHERE
			request_id = $1
	`

	order := &entity.Order{RequestID: requestID}
	err := r.DB.QueryRowContext(ctx, query, requestID).Scan(&order.ID, &order.CustomerID, &order.TotalAmount, &order.Description,
		&order.Status, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	query = `
		SELECT
			ot.id, ot.product_id, ot.quantity, ot.price
		FROM
			order_items ot
		JOIN
			orders o
		ON
			ot.order_id = o.id
		WHERE
			o.request_id = $1
	`

	rows, err := r.DB.QueryContext(ctx, query, order.RequestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*entity.OrderItem{}
	for rows.Next() {
		item := new(entity.OrderItem)
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	order.Items = items

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return order, nil
}

func (r *orderRepositoryImpl) Save(ctx context.Context, order *entity.Order) error {
	query := `
		INSERT INTO
			orders (request_id, customer_id, total_amount, description, status)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING
			id, created_at, updated_at
	`

	if err := r.DB.QueryRowContext(ctx, query, order.RequestID, order.CustomerID, order.TotalAmount, order.Description, order.Status).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt); err != nil {
		return err
	}

	params := []string{}
	args := []any{}

	orderItemsLen := len(params)
	for _, item := range order.Items {
		params = append(params, fmt.Sprintf("($%d, $%d, $%d, $%d)", orderItemsLen+1, orderItemsLen+2, orderItemsLen+3, orderItemsLen+4))
		args = append(args, item.ProductID, order.ID, item.Quantity, item.Price)
	}

	query = fmt.Sprintf(
		"INSERT INTO order_items (product_id, order_id, quantity, price) VALUES %s RETURNING id, created_at, updated_at",
		strings.Join(params, ","),
	)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		item := order.Items[i]
		if err := rows.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
