package ginutils

import (
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/constant"
)

func GetUserID(ctx *gin.Context) int64 {
	return ctx.GetInt64(constant.CTX_USER_ID)
}

func GetEmail(ctx *gin.Context) string {
	return ctx.GetString(constant.CTX_EMAIL)
}
