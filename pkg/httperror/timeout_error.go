package httperror

import (
	"errors"
	"net/http"

	"github.com/jordanmarcelino/learn-go-microservices/pkg/constant"
)

func NewTimeoutError() *ResponseError {
	msg := constant.RequestTimeoutErrorMessage
	err := errors.New(msg)

	return NewResponseError(err, http.StatusRequestTimeout, msg)
}
