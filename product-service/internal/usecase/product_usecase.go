package usecase

import (
	"context"

	. "github.com/jordanmarcelino/learn-go-microservices/pkg/dto"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/pageutils"
	. "github.com/jordanmarcelino/learn-go-microservices/product-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/entity"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/httperror"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/repository"
)

type ProductUseCase interface {
	Search(ctx context.Context, req *SearchProductRequest) ([]*ProductResponse, *PageMetaData, error)
	Get(ctx context.Context, req *GetProductRequest) (*ProductResponse, error)
	Create(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error)
	Update(ctx context.Context, req *UpdateProductRequest) (*ProductResponse, error)
	Delete(ctx context.Context, req *DeleteProductRequest) error
}

type productUseCase struct {
	DataStore              repository.DataStore
	ProductCreatedProducer mq.KafkaProducer
	ProductUpdatedProducer mq.KafkaProducer
	ProductDeletedProducer mq.KafkaProducer
}

func NewProductUseCase(
	dataStore repository.DataStore,
	productCreatedProducer mq.KafkaProducer,
	productUpdatedProducer mq.KafkaProducer,
	productDeletedProducer mq.KafkaProducer,
) ProductUseCase {
	return &productUseCase{
		DataStore:              dataStore,
		ProductCreatedProducer: productCreatedProducer,
		ProductUpdatedProducer: productUpdatedProducer,
		ProductDeletedProducer: productDeletedProducer,
	}
}

func (u *productUseCase) Search(ctx context.Context, req *SearchProductRequest) ([]*ProductResponse, *PageMetaData, error) {
	productRepository := u.DataStore.ProductRepository()

	products, total, err := productRepository.Search(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	res := ToProductResponses(products)
	metadata := pageutils.NewMetadata(total, req.Page, req.Limit)
	return res, metadata, nil
}

func (u *productUseCase) Get(ctx context.Context, req *GetProductRequest) (*ProductResponse, error) {
	productRepository := u.DataStore.ProductRepository()

	product, err := productRepository.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, httperror.NewProductNotFoundError()
	}

	return ToProductResponse(product), nil
}

func (u *productUseCase) Create(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error) {
	res := new(ProductResponse)
	err := u.DataStore.Atomic(ctx, func(ds repository.DataStore) error {
		productRepository := ds.ProductRepository()

		product := &entity.Product{Name: req.Name, Description: req.Description, Price: req.Price, Quantity: req.Quantity}
		if err := productRepository.Save(ctx, product); err != nil {
			return err
		}

		if err := u.ProductCreatedProducer.Send(ctx, ToProductCreatedEvent(product)); err != nil {
			return err
		}

		res = ToProductResponse(product)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *productUseCase) Update(ctx context.Context, req *UpdateProductRequest) (*ProductResponse, error) {
	res := new(ProductResponse)
	err := u.DataStore.Atomic(ctx, func(ds repository.DataStore) error {
		productRepository := ds.ProductRepository()

		product, err := productRepository.FindByID(ctx, req.ID)
		if err != nil {
			return err
		}
		if product == nil {
			return httperror.NewProductNotFoundError()
		}

		product.Name = req.Name
		product.Description = req.Description
		product.Price = req.Price
		product.Quantity = req.Quantity

		if err := productRepository.Update(ctx, product); err != nil {
			return err
		}

		if err := u.ProductUpdatedProducer.Send(ctx, ToProductUpdatedEvent(product)); err != nil {
			return err
		}

		res = ToProductResponse(product)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *productUseCase) Delete(ctx context.Context, req *DeleteProductRequest) error {
	return u.DataStore.Atomic(ctx, func(ds repository.DataStore) error {
		productRepository := ds.ProductRepository()

		product, err := productRepository.FindByID(ctx, req.ID)
		if err != nil {
			return err
		}
		if product == nil {
			return httperror.NewProductNotFoundError()
		}

		if err := productRepository.DeleteByID(ctx, product.ID); err != nil {
			return err
		}

		return u.ProductDeletedProducer.Send(ctx, ToProductDeletedEvent(product))
	})
}
