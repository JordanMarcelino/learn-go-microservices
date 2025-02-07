package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/ginutils"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/usecase"
)

type ProductController struct {
	productUseCase usecase.ProductUseCase
}

func NewProductController(productUseCase usecase.ProductUseCase) *ProductController {
	return &ProductController{
		productUseCase: productUseCase,
	}
}

func (c *ProductController) Route(r *gin.Engine) {
	r.POST("", c.Create)
}

func (c *ProductController) Create(ctx *gin.Context) {
	if _, ok := ginutils.GetXUserID(ctx); !ok {
		ctx.Error(httperror.NewUnauthorizedError())
		return
	}

	req := new(dto.CreateProductRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.productUseCase.Create(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}
	ginutils.ResponseCreated(ctx, res)
}
