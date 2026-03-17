package main

import (
	"net/http"
	"time"
	"url-shortener/internal/helpers"
	"url-shortener/internal/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()

	router.POST("/url", func(c *gin.Context) {
		var request models.ShortenURLRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		normalizedURL, err := helpers.NormalizeURL(request.OriginalURL)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		timeNow := time.Now()
		timeNow.Add(*time.Second)

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	err := router.Run()

	if err != nil {
		panic(err)
	}
}
