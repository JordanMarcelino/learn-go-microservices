package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/constant"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/ginutils"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/pageutils"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/middleware"
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
	g := r.Use(middleware.AuthMiddleware)
	{
		g.GET("", c.Search)
		g.POST("", c.Create)
	}
}

func (c *ProductController) Create(ctx *gin.Context) {
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

func (c *ProductController) Search(ctx *gin.Context) {
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", fmt.Sprintf("%d", constant.DefaultLimit)), 10, 64)
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", fmt.Sprintf("%d", constant.DefaultPage)), 10, 64)

	req := &dto.SearchProductRequest{Limit: limit, Page: page}
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.Error(err)
		return
	}

	res, paging, err := c.productUseCase.Search(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	paging.Links = pageutils.NewLinks(ctx.Request, int(paging.Page), int(paging.Size), int(paging.TotalItem), int(paging.TotalPage))
	ginutils.ResponseOKPagination(ctx, res, paging)
}
