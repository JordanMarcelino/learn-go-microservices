package httperror

import (
	"errors"
	"net/http"

	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
)

func NewUserNotFoundError() *httperror.ResponseError {
	msg := constant.UserNotFoundErrorErrorMessage
	err := errors.New(msg)

	return httperror.NewResponseError(err, http.StatusBadRequest, msg)
}
