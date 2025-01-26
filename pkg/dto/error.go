package dto

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
