package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var JWT_CONFIG *JwtConfig

func init() {
	configPath := parseConfigPath()
	viper.AddConfigPath(configPath)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	JWT_CONFIG = initJwtConfig()
}

func parseConfigPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(wd)
}
