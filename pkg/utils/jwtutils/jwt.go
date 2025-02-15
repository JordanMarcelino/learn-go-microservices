package jwtutils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/config"
)

type JwtUtil interface {
	Sign(payload *JWTPayload) (string, error)
	Parse(token string) (*JWTClaims, error)
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
}

type JWTPayload struct {
	UserID int64
	Email  string
}

type jwtUtil struct {
	config *config.JwtConfig
}

func NewJwtUtil() JwtUtil {
	return &jwtUtil{
		config: config.JWT_CONFIG,
	}
}

func (j *jwtUtil) Sign(payload *JWTPayload) (string, error) {
	currentTime := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID: payload.UserID,
		Email:  payload.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(j.config.TokenDuration) * time.Minute)),
			Issuer:    j.config.Issuer,
		},
	})

	s, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return "", err
	}

	return s, nil
}

func (j *jwtUtil) Parse(token string) (*JWTClaims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods(j.config.AllowedAlgs),
		jwt.WithIssuer(j.config.Issuer),
		jwt.WithIssuedAt(),
	)

	return j.parseClaims(parser, token)
}

func (j *jwtUtil) parseClaims(parser *jwt.Parser, token string) (*JWTClaims, error) {
	parsedToken, err := parser.ParseWithClaims(token, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.config.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(*JWTClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	return nil, errors.New("token is not valid")
}
