package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
)

func RequestTimeout(timeoutPeriod int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		timeoutCtx, cancel := context.WithTimeout(
			ctx.Request.Context(),
			time.Duration(timeoutPeriod)*time.Second,
		)
		defer cancel()

		done := make(chan struct{})
		ctx.Request = ctx.Request.WithContext(timeoutCtx)

		go next(ctx, done)

		select {
		case <-timeoutCtx.Done():
			ctx.Error(httperror.NewTimeoutError())
			ctx.Abort()
		case <-done:
		}
	}
}

func next(ctx *gin.Context, done chan struct{}) {
	defer func() {
		close(done)
		if err, ok := recover().(error); ok && err != nil {
			ctx.Error(err)
			ctx.Abort()
		}
	}()

	ctx.Next()
}
