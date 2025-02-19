package workers

import (
	"context"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/server"
)

func runKafkaWorker(cfg *config.Config, ctx context.Context) {
	srv := server.NewKafkaServer(cfg)
	go srv.Start()

	<-ctx.Done()
	srv.Shutdown()
}
