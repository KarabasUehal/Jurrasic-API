package main

import (
	"Dinosaurus/handlers"
	"Dinosaurus/models"
	"Dinosaurus/storage"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var jwtKey = []byte("amfkdhfneigjtnfkgmdlsvmutskgsjrg")

func main() {

	if err := storage.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(cacheMiddleware)

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

func cacheMiddleware(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.Next()
		return
	}

	cacheKey := c.Request.URL.String()
	ctx := context.Background()

	// Проверяем, инициализирован ли redisClient
	if redisClient == nil {
		log.Printf("Redis client is nil, skipping cache for key: %s", cacheKey)
		c.Next()
		return
	}

	// Проверяем кеш в Redis
	cached, err := redisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		log.Printf("Cache miss for key: %s", cacheKey)
		c.Next()
		return
	}
	if err != nil {
		log.Printf("Redis error for key %s: %v", cacheKey, err)
		c.Next()
		return
	}

	// Десериализация в зависимости от маршрута
	if cacheKey == "/dino" {
		var dinosaurus []models.Dinosaurus
		if err := json.Unmarshal([]byte(cached), &dinosaurus); err != nil {
			log.Printf("Failed to unmarshal cached items list for key %s: %v", cacheKey, err)
			c.Next()
			return
		}
		log.Printf("Cache hit for dinosaurus list")
		c.JSON(http.StatusOK, dinosaurus)
	} else {
		var dinosaur models.Dinosaurus
		if err := json.Unmarshal([]byte(cached), &dinosaur); err != nil {
			log.Printf("Failed to unmarshal cached dinosaur for key %s: %v", cacheKey, err)
			c.Next()
			return
		}
		log.Printf("Cache hit for key: %s", cacheKey)
		c.JSON(http.StatusOK, dinosaur)
	}
	c.Abort()
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
