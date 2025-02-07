package ginutils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/constant"
)

func GetXUserID(ctx *gin.Context) (int64, bool) {
	xUserID := ctx.GetHeader(constant.X_USER_ID)
	userID, err := strconv.ParseInt(xUserID, 10, 64)
	if err != nil {
		return 0, false
	}
	return userID, true
}
