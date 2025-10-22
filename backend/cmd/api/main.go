package main

import (
	"log"
	"net/http"

	"github.com/PhosFactum/TechnoLotos/backend/internal/config"
	"github.com/PhosFactum/TechnoLotos/backend/internal/database"
	"github.com/PhosFactum/TechnoLotos/backend/internal/models"
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

	// Базовый маршрут для проверки работы
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API понтончиков работает!",
			"status":  "success",
		})
	})

	// Маршрут для проверки состояния работы API
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// ВРЕМЕННЫЙ маршрут для проверки продуктов из БД
	router.GET("/products", func(c *gin.Context) {
		var products []models.Product
		result := database.DB.Find(&products)

		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Ошибка БД"})
			return
		}

		c.JSON(200, gin.H{
			"count":    result.RowsAffected,
			"products": products,
		})
	})

	// Запускаем сервер на порту 8080
	log.Printf("Сервер запускается на http://localhost:%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
