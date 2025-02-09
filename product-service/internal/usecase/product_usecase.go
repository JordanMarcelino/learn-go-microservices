package usecase

import (
	"context"

	. "github.com/jordanmarcelino/learn-go-microservices/pkg/dto"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/pageutils"
	. "github.com/jordanmarcelino/learn-go-microservices/product-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/entity"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/repository"
)

type ProductUseCase interface {
	Search(ctx context.Context, req *SearchProductRequest) ([]*ProductResponse, *PageMetaData, error)
	Create(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error)
}

type productUseCase struct {
	DataStore              repository.DataStore
	ProductCreatedProducer mq.KafkaProducer
}

func NewProductUseCase(
	dataStore repository.DataStore,
	ProductCreatedProducer mq.KafkaProducer,
) ProductUseCase {
	return &productUseCase{
		DataStore:              dataStore,
		ProductCreatedProducer: ProductCreatedProducer,
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
