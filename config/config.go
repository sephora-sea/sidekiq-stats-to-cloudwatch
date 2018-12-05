package config

import (
	"sync"

	"github.com/caarlos0/env"
)

type Config struct {
	APIKey    string `env:"API_KEY"`
	AWSRegion string `env:"AWS_REGION" envDefault:"ap-southeast-1"`
	AppName   string `env:"APP_NAME" envDefault:"nameless app"`
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
		env.Parse(instance)
	})
	return instance
}
