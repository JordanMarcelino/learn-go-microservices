package httperror

import (
	"errors"
	"net/http"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
)

func NewProductNotFoundError() *httperror.ResponseError {
	msg := constant.ProductNotFoundErrorMessage
	err := errors.New(msg)

	return httperror.NewResponseError(err, http.StatusNotFound, msg)
}
