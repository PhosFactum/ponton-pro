package database

import (
	"log"

	"github.com/PhosFactum/TechnoLotos/backend/internal/config"
	"github.com/PhosFactum/TechnoLotos/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect - подключение к БД
func Connect(cfg *config.Config) error {
	var err error
	var dbName string

	if cfg.Environment == "production" {
		dbName = "/var/lib/TechnoLotos/production.db"
	} else {
		dbName = "lotos.db"
	}

	// SQLite - идеальная легковесная БД для SPA
	DB, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return err
	}

	log.Println("Database has been connected!")
	return nil
}

// Migrate - создание таблиц
func Migrate() error {
	// Создаём таблицы Продуктов и Заявок
	err := DB.AutoMigrate(
		&models.Product{},
		&models.Request{},
	)

	if err != nil {
		return err
	}

	log.Println("Tables have been created: products, requests")
	return nil
}
