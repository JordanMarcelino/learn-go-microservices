package config

import (
	"log"

	"github.com/spf13/viper"
)

type RedisConfig struct {
	Addrs           []string `mapstructure:"REDIS_ADDRS"`
	Password        string   `mapstructure:"REDIS_PASSWORD"`
	DialTimeout     int      `mapstructure:"REDIS_DIAL_TIMEOUT"`
	ReadTimeout     int      `mapstructure:"REDIS_READ_TIMEOUT"`
	WriteTimeout    int      `mapstructure:"REDIS_WRITE_TIMEOUT"`
	MinIdleConn     int      `mapstructure:"REDIS_MIN_IDLE_CONN"`
	MaxIdleConn     int      `mapstructure:"REDIS_MAX_IDLE_CONN"`
	MaxActiveConn   int      `mapstructure:"REDIS_MAX_ACTIVE_CONN"`
	MaxConnLifetime int      `mapstructure:"REDIS_MAX_CONN_LIFETIME"`
}

func initRedisConfig() *RedisConfig {
	rdConfig := &RedisConfig{}

	if err := viper.Unmarshal(&rdConfig); err != nil {
		log.Fatalf("error mapping redis config: %v", err)
	}

	return rdConfig
}
