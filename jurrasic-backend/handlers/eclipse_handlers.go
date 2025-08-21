package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"Dinosaurus/models"
	"Dinosaurus/storage"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

var redisClient *redis.Client
var jwtKey = []byte("amfkdhfneigjtnfkgmdlsvmutskgsjrg")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GetAllDinosaurus(c *gin.Context) {
	cacheKey := c.Request.URL.String()
	ctx := context.Background()

	if redisClient != nil {
		cached, err := redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var Dinosaurus []models.Dinosaurus
			if json.Unmarshal([]byte(cached), &Dinosaurus) == nil {
				log.Printf("Cache hit for dinosaurus list")
				c.JSON(http.StatusOK, Dinosaurus)
				return
			}
			log.Printf("Failed to unmarshal cached dinosaurus list: %v", err)
		} else if err != redis.Nil {
			log.Printf("Redis error for dinosaurus list: %v", err)
		}
	} else {
		log.Printf("Redis client is nil, skipping cache for dinosaurus list")
	}

	dinosaurus := storage.GetAllDinosaurus()
	if dinosaurus == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query dinosaurus"})
		return
	}

	dinosaurusJSON, err := json.Marshal(dinosaurus)
	if err != nil {
		log.Printf("Failed to serialize dinosaurus: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize dinosaurus"})
		return
	}

	if redisClient != nil {
		err = redisClient.Set(ctx, cacheKey, dinosaurusJSON, 5*time.Minute).Err()
		if err != nil {
			log.Printf("Failed to cache items: %v", err)
		}
	} else {
		log.Printf("Redis client is nil, skipping cache for dinosaurus list")
	}

	c.JSON(http.StatusOK, dinosaurus)
}

func GetDinosaurByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("Invalid ID format: %d, error: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid dinosaur ID"})
		return
	}

	cacheKey := c.Request.URL.String()
	ctx := context.Background()

	if redisClient != nil {
		cached, err := redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var dino models.Dinosaurus
			if json.Unmarshal([]byte(cached), &dino) == nil {
				log.Printf("Cache hit for dinosaur ID: %d", id)
				c.JSON(http.StatusOK, dino)
				return
			}
			log.Printf("Failed to unmarshal cached dinosaur ID %d: %v", id, err)
		} else if err != redis.Nil {
			log.Printf("Redis error for dinosaur ID %d: %v", id, err)
		}
	} else {
		log.Printf("Redis client is nil, skipping cache for dinosaur ID %d", id)
	}

	dinosaur := storage.GetDinosaurByID(id)
	if dinosaur == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dinosaur not found"})
		return
	}

	dinosaurJSON, err := json.Marshal(dinosaur)
	if err != nil {
		log.Printf("Failed to serialize dinosaur ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize dinosaur"})
		return
	}

	if redisClient != nil {
		err = redisClient.Set(ctx, cacheKey, dinosaurJSON, 5*time.Minute).Err()
		if err != nil {
			log.Printf("Failed to cache dinosaur ID %d: %v", id, err)
		}
	} else {
		log.Printf("Redis client is nil, skipping cache for dinosaur ID %d", id)
	}

	c.JSON(http.StatusOK, dinosaur)
}

func AddDinosaur(c *gin.Context) {
	var newDino models.Dinosaurus

	if err := c.ShouldBindJSON(&newDino); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Invalid of input data": err.Error()})
		return
	}

	dinosaur := storage.AddDinosaur(
		&newDino,
	)

	ctx := context.Background()
	if redisClient != nil {
		err := redisClient.Del(ctx, "/api/dinosaurus").Err()
		if err != nil {
			log.Printf("Failed to invalidate cache for /api/dinosaurus: %v", err)
		} else {
			log.Printf("Created dinosaur ID %d, invalidated cache for /api/dinosaurus", newDino.ID)
		}
	} else {
		log.Printf("Redis client is nil, skipping cache invalidation for /api/dinosaurus")
	}

	c.JSON(http.StatusCreated, dinosaur)
}

func UpdateDinosaurByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("Invalid ID format: %d, error: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid eclipse ID",
		})
		return
	}
	var updatedDino models.Dinosaurus

	if err := c.ShouldBindJSON(&updatedDino); err != nil {
		log.Printf("Failed to bind JSON for item ID %d: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Invalid of input data": err.Error()})
		return
	}

	dinosaur := storage.UpdateDinosaurByID(
		id,
		updatedDino,
	)

	if dinosaur == nil {
		NewDinosaur := storage.AddDinosaur(
			&updatedDino,
		)
		c.JSON(http.StatusCreated, NewDinosaur)
		return
	}

	ctx := context.Background()
	if redisClient != nil {
		err = redisClient.Del(ctx, c.Request.URL.String()).Err()
		if err != nil {
			log.Printf("Failed to invalidate cache for %s: %v", c.Request.URL.String(), err)
		}
		err = redisClient.Del(ctx, "/api/dinosaurus").Err()
		if err != nil {
			log.Printf("Failed to invalidate cache for /api/dinosaurus: %v", err)
		} else {
			log.Printf("Updated dinosaur ID %d, invalidated cache for %s and /api/dinosaurus", id, c.Request.URL.String())
		}
	} else {
		log.Printf("Redis client is nil, skipping cache invalidation for dinosaur ID %d", id)
	}

	c.JSON(http.StatusOK, dinosaur)
}

func DeleteDinosaurByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("Invalid ID format: %d, error: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid dinosaur ID",
		})
		return
	}

	if success := storage.DeleteDinosaurByID(id); !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dinosaur not found"})
		return
	}

	ctx := context.Background()
	if redisClient != nil {
		err = redisClient.Del(ctx, c.Request.URL.String()).Err()
		if err != nil {
			log.Printf("Failed to invalidate cache for %s: %v", c.Request.URL.String(), err)
		}
		err = redisClient.Del(ctx, "/api/dinosaurus").Err()
		if err != nil {
			log.Printf("Failed to invalidate cache for /api/dinosaurus: %v", err)
		} else {
			log.Printf("Deleted dinosaur ID %d, invalidated cache for %s and /api/dinosaurus", id, c.Request.URL.String())
		}
	} else {
		log.Printf("Redis client is nil, skipping cache invalidation for dinosaur ID %d", id)
	}

	c.Status(http.StatusNoContent)
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(hashedPassword)
	storage.DB.Create(&user)

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	c.JSON(http.StatusOK, gin.H{"message": "User registered", "token": tokenString})
}

func Login(c *gin.Context) {
	var user models.User
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	storage.DB.Where("username = ?", input.Username).First(&user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: input.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
