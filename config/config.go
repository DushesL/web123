package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSslMode  string
	Port       string
}

func Load() (*Config, error) {
	// загрузка .env
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error while loading .env")
	}

	// конфиг из переменных окружения
	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSslMode:  os.Getenv("DB_SSLMODE"),
		Port:       os.Getenv("PORT"),
	}

	if cfg.DBHost == "" {
		return nil, fmt.Errorf("missing DB_HOST in config")
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("missing DB_NAME in config")
	}

	return cfg, nil
}
