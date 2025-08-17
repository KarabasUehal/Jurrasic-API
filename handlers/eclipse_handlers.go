package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang-gin/models"
	"golang-gin/storage"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

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

	// Сохраняем в Redis с TTL 5 минут
	if redisClient != nil {
		err = redisClient.Set(ctx, cacheKey, dinosaurusJSON, 5*time.Minute).Err()
		if err != nil {
			log.Printf("Failed to cache items: %v", err)
			// Продолжаем, так как данные из БД уже получены
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
			// Продолжаем, так как данные из БД уже получены
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
		newDino,
	)

	ctx := context.Background()
	if redisClient != nil {
		err := redisClient.Del(ctx, "/dino").Err()
		if err != nil {
			log.Printf("Failed to invalidate cache for /dino: %v", err)
			// Не возвращаем ошибку клиенту, так как операция с БД успешна
		} else {
			log.Printf("Created dinosaur ID %d, invalidated cache for /dino", newDino.ID)
		}
	} else {
		log.Printf("Redis client is nil, skipping cache invalidation for /dino")
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
			updatedDino,
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
		err = redisClient.Del(ctx, "/dino").Err()
		if err != nil {
			log.Printf("Failed to invalidate cache for /dino: %v", err)
		} else {
			log.Printf("Updated dinosaur ID %d, invalidated cache for %s and /dino", id, c.Request.URL.String())
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
		err = redisClient.Del(ctx, "/dino").Err()
		if err != nil {
			log.Printf("Failed to invalidate cache for /dino: %v", err)
		} else {
			log.Printf("Deleted dinosaur ID %d, invalidated cache for %s and /dino", id, c.Request.URL.String())
		}
	} else {
		log.Printf("Redis client is nil, skipping cache invalidation for dinosaur ID %d", id)
	}

	c.Status(http.StatusNoContent)
}
