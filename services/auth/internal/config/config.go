package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env  string     `yaml:"env" env-required:"true"`
	GRPC GRPCConfig `yaml:"grpc"`
	JWT  JWTConfig
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

type JWTConfig struct {
	TTL    time.Duration `env:"JWT_TTL" env-required:"true"`
	Secret string        `env:"JWT_SECRET" env-required:"true"`
}

func MustLoad() *Config {
	_ = godotenv.Load()

	configPath, ok := os.LookupEnv("CONFIG_PATH")
	if !ok {
		log.Fatal("environment variable CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file is not exists: %s", configPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("config(%s) read error: %s", configPath, err)
	}

	return cfg
}
