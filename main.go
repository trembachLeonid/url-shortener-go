package main

import (
	"net/http"
	"time"
	"url-shortener/internal/base62"
	"url-shortener/internal/helpers"
	"url-shortener/internal/models"
	"url-shortener/internal/persistence"

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

		var timeNow time.Time
		if request.ExpireTime != nil {
			timeNow = time.Now()
			timeNow = timeNow.Add(time.Duration(*request.ExpireTime) * time.Second)
		}

		urlId, err := persistence.InsertURL(normalizedURL, &timeNow)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		shortUrl := base62.ToBase62(urlId)

		c.JSON(http.StatusOK, gin.H{
			"short_url": shortUrl,
		})
	})

	err := router.Run()

	if err != nil {
		panic(err)
	}
}
