package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Node         string   `mapstructure:"NODE"`
	PrivateKeys  []string `mapstructure:"PRIVATE_KEYS"`
	DatabasePath string   `mapstructure:"DATABASE_PATH"`
}

// Load reads configuration from environment variables and config files
func Load() (*Config, error) {
	// Set default values
	viper.SetDefault("DATABASE_PATH", "./main.db")

	// Try to read from main.env file
	viper.SetConfigFile("main.env")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		// If main.env doesn't exist, that's okay
		if _, ok := err.(*os.PathError); !ok {
			return nil, fmt.Errorf("error reading main.env: %w", err)
		}
	}

	// Try to read from .env file (has higher priority)
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		if err := viper.MergeInConfig(); err != nil {
			return nil, fmt.Errorf("error reading .env: %w", err)
		}
	}

	// Override config parameters from environment variables
	viper.AutomaticEnv()

	// Bind environment variables
	if err := viper.BindEnv("NODE"); err != nil {
		return nil, fmt.Errorf("error binding NODE env: %w", err)
	}
	if err := viper.BindEnv("PRIVATE_KEYS"); err != nil {
		return nil, fmt.Errorf("error binding PRIVATE_KEYS env: %w", err)
	}
	if err := viper.BindEnv("DATABASE_PATH"); err != nil {
		return nil, fmt.Errorf("error binding DATABASE_PATH env: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate required fields
	if err := config.validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// validate checks that all required configuration values are present
func (c *Config) validate() error {
	if c.Node == "" {
		return fmt.Errorf("NODE is required")
	}

	if len(c.PrivateKeys) == 0 {
		return fmt.Errorf("PRIVATE_KEYS is required")
	}

	if c.DatabasePath == "" {
		return fmt.Errorf("DATABASE_PATH is required")
	}

	return nil
}
