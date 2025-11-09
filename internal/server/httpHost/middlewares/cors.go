package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMW() gin.HandlerFunc {
	origins := []string{}

	// мок для разработки
	if false {
		origins = []string{
			"http://localhost:3000",
			"https://polonium.ws",
		}
	}

	config := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Referer", "Origin", "User-Agent", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}

	if len(origins) == 0 {
		config.AllowAllOrigins = true
	}

	return cors.New(config)
}
