package config

import (
	"log"
	"log/slog"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server Server `yaml:"server"`
	Redis  Redis  `yaml:"redis"`
}

var instance *Config
var once sync.Once

func Init(pathConfig string) *Config {
	once.Do(func() {
		slog.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig(pathConfig, instance); err != nil {
			log.Fatalf("fail read config: %v", err)
		}
	})
	return instance
}
