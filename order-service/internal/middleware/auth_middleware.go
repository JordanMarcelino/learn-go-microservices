package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/ginutils"
)

func AuthMiddleware(ctx *gin.Context) {
	if _, ok := ginutils.GetXUserID(ctx); !ok {
		ctx.Error(httperror.NewUnauthorizedError())
		ctx.Abort()
		return
	}
}
