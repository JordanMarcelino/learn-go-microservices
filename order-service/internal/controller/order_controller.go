package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/middleware"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/usecase"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/ginutils"
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
		g.POST("", c.Create)
	}
}

func (c *OrderController) Create(ctx *gin.Context) {
	req := &dto.CreateOrderRequest{CustomerID: ginutils.GetUserID(ctx)}
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
