package httperror

import (
	"errors"
	"net/http"

	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/constant"
)

func NewProductNotFoundError() *httperror.ResponseError {
	msg := constant.ProductNotFoundErrorMessage
	err := errors.New(msg)

	return httperror.NewResponseError(err, http.StatusNotFound, msg)
}
