package provider

import (
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/controller"
)

func BootstrapHttp(cfg *config.Config, router *gin.Engine) {
	appController := controller.NewAppController()
	appController.Route(router)

	BootstrapUser(router)
}
