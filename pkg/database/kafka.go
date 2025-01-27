package database

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

type KafkaProducerOptions struct {
	SaramaConfig   *sarama.Config
	Brokers        []string
	Topic          string
	RetryMax       int
	FlushFrequency int
	Acks           sarama.RequiredAcks
	ReturnSuccess  bool
}

type KafkaConsumerOptions struct {
	SaramaConfig  *sarama.Config
	Brokers       []string
	ConsumerGroup string
	InitialOffset int64
}

func NewKafkaSyncProducer(opt *KafkaProducerOptions) *sarama.SyncProducer {
	config := opt.SaramaConfig
	if config == nil {
		config = sarama.NewConfig()
	}
	config.Producer.RequiredAcks = opt.Acks
	config.Producer.Retry.Max = opt.RetryMax
	config.Producer.Return.Successes = opt.ReturnSuccess
	config.Producer.Compression = sarama.CompressionGZIP
	config.Producer.Flush.Frequency = time.Duration(opt.FlushFrequency) * time.Millisecond

	producer, err := sarama.NewSyncProducer(opt.Brokers, config)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}

	return &producer
}

func NewKafkaAsyncProducer(opt *KafkaProducerOptions) *sarama.AsyncProducer {
	config := opt.SaramaConfig
	if config == nil {
		config = sarama.NewConfig()
	}
	config.Producer.RequiredAcks = opt.Acks
	config.Producer.Retry.Max = opt.RetryMax
	config.Producer.Return.Successes = opt.ReturnSuccess
	config.Producer.Compression = sarama.CompressionGZIP
	config.Producer.Flush.Frequency = time.Duration(opt.FlushFrequency) * time.Millisecond

	producer, err := sarama.NewAsyncProducer(opt.Brokers, config)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}

	return &producer
}

func NewKafkaConsumerGroup(opt *KafkaConsumerOptions) *sarama.ConsumerGroup {
	config := opt.SaramaConfig
	if config == nil {
		config = sarama.NewConfig()
	}
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = opt.InitialOffset

	consumerGroup, err := sarama.NewConsumerGroup(opt.Brokers, opt.ConsumerGroup, config)
	if err != nil {
		log.Fatalf("failed to create kafka consumer group: %v", err)
	}

	return &consumerGroup
}
