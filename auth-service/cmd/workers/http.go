package workers

import (
	"context"

	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/auth-service/internal/server"
)

func runHttpWorker(cfg *config.Config, ctx context.Context) {
	srv := server.NewHttpServer(cfg)
	go srv.Start()

	<-ctx.Done()
	srv.Shutdown()
}
