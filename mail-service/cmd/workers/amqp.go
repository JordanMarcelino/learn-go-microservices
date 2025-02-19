package workers

import (
	"context"

	"github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/server"
)

func runAMQPWorker(cfg *config.Config, ctx context.Context) {
	srv := server.NewAMQPServer(cfg)
	go srv.Start()

	<-ctx.Done()
	srv.Shutdown()
}
