package config

import "os"

type Config struct {
	Port        string
	DatabaseUrl string
	Environment string
}

// Load - загрузка конфигурации
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseUrl: getEnv("DATABASE_URL", "lotos.db"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

// getEnv - получение переменной окружения или значения по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
