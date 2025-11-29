package handlers

import (
	"net/http"

	"github.com/PhosFactum/TechnoLotos/backend/internal/database"
	"github.com/PhosFactum/TechnoLotos/backend/internal/models"
	"github.com/gin-gonic/gin"
)

// GetProducts - возвращает все продукты
func (h *Handler) GetProducts(c *gin.Context) {
	var products []models.Product
	result := database.DB.Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении продукта",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": products,
	})
}
