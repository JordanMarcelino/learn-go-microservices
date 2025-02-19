package config

import (
	"log"

	"github.com/spf13/viper"
)

type FeignConfig struct {
	OrderURL string `mapstructure:"ORDER_URL"`
}

func initFeignConfig() *FeignConfig {
	feignCfg := &FeignConfig{}

	if err := viper.Unmarshal(&feignCfg); err != nil {
		log.Fatalf("error mapping feign config: %v", err)
	}

	return feignCfg
}
