package mq

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/log"
)

type ProductUpdatedProducer struct {
	Producer sarama.SyncProducer
	topic    string
}

func NewProductUpdatedProducer(producer sarama.SyncProducer) mq.KafkaProducer {
	return &ProductUpdatedProducer{
		Producer: producer,
		topic:    constant.ProductUpdatedTopic,
	}
}

func (p *ProductUpdatedProducer) Send(ctx context.Context, event mq.KafkaEvent) error {
	bytes, err := sonic.Marshal(event)
	if err != nil {
		return err
	}

	log.Logger.Infof("sending message to topic %s, value %s", p.topic, string(bytes))

	errCh := make(chan error)

	go func() {
		_, _, err := p.Producer.SendMessage(&sarama.ProducerMessage{
			Topic: p.topic,
			Key:   sarama.StringEncoder(event.ID()),
			Value: sarama.ByteEncoder(bytes),
		})
		errCh <- err
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return err
		}

		log.Logger.Infof("message sent to topic %s", p.topic)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *ProductUpdatedProducer) Topic() string {
	return p.topic
}
