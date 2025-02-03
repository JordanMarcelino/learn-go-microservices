package config

import (
	"log"

	"github.com/spf13/viper"
)

type SMTPConfig struct {
	Host  string `mapstructure:"SMTP_HOST"`
	Email string `mapstructure:"SMTP_EMAIL"`
	Port  int    `mapstructure:"SMTP_PORT"`
}

func initSMTPConfig() *SMTPConfig {
	smtpConfig := &SMTPConfig{}

	if err := viper.Unmarshal(&smtpConfig); err != nil {
		log.Fatalf("error mapping jwt config: %v", err)
	}

	return smtpConfig
}
