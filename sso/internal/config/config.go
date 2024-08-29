package config

import (
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-required:"true"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-default:"24h"`
	Grpc        GrpcConfig    `yaml:"grpc"`
}

type GrpcConfig struct {
	Port    int           `yaml:"port" env-default:"8082"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		cfg = MustLoad()
	})
	return cfg
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadWithArg(configPath)
}

func MustLoadWithArg(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}
