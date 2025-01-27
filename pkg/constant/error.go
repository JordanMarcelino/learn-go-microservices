package constant

const (
	InternalServerErrorMessage     = "currently our server is facing unexpected error, please try again later"
	EOFErrorMessage                = "missing body request"
	JsonSyntaxErrorMessage         = "invalid JSON syntax"
	JsonUnMarshallTypeErrorMessage = "invalid value for %s"
	UnauthorizedErrorMessage       = "unauthorized"
	RequestTimeoutErrorMessage     = "failed to process request in time, please try again"
	ValidationErrorMessage         = "input validation error"
)
