package httperror

import (
	"errors"
	"net/http"

	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
)

func NewUserNotVerifiedError() *httperror.ResponseError {
	msg := constant.UserNotVerifiedErrorMessage
	err := errors.New(msg)

	return httperror.NewResponseError(err, http.StatusUnauthorized, msg)
}
