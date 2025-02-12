package config

import (
	"log"

	"github.com/spf13/viper"
)

type PostgresConfig struct {
	Host            string `mapstructure:"DB_HOST"`
	DbName          string `mapstructure:"DB_NAME"`
	Username        string `mapstructure:"DB_USER"`
	Password        string `mapstructure:"DB_PASSWORD"`
	Sslmode         string `mapstructure:"DB_SSL_MODE"`
	Port            int    `mapstructure:"DB_PORT"`
	MaxIdleConn     int    `mapstructure:"DB_MAX_IDLE_CONN"`
	MaxOpenConn     int    `mapstructure:"DB_MAX_OPEN_CONN"`
	MaxConnLifetime int    `mapstructure:"DB_CONN_MAX_LIFETIME"`
}

func initPostgresConfig() *PostgresConfig {
	pgConfig := &PostgresConfig{}

	if err := viper.Unmarshal(&pgConfig); err != nil {
		log.Fatalf("error mapping database config: %v", err)
	}

	return pgConfig
}
