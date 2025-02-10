package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/constant"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
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
		g.POST("", c.Create)
		g.PUT("/:productId", c.Update)
		g.DELETE("/:productId", c.Delete)
	}

	r.GET("", c.Search)
	r.GET("/:productId", c.Get)
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

func (c *ProductController) Get(ctx *gin.Context) {
	param := ctx.Param("productId")
	productId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		ctx.Error(httperror.NewInvalidURLParamError(param))
		return
	}

	req := &dto.GetProductRequest{ID: productId}
	res, err := c.productUseCase.Get(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOK(ctx, res)
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

func (c *ProductController) Update(ctx *gin.Context) {
	param := ctx.Param("productId")
	productId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		ctx.Error(httperror.NewInvalidURLParamError(param))
		return
	}

	req := &dto.UpdateProductRequest{ID: productId}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.productUseCase.Update(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOK(ctx, res)
}

func (c *ProductController) Delete(ctx *gin.Context) {
	param := ctx.Param("productId")
	productId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		ctx.Error(httperror.NewInvalidURLParamError(param))
		return
	}

	req := &dto.DeleteProductRequest{ID: productId}
	if err := c.productUseCase.Delete(ctx, req); err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOKPlain(ctx)
}
