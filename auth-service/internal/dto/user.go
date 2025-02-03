package dto

import (
	"time"

	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/entity"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterResponse struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type VerificationResponse struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type VerificationRequest struct {
	Email             string `json:"email" binding:"required,email"`
	VerificationToken string `json:"verification_token" binding:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func ToRegisterResponse(user *entity.User) *RegisterResponse {
	return &RegisterResponse{
		ID:         user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func ToVerificationResponse(user *entity.User) *VerificationResponse {
	return &VerificationResponse{
		ID:         user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}
}
