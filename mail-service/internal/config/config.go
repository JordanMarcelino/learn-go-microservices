package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Amqp   *AMQPConfig
	SMTP   *SMTPConfig
	Logger *LoggerConfig
}

func InitConfig() *Config {
	configPath := parseConfigPath()
	viper.AddConfigPath(configPath)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	return &Config{
		Amqp:   initAmqpConfig(),
		SMTP:   initSMTPConfig(),
		Logger: initLoggerConfig(),
	}
}

func parseConfigPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(wd)
}
