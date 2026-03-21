package main

import (
	"log"
	"net/http"
	"time"
	"url-shortener/internal/base62"
	"url-shortener/internal/helpers"
	"url-shortener/internal/models"
	"url-shortener/internal/persistence"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()

	err := godotenv.Load()

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

		var timeNow *time.Time
		if request.ExpireTime != nil {
			log.Printf("CREATE URL >> Expire time: %v", request.ExpireTime)
			expireTime := time.Now().UTC().Add(time.Duration(*request.ExpireTime) * time.Second)
			timeNow = &expireTime
		}

		log.Printf("CREATE URL >> URL: %s, Expire time: %v", normalizedURL, request.ExpireTime)
		urlId, err := persistence.InsertURL(normalizedURL, timeNow)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		shortUrl := base62.ToBase62(urlId)

		err = persistence.SetShortURL(urlId, shortUrl)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"short_url": shortUrl,
		})
	})

	router.GET("/:shortUrl", func(c *gin.Context) {
		shortUrl := c.Param("shortUrl")

		originalURL, expireTime, err := persistence.GetURL(shortUrl)

		log.Printf("Original URL: %s, Expire Time: %v", originalURL, expireTime)

		if err != nil || (expireTime != nil && expireTime.Before(time.Now().UTC())) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "URL not found",
			})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, originalURL)
	})

	err = router.Run()

	if err != nil {
		panic(err)
	}
}
