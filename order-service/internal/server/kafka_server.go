package server

import (
	"context"
	"time"

	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/config"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/log"
	"github.com/jordanmarcelino/learn-go-microservices/order-service/internal/provider"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
)

type KafkaServer struct {
	ctx       context.Context
	cancel    context.CancelFunc
	consumers []mq.KafkaConsumer
}

func NewKafkaServer(cfg *config.Config) *KafkaServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &KafkaServer{
		ctx:       ctx,
		cancel:    cancel,
		consumers: provider.BootstrapKafka(cfg.Kafka),
	}
}

func (s *KafkaServer) Start() {
	log.Logger.Info("Running Kafka server...")
	for _, consumer := range s.consumers {
		go consumer.Consume(s.ctx)
	}
	log.Logger.Info("Kafka server is running...")
}

func (s *KafkaServer) Shutdown() {
	log.Logger.Info("Attempting to shut down the Kafka server...")
	time.Sleep(10 * time.Second)

	s.cancel()
	for _, consumer := range s.consumers {
		if err := consumer.Close(); err != nil {
			log.Logger.Errorf("Error closing consumer %v: %v", consumer.Topic(), err)
		}
	}

	log.Logger.Info("Kafka server shut down gracefully")
}
