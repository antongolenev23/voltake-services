package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	Storage    ConfigStorage
	HTTPServer ConfigHTTPServer `yaml:"http_server"`
	JWT        ConfigJWT
}

type ConfigStorage struct {
	Host     string `env:"BOOKING_STORAGE_HOST" env-required:"true"`
	Port     int    `env:"BOOKING_STORAGE_PORT" env-required:"true"`
	User     string `env:"BOOKING_STORAGE_USER" env-required:"true"`
	Password string `env:"BOOKING_STORAGE_PASSWORD" env-required:"true"`
	Name     string `env:"BOOKING_STORAGE_NAME" env-required:"true"`
	SSLMode  string `env:"BOOKING_STORAGE_SSLMODE" env-required:"true"`
}

type ConfigHTTPServer struct {
	Address              string        `yaml:"address" env-required:"true"`
	RequestReadTimeout   time.Duration `yaml:"request_read_timeout" env-required:"true"`
	ResponseWriteTimeout time.Duration `yaml:"response_write_timeout" env-required:"true"`
	IdleTimeout          time.Duration `yaml:"idle_timeout" env-required:"true"`
}

type ConfigJWT struct {
	TTL    time.Duration `env:"JWT_TTL" env-required:"true"`
	Secret string        `env:"JWT_SECRET" env-required:"true"`
}

func MustLoad() *Config {
	_ = godotenv.Load()

	const configPathEnvVar = "BOOKING_CONFIG_PATH"
	configPath, ok := os.LookupEnv(configPathEnvVar)
	if !ok {
		log.Fatalf("environment variable %s is not set", configPathEnvVar)
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
