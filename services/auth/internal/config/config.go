package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env  string     `yaml:"env" env-required:"true"`
	GRPC ConfigGRPC `yaml:"grpc"`
	JWT  ConfigJWT
	Storage ConfigStorage
}

type ConfigStorage struct {
	Host string `env:"AUTH_STORAGE_HOST" env-required:"true"`
	Port int `env:"AUTH_STORAGE_PORT" env-required:"true"`
	User string `env:"AUTH_STORAGE_USER" env-required:"true"`
	Password string `env:"AUTH_STORAGE_PASSWORD" env-required:"true"`
	Name string `env:"AUTH_STORAGE_NAME" env-required:"true"`
	SSLMode string `env:"AUTH_STORAGE_SSLMODE" env-required:"true"`
}

type ConfigGRPC struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

type ConfigJWT struct {
	TTL    time.Duration `env:"AUTH_JWT_TTL" env-required:"true"`
	Secret string        `env:"AUTH_JWT_SECRET" env-required:"true"`
}

func MustLoad() *Config {
	_ = godotenv.Load()

	configPath, ok := os.LookupEnv("AUTH_CONFIG_PATH")
	if !ok {
		log.Fatal("environment variable AUTH_CONFIG_PATH is not set")
	}

	absPath, err := filepath.Abs(configPath)
	if err != nil {
		log.Fatalf("cannot resolve config path: %v", err)
	}

	cfg := new(Config)

	if err := cleanenv.ReadConfig(absPath, cfg); err != nil {
		log.Fatalf("config read error: %s", err)
	}

	return cfg
}
