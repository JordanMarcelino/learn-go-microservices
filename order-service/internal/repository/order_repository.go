package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/entity"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/pageutils"
)

type OrderRepository interface {
	Search(ctx context.Context, req *dto.SearchOrderRequest) ([]*entity.Order, int64, error)
	FindByRequestID(ctx context.Context, requestID string) (*entity.Order, error)
	FindByID(ctx context.Context, id int64) (*entity.Order, error)
	Save(ctx context.Context, order *entity.Order) error
	UpdateStatus(ctx context.Context, order *entity.Order) error
}

type orderRepositoryImpl struct {
	DB DBTX
}

func NewOrderRepository(db DBTX) *orderRepositoryImpl {
	return &orderRepositoryImpl{
		DB: db,
	}
}

func (r *orderRepositoryImpl) Search(ctx context.Context, req *dto.SearchOrderRequest) ([]*entity.Order, int64, error) {
	query := `
		SELECT
			o.id, o.customer_id, o.total_amount, o.description, o.status, o.created_at, o.updated_at,
			ot.id, ot.product_id, ot.quantity, ot.price,
			COUNT(o.id) OVER(PARTITION BY 1)
		FROM
			orders o
		JOIN
			order_items ot
		ON
			o.id = ot.order_id
		WHERE
			($1 = '' OR o.status = $1)
			AND o.created_at BETWEEN $2 AND $3
		LIMIT $4 OFFSET $5
	`

	rows, err := r.DB.QueryContext(ctx, query, req.Status, req.StartDate, req.EndDate, req.Limit, pageutils.GetOffset(req.Page, req.Limit))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	total := int64(0)
	orderMap := map[int64]*entity.Order{}
	for rows.Next() {
		order := new(entity.Order)
		item := new(entity.OrderItem)
		if err := rows.Scan(&order.ID, &order.CustomerID, &order.TotalAmount, &order.Description, &order.Status, &order.CreatedAt, &order.UpdatedAt,
			&item.ID, &item.ProductID, &item.Quantity, &item.Price, &total); err != nil {
			return nil, 0, err
		}

		if _, ok := orderMap[order.ID]; !ok {
			orderMap[order.ID] = order
		}
		orderMap[order.ID].Items = append(orderMap[order.ID].Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	orders := []*entity.Order{}
	for _, order := range orderMap {
		orders = append(orders, order)
	}

	return orders, total, nil
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

func (r *orderRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.Order, error) {
	query := `
		SELECT
			request_id, customer_id, total_amount, description, status, created_at, updated_at
		FROM
			orders
		WHERE
			id = $1
	`

	order := &entity.Order{ID: id}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&order.RequestID, &order.CustomerID, &order.TotalAmount, &order.Description,
		&order.Status, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	query = `
		SELECT
			id, product_id, quantity, price
		FROM
			order_items
		WHERE
			order_id = $1
	`

	rows, err := r.DB.QueryContext(ctx, query, id)
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

func (r *orderRepositoryImpl) UpdateStatus(ctx context.Context, order *entity.Order) error {
	query := `
		UPDATE
			orders
		SET
			status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $2
	`

	_, err := r.DB.ExecContext(ctx, query, order.Status, order.ID)
	return err
}
