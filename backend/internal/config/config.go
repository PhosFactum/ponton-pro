package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	IP	    string
	Port        string
	DatabaseUrl string
	Environment string
	BotToken    string
	AdminChatID int64 // ID чата, куда бот будет кидать заявки
}

// Load - загрузка конфигурации
func Load() *Config {
	envPath, err := findEnvFile()
	if err == nil {
		log.Printf("Loading .env from: %s", envPath)
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Error while loading .env file: %v", err)
		}
	} else {
		log.Println(".env file was not found!")
	}

	return &Config{
		IP: 	     getEnv("IP", "localhost"),
		Port:        getEnv("PORT", "8080"),
		DatabaseUrl: getEnv("DATABASE_URL", "lotos.db"),
		Environment: getEnv("ENVIRONMENT", "development"),
		BotToken:    getEnv("BOT_TOKEN", ""),
		AdminChatID: getEnvAsInt64("ADMIN_CHAT_ID", 0),
	}
}

// findEnvFile - рекурсивно ищем файл .env
func findEnvFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		path := filepath.Join(dir, ".env")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			return "", os.ErrNotExist
		}

		dir = parent
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

// getEnvAsInt64 - спецфункция для численных параметров
func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}
