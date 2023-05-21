package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Debug bool   `mapstructure:"debug"`
	Port  string `mapstructure:"port"`

	Database Database `mapstructure:"db"`
	RabbitMQ RabbitMQ `mapstructure:"rabbitmq"`
}

type Database struct {
	Driver string `mapstructure:"driver"`
	URL    string `mapstructure:"url"`
}

type RabbitMQ struct {
	URL string `mapstructure:"url"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
