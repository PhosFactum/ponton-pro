package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck - проверка состояния API
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "API ТехноЛотос работает!",
	})
}
