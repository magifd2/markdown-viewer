package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config stores all configuration for the application.
// The values are read by viper from a config file, environment variables, or command line flags.
type Config struct {
	Port       int    `mapstructure:"port"`
	Open       bool   `mapstructure:"open"`	TargetDir  string `mapstructure:"target_dir"`
}

// LoadConfig reads configuration from file and environment variables.
func LoadConfig() (config Config, err error) {
	// Set defaults
	viper.SetDefault("port", 8888)
	viper.SetDefault("open", false)
	viper.SetDefault("target_dir", ".")

	// Set up config file
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/mdv")
	viper.AddConfigPath(".")

	// Read config file
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			return
		}
	}

	// Set up environment variables
	viper.SetEnvPrefix("MDV")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Unmarshal config into struct
	err = viper.Unmarshal(&config)
	return
}