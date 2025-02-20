package mq

import (
	"context"
	"math"
	"time"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"github.com/jordanmarcelino/learn-go-microservices/pkg/mq"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/constant"
	"github.com/jordanmarcelino/learn-go-microservices/product-service/internal/log"
)

type ProductCreatedProducer struct {
	Producer  sarama.AsyncProducer
	RetryChan chan *sarama.ProducerMessage
	topic     string
}

func NewProductCreatedProducer(producer sarama.AsyncProducer) mq.KafkaProducer {
	pr := &ProductCreatedProducer{
		Producer:  producer,
		topic:     constant.ProductCreatedTopic,
		RetryChan: make(chan *sarama.ProducerMessage, 100),
	}

	go pr.HandleErrors()
	go pr.Retry()
	return pr
}

func (p *ProductCreatedProducer) Send(ctx context.Context, event mq.KafkaEvent) error {
	bytes, err := sonic.Marshal(event)
	if err != nil {
		return err
	}

	select {
	case p.Producer.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(event.ID()),
		Value: sarama.ByteEncoder(bytes),
	}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *ProductCreatedProducer) HandleErrors() {
	for err := range p.Producer.Errors() {
		if err != nil {
			log.Logger.Errorf("failed to send message: %s", err)
		}

		p.RetryChan <- err.Msg
	}
}

func (p *ProductCreatedProducer) Retry() {
	for msg := range p.RetryChan {
		metadata := new(mq.KafkaMetadata)
		p.UnmarshalMetadata(msg, metadata)

		if metadata.Retry > constant.KafkaProducerRetryLimit {
			log.Logger.Errorf("failed to send message after %d retries [partition-%v]-[offset-%v] %v", constant.KafkaProducerRetryLimit, msg.Partition, msg.Offset, msg.Value)
			return
		}

		metadata.Retry++
		metaBytes, _ := sonic.Marshal(metadata)
		msg.Metadata = metaBytes

		backoff := time.Duration(math.Pow(constant.KafkaProducerRetryDelay, float64(metadata.Retry))) * time.Second
		time.Sleep(backoff)

		select {
		case p.Producer.Input() <- msg:
			log.Logger.Infof("retrying message (attempt %d) to topic %s", metadata.Retry, msg.Topic)
		default:
			log.Logger.Infof("failed to retry message (attempt %d) to topic %s", metadata.Retry, msg.Topic)
		}
	}
}

func (p *ProductCreatedProducer) UnmarshalMetadata(msg *sarama.ProducerMessage, metadata *mq.KafkaMetadata) {
	if msg.Metadata == nil {
		return
	}

	metaBytes, ok := msg.Metadata.([]byte)
	if !ok {
		return
	}

	_ = sonic.Unmarshal(metaBytes, &metadata)
}

func (p *ProductCreatedProducer) Topic() string {
	return p.topic
}
