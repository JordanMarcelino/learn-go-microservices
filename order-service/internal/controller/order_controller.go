package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/middleware"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/usecase"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/ginutils"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/pageutils"
)

type OrderController struct {
	OrderUseCase usecase.OrderUseCase
}

func NewOrderController(orderUseCase usecase.OrderUseCase) *OrderController {
	return &OrderController{
		OrderUseCase: orderUseCase,
	}
}

func (c *OrderController) Route(r *gin.Engine) {
	g := r.Use(middleware.AuthMiddleware)
	{
		g.GET("", c.Search)
		g.POST("", c.Create)
		g.GET("/:orderId", c.Get)

		g.POST("/pay", c.Pay)
		g.POST("/cancel", c.Cancel)
	}

}

func (c *OrderController) Search(ctx *gin.Context) {
	req := new(dto.SearchOrderRequest)
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.Error(err)
		return
	}

	res, paging, err := c.OrderUseCase.Search(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	paging.Links = pageutils.NewLinks(ctx.Request, int(paging.Page), int(paging.Size), int(paging.TotalItem), int(paging.TotalPage))
	ginutils.ResponseOKPagination(ctx, res, paging)
}

func (c *OrderController) Get(ctx *gin.Context) {
	param := ctx.Param("orderId")
	orderId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		ctx.Error(httperror.NewInvalidURLParamError(param))
		return
	}

	req := &dto.GetOrderRequest{OrderID: orderId}
	res, err := c.OrderUseCase.Get(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOK(ctx, res)
}

func (c *OrderController) Create(ctx *gin.Context) {
	req := &dto.CreateOrderRequest{CustomerID: ginutils.GetUserID(ctx), CustomerEmail: ginutils.GetEmail(ctx)}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.OrderUseCase.Save(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseCreated(ctx, res)
}

func (c *OrderController) Pay(ctx *gin.Context) {
	req := &dto.PayOrderRequest{CustomerID: ginutils.GetUserID(ctx), CustomerEmail: ginutils.GetEmail(ctx)}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.OrderUseCase.Pay(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOK(ctx, res)
}

func (c *OrderController) Cancel(ctx *gin.Context) {
	req := &dto.CancelOrderRequest{CustomerID: ginutils.GetUserID(ctx), CustomerEmail: ginutils.GetEmail(ctx)}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.OrderUseCase.Cancel(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOK(ctx, res)
}
