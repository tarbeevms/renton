package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config struct represents the configuration structure
type Config struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
}

// MySQLConfig struct represents the MySQL specific configuration
type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// GetConfig reads and parses the configuration file
func GetConfig() (*Config, error) {
	viper.SetConfigName("config") // name of the config file (without extension)
	viper.SetConfigType("yml")    // config file type
	viper.AddConfigPath("../../") // look for the config file in the current directory
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}
