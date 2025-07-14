package config

import "github.com/spf13/viper"

type Config struct {
	Port   int      `mapstructure:"port"`
	Target []string `mapstructure:"target"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	// setup configuration settings
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("settings")

	// default values
	viper.SetDefault("port", 8080)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
