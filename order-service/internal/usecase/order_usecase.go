package usecase

import (
	"context"

	"github.com/bsm/redislock"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"
	. "github.com/jordanmarcelino/learn-go-microservices/order-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/entity"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/httperror"
	. "github.com/jordanmarcelino/learn-go-microservices/pkg/dto"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/pageutils"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/repository"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/utils/redisutils"
	"github.com/shopspring/decimal"
)

type OrderUseCase interface {
	Search(ctx context.Context, req *SearchOrderRequest) ([]*OrderResponse, *PageMetaData, error)
	Get(ctx context.Context, req *GetOrderRequest) (*OrderResponse, error)
	Save(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, error)
}

type orderUseCaseImpl struct {
	DataStore               repository.DataStore
	LockRepository          repository.LockRepository
	OrderCreatedProducer    mq.KafkaProducer
	PaymentReminderProducer mq.AMQPProducer
	AutoCancelProducer      mq.AMQPProducer
}

func NewOrderUseCase(
	dataStore repository.DataStore,
	lockRepository repository.LockRepository,
	orderCreatedProducer mq.KafkaProducer,
	paymentReminderProducer mq.AMQPProducer,
	autoCancelProducer mq.AMQPProducer,
) *orderUseCaseImpl {
	return &orderUseCaseImpl{
		DataStore:               dataStore,
		LockRepository:          lockRepository,
		OrderCreatedProducer:    orderCreatedProducer,
		PaymentReminderProducer: paymentReminderProducer,
		AutoCancelProducer:      autoCancelProducer,
	}
}

func (u *orderUseCaseImpl) Search(ctx context.Context, req *SearchOrderRequest) ([]*OrderResponse, *PageMetaData, error) {
	orderRepository := u.DataStore.OrderRepository()

	orders, total, err := orderRepository.Search(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	res := ToOrderResponses(orders)
	metadata := pageutils.NewMetadata(total, req.Page, req.Limit)
	return res, metadata, nil
}

func (u *orderUseCaseImpl) Get(ctx context.Context, req *GetOrderRequest) (*OrderResponse, error) {
	order, err := u.DataStore.OrderRepository().FindByID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, httperror.NewOrderNotFoundError()
	}

	return ToOrderResponse(order), nil
}

func (u *orderUseCaseImpl) Save(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, error) {
	lockKey := redisutils.NewLockKey(req.RequestID, req.CustomerID)
	ttl := constant.CreateOrderTTL
	opt := &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(constant.CreateOrderRetryInterval), constant.CreateOrderRetryLimit),
	}

	lock, err := u.LockRepository.Get(ctx, lockKey, ttl, opt)
	if err != nil {
		return nil, err
	}
	defer lock.Release(ctx)

	res := new(OrderResponse)
	err = u.DataStore.Atomic(ctx, func(ds repository.DataStore) error {
		orderRepository := ds.OrderRepository()
		productRepository := ds.ProductRepository()

		order, err := orderRepository.FindByRequestID(ctx, req.RequestID)
		if err != nil {
			return err
		}
		if order != nil && order.CustomerID == req.CustomerID {
			res = ToOrderResponse(order)
			return nil
		}

		order = &entity.Order{
			RequestID:   req.RequestID,
			CustomerID:  req.CustomerID,
			Description: req.Description,
			Status:      constant.ORDER_PENDING,
			Items:       []*entity.OrderItem{},
		}
		productIds := []int64{}
		for _, item := range req.Items {
			productIds = append(productIds, item.ProductID)
		}

		products, err := productRepository.FindAllByIDForUpdate(ctx, productIds)
		if err != nil {
			return err
		}
		if len(products) != len(productIds) {
			return httperror.NewProductNotFoundError()
		}

		totalAmount := decimal.Decimal{}
		for i, product := range products {
			if product.Quantity < req.Items[i].Quantity {
				return httperror.NewInsufficientProductStockError()
			}
			product.Quantity -= req.Items[i].Quantity

			totalAmount = totalAmount.Add(product.Price.Mul(decimal.NewFromInt(int64(req.Items[i].Quantity))))
			order.Items = append(order.Items, &entity.OrderItem{ProductID: product.ID, Price: product.Price, Quantity: req.Items[i].Quantity})
		}
		order.TotalAmount = totalAmount

		if err := productRepository.UpdateAllQuantity(ctx, products); err != nil {
			return err
		}
		if err := orderRepository.Save(ctx, order); err != nil {
			return err
		}

		if err := u.OrderCreatedProducer.Send(ctx, ToOrderCreatedEvent(order)); err != nil {
			return err
		}

		if err := u.PaymentReminderProducer.Send(ctx, &PaymentReminderEvent{OrderID: order.ID, Email: req.CustomerEmail}); err != nil {
			return err
		}

		if err := u.AutoCancelProducer.Send(ctx, &AutoCancelEvent{OrderID: order.ID, UserID: req.CustomerID, Email: req.CustomerEmail}); err != nil {
			return err
		}

		res = ToOrderResponse(order)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
