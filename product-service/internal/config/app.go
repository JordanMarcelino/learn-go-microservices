package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Environment string `mapstructure:"APP_ENVIRONMENT"`
}

func initAppConfig() *AppConfig {
	appConfig := &AppConfig{}

	if err := viper.Unmarshal(&appConfig); err != nil {
		log.Fatalf("error mapping app config: %v", err)
	}

	return appConfig
}
