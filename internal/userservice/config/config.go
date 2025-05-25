package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
}

type GRPCConfig struct {
}

func LoadConfig() *Config {
	var config Config
	if err := cleanenv.ParseYAML(); err != nil {
		panic()
	}

	return &config
}
