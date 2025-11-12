package main

import (
	"log"

	"github.com/PhosFactum/TechnoLotos/backend/internal/bot"
	"github.com/PhosFactum/TechnoLotos/backend/internal/config"
	"github.com/PhosFactum/TechnoLotos/backend/internal/database"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()
	log.Printf("Launching bot in mode: %s", cfg.Environment)

	// Проверяем наличие токена бота
	if cfg.BotToken == "" {
		log.Fatal("BOT_TOKEN env variable is required!")
	}

	// Подключаемся к БД
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Error while connecting to DB:", err)
	}

	// Создаём и запускаем бота
	bot, err := bot.NewBot(cfg.BotToken)
	if err != nil {
		log.Fatal("Error while creating bot:", err)
	}

	log.Println("Bot is starting...")
	bot.Start()
}
