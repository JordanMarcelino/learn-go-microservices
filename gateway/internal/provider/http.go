package provider

import (
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/gateway/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/gateway/internal/controller"
)

func BootstrapHttp(cfg *config.Config, router *gin.Engine) {
	appController := controller.NewAppController()
	gatewayController := controller.NewGatewayController(authMiddleware, cfg.ServiceConfig)

	appController.Route(router)
	gatewayController.Route(router)
}
