package httperror

import (
	"errors"
	"net/http"

	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
)

func NewTokenExpiredError() *httperror.ResponseError {
	msg := constant.TokenExpiredErrorErrorMessage
	err := errors.New(msg)

	return httperror.NewResponseError(err, http.StatusBadRequest, msg)
}
