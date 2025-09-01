package main

import (
	"Dinosaurus/handlers"
	"Dinosaurus/storage"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("amfkdhfneigjtnfkgmdlsvmutskgsjrg")

func main() {

	if err := storage.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	redisClient := storage.NewRedis()
	if redisClient == nil {
		log.Fatalf("Error to start redis")
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	// Защищенные роуты с JWT
	api := router.Group("/api")
	api.Use(AuthMiddleware())
	api.GET("/dinosaurus", handlers.GetAllDinosaurus)
	api.GET("/dinosaurus/:id", handlers.GetDinosaurByID)
	api.POST("/dinosaurus", handlers.AddDinosaur)
	api.PUT("/dinosaurus/:id", handlers.UpdateDinosaurByID)
	api.DELETE("/dinosaurus/:id", handlers.DeleteDinosaurByID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
