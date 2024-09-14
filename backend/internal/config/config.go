package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Supabase SupabaseConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type SupabaseConfig struct {
	URL string
	Key string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
