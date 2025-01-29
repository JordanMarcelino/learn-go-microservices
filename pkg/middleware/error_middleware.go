package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/constant"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/dto"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/httperror"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/validationutils"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		errLen := len(ctx.Errors)
		if errLen > 0 {
			err := ctx.Errors.Last()

			switch e := err.Err.(type) {
			case validator.ValidationErrors:
				handleValidationError(ctx, e)
			case *json.SyntaxError:
				handleJsonSyntaxError(ctx)
			case *json.UnmarshalTypeError:
				handleJsonUnmarshalTypeError(ctx, e)
			case *time.ParseError:
				handleParseTimeError(ctx, e)
			case *httperror.ResponseError:
				ctx.AbortWithStatusJSON(e.GetCode(), dto.WebResponse[any]{
					Message: e.DisplayMessage(),
				})
			default:
				if errors.Is(e, io.EOF) {
					ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
						Message: constant.EOFErrorMessage,
					})
					return
				}

				ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.WebResponse[any]{
					Message: constant.InternalServerErrorMessage,
				})
			}
		}
	}
}

func handleJsonSyntaxError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
		Message: constant.JsonSyntaxErrorMessage,
	})
}

func handleJsonUnmarshalTypeError(ctx *gin.Context, err *json.UnmarshalTypeError) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
		Message: fmt.Sprintf(constant.JsonUnMarshallTypeErrorMessage, err.Field),
	})
}

func handleParseTimeError(ctx *gin.Context, err *time.ParseError) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
		Message: fmt.Sprintf("please send time in format of %s, got: %s", constant.ConvertGoTimeLayoutToReadable(err.Layout), err.Value),
	})
}

func handleValidationError(ctx *gin.Context, err validator.ValidationErrors) {
	ve := []dto.FieldError{}

	for _, fe := range err {
		ve = append(ve, dto.FieldError{
			Field:   fe.Field(),
			Message: validationutils.TagToMsg(fe),
		})
	}

	ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
		Message: constant.ValidationErrorMessage,
		Errors:  ve,
	})
}
