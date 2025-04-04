package config

import (
	"ApiSmart/pkg/database"
	"os"
)

type Config struct {
	ServerPort string
	DBConfig   database.DBConfig
	JWTSecret  string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8000"),
		DBConfig: database.DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "manuel"),
			Password: getEnv("DB_PASSWORD", "manuel"),
			DBName:   getEnv("DB_NAME", "sensores_db"),
		},
		JWTSecret: getEnv("JWT_SECRET", "secret_key_cambiar_en_produccion"),
	}
}


func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
