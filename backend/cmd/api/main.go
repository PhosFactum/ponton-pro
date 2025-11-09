package main

import (
	"log"

	"github.com/PhosFactum/TechnoLotos/backend/internal/config"
	"github.com/PhosFactum/TechnoLotos/backend/internal/database"
	"github.com/PhosFactum/TechnoLotos/backend/internal/handlers"
	"github.com/PhosFactum/TechnoLotos/backend/internal/middleware"
	"github.com/PhosFactum/TechnoLotos/backend/internal/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()
	log.Printf("Launching in mode: %s", cfg.Environment)

	// Подключаемся к БД
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Error while connecting to DB:", err)
	}

	// Создаём таблицы
	if err := database.Migrate(); err != nil {
		log.Fatal("Error while creating tables:", err)
	}

	// Добавляем тестовые данные
	database.SeedTestData()

	// Инициализируем роутер Gin
	router := gin.Default()

	// CORS для фронтенда
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Логирование (на всякий случай)
	router.Use(middleware.Logger())

	// API-роуты
	api := router.Group("/api")
	{
		api.GET("/", handlers.HealthCheck)
		api.GET("/health", handlers.HealthCheck)
		api.GET("/products", handlers.GetProducts)    // Каталог товаров
		api.POST("/requests", handlers.CreateRequest) // Создание заявки

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
