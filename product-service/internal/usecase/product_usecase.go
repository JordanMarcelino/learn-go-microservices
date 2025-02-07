package usecase

import (
	"context"

	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/entity"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/repository"
)

type ProductUseCase interface {
	Create(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error)
}

type productUseCase struct {
	DataStore repository.DataStore
}

func NewProductUseCase(
	dataStore repository.DataStore,
) ProductUseCase {
	return &productUseCase{
		DataStore: dataStore,
	}
}

func (u *productUseCase) Create(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	res := new(dto.ProductResponse)
	err := u.DataStore.Atomic(ctx, func(ds repository.DataStore) error {
		productRepository := ds.ProductRepository()

		product := &entity.Product{Name: req.Name, Description: req.Description, Price: req.Price, Quantity: req.Quantity}
		if err := productRepository.Save(ctx, product); err != nil {
			return err
		}

		res = dto.ToProductResponse(product)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}
