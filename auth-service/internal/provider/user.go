package provider

import (
	"github.com/gin-gonic/gin"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/controller"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/mq"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/repository"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/usecase"
)

func BootstrapUser(router *gin.Engine) {
	redisRepository := repository.NewRedisRepository(rdb)
	sendVerificationProducer := mq.NewSendVerificationProducer(rabbitmq)
	accountVerifiedProducer := mq.NewAccountVerifiedProducer(rabbitmq)

	userUseCase := usecase.NewUserUseCase(bcryptHasher, jwtUtil, dataStore, redisRepository, sendVerificationProducer, accountVerifiedProducer)
	userController := controller.NewUserController(userUseCase)

	userController.Route(router)
}
