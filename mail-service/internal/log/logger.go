package log

import "github.com/jordanmarcelino/learn-go-microservices/pkg/logger"

var Logger logger.Logger

func SetLogger(logger logger.Logger) {
	Logger = logger
}
