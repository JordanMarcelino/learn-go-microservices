package config

import (
	"log"

	"github.com/spf13/viper"
)

type SMTPConfig struct {
	SMTPHost  string `mapstructure:"SMTP_HOST"`
	SMTPEMail string `mapstructure:"SMTP_EMAIL"`
	SMTPPort  int    `mapstructure:"SMTP_PORT"`
}

func initSMTPConfig() *SMTPConfig {
	smtpConfig := &SMTPConfig{}

	if err := viper.Unmarshal(&smtpConfig); err != nil {
		log.Fatalf("error mapping smtp config: %v", err)
	}

	return smtpConfig
}
