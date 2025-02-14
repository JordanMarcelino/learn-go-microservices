package httperror

import (
	"errors"
	"net/http"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
)

func NewInsufficientProductStockError() *httperror.ResponseError {
	msg := constant.InsufficientProductStockErrorMessage
	err := errors.New(msg)

	return httperror.NewResponseError(err, http.StatusConflict, msg)
}
