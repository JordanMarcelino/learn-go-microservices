package config

import (
	"log"

	"github.com/spf13/viper"
)

type AMQPConfig struct {
	Username string `mapstructure:"AMQP_USERNAME"`
	Password string `mapstructure:"AMQP_PASSWORD"`
	Host     string `mapstructure:"AMQP_HOST"`
	VHost    string `mapstructure:"AMQP_VHOST"`
	Port     int    `mapstructure:"AMQP_PORT"`
}

func initAmqpConfig() *AMQPConfig {
	amqpConfig := &AMQPConfig{}

	if err := viper.Unmarshal(&amqpConfig); err != nil {
		log.Fatalf("error mapping amqp config: %v", err)
	}

	return amqpConfig
}
