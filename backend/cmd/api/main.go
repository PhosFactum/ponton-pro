package main

import (
	"log"

	"github.com/PhosFactum/TechnoLotos/backend/internal/bot"
	"github.com/PhosFactum/TechnoLotos/backend/internal/config"
	"github.com/PhosFactum/TechnoLotos/backend/internal/database"
	"github.com/PhosFactum/TechnoLotos/backend/internal/handlers"
	"github.com/PhosFactum/TechnoLotos/backend/internal/middleware"
	"github.com/PhosFactum/TechnoLotos/backend/internal/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Загружаем конфигурацию
	cfg := config.Load()
	log.Printf("Launching in mode: %s", cfg.Environment)

	// 2. Подключаемся к БД
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Error while connecting to DB:", err)
	}

	// 3. Миграции и сиды
	if err := database.Migrate(); err != nil {
		log.Fatal("Error while creating tables:", err)
	}
	database.SeedTestData()

	// 4. Инициализируем бота
	if cfg.BotToken == "" {
		log.Fatal("BOT_TOKEN is required on .env!")
	}

	tgBot, err := bot.NewBot(cfg.BotToken)
	if err != nil {
		log.Fatal("Error while creating telegram bot:", err)
	}

	log.Println("Starting Telegram Bot in background...")
	go tgBot.Start() // Запускаем бота в отдельной горутине

	// 5. Инициализируем хэндлеры (внедряем бота и конфиг)
	h := handlers.NewHandler(tgBot, cfg)

	// 6. Инициализируем роутер Gin
	router := gin.Default()

	// CORS для фронтенда
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 7. Логирование и прочие подключения
	router.Use(middleware.Logger())

	// 8. API-роуты
	api := router.Group("/api")
	{
		api.GET("/", h.HealthCheck)
		api.GET("/health", h.HealthCheck)      // Проверка состояния сервера
		api.GET("/products", h.GetProducts)    // Каталог товаров
		api.POST("/requests", h.CreateRequest) // Создание заявки

		// ВРЕМЕННЫЙ ЭНДПОЙНТ ДЛЯ ПРОСМОТРА ВСЕХ ЗАЯВОК
		api.GET("/debug/requests", func(c *gin.Context) {
			var requests []models.Request
			database.DB.Find(&requests)

			c.JSON(200, gin.H{
				"count":    len(requests),
				"requests": requests,
			})
		})
	}

	// Запускаем сервер на порту 8080
	log.Printf("Сервер запускается на http://localhost:%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
