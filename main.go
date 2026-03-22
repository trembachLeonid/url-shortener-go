package main

import (
	"crypto/sha256"
	"log"
	"net/http"
	"os"
	"time"
	"url-shortener/internal/base62"
	"url-shortener/internal/helpers"
	"url-shortener/internal/models"
	"url-shortener/internal/persistence"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var base62Secret = []byte(os.Getenv("BASE62_SECRET"))
var hostURL = os.Getenv("HOST_URL")

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

		for i := 0; i < 10; i++ {
			urlId, err := persistence.InsertURL(normalizedURL, timeNow)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			idInBytes := helpers.Int32ToBytes(urlId)
			toHash := append(idInBytes, base62Secret...)
			hashed := sha256.Sum256(toHash)
			shortURL := base62.ToBase62(hashed[:6])

			if exists, err := persistence.ShortURLExists(shortURL); !exists && err == nil {
				err = persistence.SetShortURL(urlId, shortURL)
				result := hostURL + shortURL
				c.JSON(http.StatusOK, gin.H{
					"short_url": result,
				})
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			}
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
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
