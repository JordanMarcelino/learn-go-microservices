package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/gateway/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/gateway/internal/proxy"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/middleware"
)

type GatewayController struct {
	authMiddleware *middleware.AuthMiddleware
	serviceCfg     *config.ServiceConfig
}

func NewGatewayController(authMiddleware *middleware.AuthMiddleware, serviceCfg *config.ServiceConfig) *GatewayController {
	return &GatewayController{
		authMiddleware: authMiddleware,
		serviceCfg:     serviceCfg,
	}
}

func (c *GatewayController) Route(r *gin.Engine) {
	authProxy := proxy.NewReverseProxy(c.serviceCfg.AuthURL)
	productProxy := proxy.NewReverseProxy(c.serviceCfg.ProductURL)
	orderProxy := proxy.NewReverseProxy(c.serviceCfg.OrderURL)
	mailProxy := proxy.NewReverseProxy(c.serviceCfg.MailURL)

	r.Any("/auth/*path", authProxy)

	protected := r.Group("", c.authMiddleware.Authorization())
	{
		protected.Any("/products/*path", productProxy)
		protected.Any("/orders/*path", orderProxy)
		protected.Any("/mail/*path", mailProxy)
	}
}
