package workers

import (
	"context"

	"github.com/jordanmarcelino/learn-go-microservices/mail-service/internal/server"
)

func runAMQPWorker(ctx context.Context) {
	srv := server.NewAMQPServer()
	go srv.Start()

	<-ctx.Done()
	srv.Shutdown()
}
