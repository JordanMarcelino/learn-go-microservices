package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/constant"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/ginutils"
)

func AuthMiddleware(ctx *gin.Context) {
	userID, ok := ginutils.GetXUserID(ctx)
	if !ok {
		ctx.Error(httperror.NewUnauthorizedError())
		ctx.Abort()
		return
	}

	email, ok := ginutils.GetXEmail(ctx)
	if !ok {
		ctx.Error(httperror.NewUnauthorizedError())
		ctx.Abort()
		return
	}

	ctx.Set(constant.CTX_USER_ID, userID)
	ctx.Set(constant.CTX_EMAIL, email)
}
