package config

import (
	"log"

	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Brokers        []string `mapstructure:"KAFKA_BROKERS"`
	RetryMax       int      `mapstructure:"KAFKA_RETRY_MAX"`
	FlushFrequency int      `mapstructure:"KAFKA_FLUSH_FREQUENCY"`
	ReturnSuccess  bool     `mapstructure:"KAFKA_RETURN_SUCCESS"`
}

func initKafkaConfig() *KafkaConfig {
	kafkaConfig := &KafkaConfig{}

	if err := viper.Unmarshal(&kafkaConfig); err != nil {
		log.Fatalf("error mapping kafka config: %v", err)
	}

	return kafkaConfig
}
