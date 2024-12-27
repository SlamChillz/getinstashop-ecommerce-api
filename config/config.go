package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	DatabaseURL       string `mapstructure:"DATABASE_URL"`
	MigrationURL      string `mapstructure:"MIGRATION_URL"`
	JwtSecret         string `mapstructure:"JWT_SECRET"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(fmt.Sprintf("%s%s", path, ".env"))
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
