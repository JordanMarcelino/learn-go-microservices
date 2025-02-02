package config

import (
	"log"

	"github.com/spf13/viper"
)

type LoggerConfig struct {
	Level int `mapstructure:"LOGGER_LEVEL"`
}

func initLoggerConfig() *LoggerConfig {
	loggerConfig := &LoggerConfig{}

	if err := viper.Unmarshal(&loggerConfig); err != nil {
		log.Fatalf("error mapping logger config: %v", err)
	}

	return loggerConfig
}
