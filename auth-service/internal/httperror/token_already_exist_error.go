package httperror

import (
	"errors"
	"net/http"

	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
)

func NewTokenAlreadyExistError() *httperror.ResponseError {
	msg := constant.TokenAlreadyExistErrorMessage
	err := errors.New(msg)

	return httperror.NewResponseError(err, http.StatusTooManyRequests, msg)
}
