package database

import (
	"log"

	"github.com/PhosFactum/TechnoLotos/backend/internal/models"
)

// SeedTestData - заполнение тестовыми продуктами
func SeedTestData() {
	var productCount int64
	DB.Model(&models.Product{}).Count(&productCount)

	// Добавляем тестовые данные, только если таблица пустая
	if productCount == 0 {
		products := []models.Product{
			{
				Title:       "Модульные понтоны из ПНД",
				Description: "Прочные и долговечные понтоны для пирсов, причалов и плавучих платформ. Устойчивы к коррозии и УФ-излучению.",
			},
			{
				Title:       "Понтоны для катеров и яхт",
				Description: "Специализированные понтонные системы для безопасной швартовки маломерных судов. Повышенная грузоподъемность.",
			},
			{
				Title:       "Плавучие платформы",
				Description: "Многофункциональные платформы для мероприятий, кафе и зон отдыха. Быстрый монтаж и демонтаж.",
			},
			{
				Title:       "Лодки из ПНД",
				Description: "Прочные и надежные лодки из пищевого полиэтилена. Идеальны для рыбалки и отдыха на воде.",
			},
		}

		DB.Create(&products)
		log.Println("Added test products:", len(products))
	}
}

// Заявки не добавляются, ибо они создаются через само API
