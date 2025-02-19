package workers

import (
	"context"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/server"
)

func runAMQPWorker(cfg *config.Config, ctx context.Context) {
	srv := server.NewAMQPServer(cfg)
	go srv.Start()

	<-ctx.Done()
	srv.Shutdown()
}
