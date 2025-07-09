package main

import "github.com/spf13/viper"

type config struct {
	Port      int    `mapstructure:"port"`
	MediaPath string `mapstructure:"media_path"`
}

func loadConfig() (*config, error) {
	var cfg config

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
