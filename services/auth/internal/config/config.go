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
	Env        string     `yaml:"env" env-required:"true"`
	GRPC       ConfigGRPC `yaml:"grpc"`
	JWT        ConfigJWT
	Repository ConfigRepository
}

type ConfigRepository struct {
	Host     string `env:"AUTH_REPOSITORY_HOST" env-required:"true"`
	Port     int    `env:"AUTH_REPOSITORY_PORT" env-required:"true"`
	User     string `env:"AUTH_REPOSITORY_USER" env-required:"true"`
	Password string `env:"AUTH_REPOSITORY_PASSWORD" env-required:"true"`
	Name     string `env:"AUTH_REPOSITORY_NAME" env-required:"true"`
	SSLMode  string `env:"AUTH_REPOSITORY_SSLMODE" env-required:"true"`
}

type ConfigGRPC struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

type ConfigJWT struct {
	TTL    time.Duration `env:"JWT_TTL" env-required:"true"`
	Secret string        `env:"JWT_SECRET" env-required:"true"`
}

func MustLoad() *Config {
	_ = godotenv.Load()

	const configPathEnvVar = "AUTH_CONFIG_PATH"
	configPath, ok := os.LookupEnv(configPathEnvVar)
	if !ok {
		log.Fatalf("environment variable %s is not set", configPathEnvVar)
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
