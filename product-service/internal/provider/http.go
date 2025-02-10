package provider

import (
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/controller"
)

func BootstrapHttp(cfg *config.Config, router *gin.Engine) {
	appController := controller.NewAppController()
	appController.Route(router)

	BootstrapProduct(cfg, router)
}
