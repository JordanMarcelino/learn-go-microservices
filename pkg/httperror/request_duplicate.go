package httperror

import (
	"errors"
	"net/http"

	"github.com/jordanmarcelino/learn-go-microservices/pkg/constant"
)

func NewRequestDuplicateError() *ResponseError {
	msg := constant.RequestDuplicateErrorMessage
	err := errors.New(msg)

	return NewResponseError(err, http.StatusConflict, msg)
}
