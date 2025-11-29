package handlers

import (
	"net/http"

	"github.com/PhosFactum/TechnoLotos/backend/internal/database"
	"github.com/PhosFactum/TechnoLotos/backend/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateRequest - создание новой заявки
func (h *Handler) CreateRequest(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Phone       string `json:"phone" binding:"required"`
		Email       string `json:"email,omitempty"`
		Description string `json:"description,omitempty"`
		ProductID   *uint  `json:"product_id,omitempty"`
	}

	// Валидация входных данных
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверные данные заявки: имя и телефон обязательны!",
		})
		return
	}

	// Простая валидация телефона (минимум 10 цифр)
	if len(input.Phone) < 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Телефон должен содержать минимум 10 цифр",
		})
		return
	}

	// Создаём заявку
	request := models.Request{
		Name:        input.Name,
		Phone:       input.Phone,
		Email:       input.Email,
		Description: input.Description,
		ProductID:   input.ProductID,
	}

	result := database.DB.Create(&request)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при создании заявки",
		})
		return
	}

	// Подгрузка данных о товаре для красивого вывода в ТГ
	if input.ProductID != nil {
		database.DB.Preload("Product").First(&request, request.ID)
	}

	// TODO: тут позже будет отправка уведомления
	go h.bot.SendNewRequestNotification(h.cfg.AdminChatID, request)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Заявка успешно создана!",
		"data":    request,
	})

}
