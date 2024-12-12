package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DatabasePath  string
}

func Load() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	token := os.Getenv("TELEGRAM_APITOKEN")
	if token == "" {
		return nil, ErrMissingTelegramToken
	}
	db_path := os.Getenv("STORAGE_PATH")
	if db_path == "" {
		return nil, ErrMissingDBPath
	}

	return &Config{
		TelegramToken: token,
		DatabasePath:  db_path,
	}, nil
}

var ErrMissingTelegramToken = &ConfigError{"TELEGRAM_TOKEN не задан"}
var ErrMissingDBPath = &ConfigError{"TELEGRAM_TOKEN не задан"}

type ConfigError struct {
	msg string
}

func (e *ConfigError) Error() string {
	return e.msg
}
