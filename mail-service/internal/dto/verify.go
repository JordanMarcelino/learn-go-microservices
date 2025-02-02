package dto

import "github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/constant"

type SendVerificationEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func (e SendVerificationEvent) Key() string {
	return constant.SendVerificationKey
}

type AccountVerifiedEvent struct {
	Email string `json:"email"`
}

func (e AccountVerifiedEvent) Key() string {
	return constant.AccountVerifiedKey
}
