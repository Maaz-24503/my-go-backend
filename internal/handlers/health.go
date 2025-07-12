package handlers

import (
	"github.com/gin-gonic/gin"
	"my-go-backend/pkg/models"
	"net/http"
)

func HealthCheck(c *gin.Context) {
	response := models.APIResponse{
		Success: true,
		Message: "Server is running",
		Data:    gin.H{"status": "healthy"},
	}
	c.JSON(http.StatusOK, response)
}
