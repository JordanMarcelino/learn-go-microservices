package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/dto"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/usecase"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/ginutils"
)

type UserController struct {
	userUseCase usecase.UserUseCase
}

func NewUserController(userUseCase usecase.UserUseCase) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

func (c *UserController) Route(r *gin.Engine) {
	r.POST("/login", c.Login)
	r.POST("/register", c.Register)
	r.POST("/verify", c.Verify)
}

func (c *UserController) Login(ctx *gin.Context) {

}

func (c *UserController) Register(ctx *gin.Context) {
	req := new(dto.RegisterRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.userUseCase.Register(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseCreated(ctx, res)
}

func (c *UserController) Verify(ctx *gin.Context) {

}
