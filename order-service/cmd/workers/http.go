package workers

import (
	"context"
	"sync"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/server"
)

func runHttpWorker(cfg *config.Config, ctx context.Context, wg *sync.WaitGroup) {
	srv := server.NewHttpServer(cfg)
	go srv.Start()

	wg.Done()
	<-ctx.Done()
	srv.Shutdown()
}
