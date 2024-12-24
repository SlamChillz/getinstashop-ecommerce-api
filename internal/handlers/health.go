package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Service is running",
		"data":    gin.H{},
	})
}
