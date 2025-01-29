package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	AuthURL    string `mapstructure:"AUTH_SERVICE_URL"`
	ProductURL string `mapstructure:"PRODUCT_SERVICE_URL"`
	OrderURL   string `mapstructure:"ORDER_SERVICE_URL"`
	MailURL    string `mapstructure:"MAIL_SERVICE_URL"`
}

func initServiceConfig() *ServiceConfig {
	serviceConfig := &ServiceConfig{}

	if err := viper.Unmarshal(&serviceConfig); err != nil {
		log.Fatalf("error mapping service config: %v", err)
	}

	return serviceConfig
}
