package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Debug      bool   `mapstructure:"debug"`
	Port       string `mapstructure:"port"`
	Passphrase string `mapstructure:"passphrase"`
	APIURL     string `mapstructure:"api_url"`
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
