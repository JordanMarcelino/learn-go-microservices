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
	r.POST("/resend-verification", c.ResendVerification)
}

func (c *UserController) Login(ctx *gin.Context) {
	req := new(dto.LoginRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.userUseCase.Login(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}
	ginutils.ResponseOK(ctx, res)
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
	req := new(dto.VerificationRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	res, err := c.userUseCase.Verify(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOK(ctx, res)
}

func (c *UserController) ResendVerification(ctx *gin.Context) {
	req := new(dto.ResendVerificationRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.Error(err)
		return
	}

	err := c.userUseCase.ResendVerification(ctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ginutils.ResponseOKPlain(ctx)
}
