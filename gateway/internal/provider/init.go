package provider

import (
	"github.com/jordanmarcelino/learn-go-microservices/pkg/middleware"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/utils/jwtutils"
)

var (
	jwtUtil        jwtutils.JwtUtil
	authMiddleware *middleware.AuthMiddleware
)

func BootstrapGlobal() {
	jwtUtil = jwtutils.NewJwtUtil()
	authMiddleware = middleware.NewAuthMiddleware(jwtUtil)
}
