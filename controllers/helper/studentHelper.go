package helper

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"VincentLimarus/go-redis/configs"
	"VincentLimarus/go-redis/models/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetAllStudent() (int, interface{}) {
	ctx := context.Background()
	redisClient := configs.GetRedis()
	db := configs.GetDB()

	if redisClient == nil || db == nil {
		return 500, gin.H{"message": "Database or Redis client not initialized"}
	}

	cacheKey := "students:all"
	cachedStudents, err := redisClient.Get(ctx, cacheKey).Result()

	if err == nil {
		log.Println("Cache Hit!")
		var students []database.Students
		if jsonErr := json.Unmarshal([]byte(cachedStudents), &students); jsonErr == nil {
			return 200, students
		} else {
			log.Printf("Cache unmarshal error: %v", jsonErr)
		}
	} else if err != redis.Nil {
		log.Printf("Redis error: %v", err)
	}

	log.Println("Cache Miss!")
	var students []database.Students
	if result := db.Find(&students); result.Error != nil {
		return 500, gin.H{"message": "Database error", "error": result.Error.Error()}
	}

	studentData, err := json.Marshal(students)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return 500, gin.H{"message": "Internal server error"}
	}

	if err := redisClient.Set(ctx, cacheKey, studentData, 30*time.Minute).Err(); err != nil {
		log.Printf("Redis set error: %v", err)
	}

	return 200, students
}

func GetStudent(studentID string) (int, interface{}) {
	ctx := context.Background()
	redisClient := configs.GetRedis()
	db := configs.GetDB()

	if redisClient == nil || db == nil {
		return 500, gin.H{"message": "Database or Redis client not initialized"}
	}

	studentUUID, err := uuid.Parse(studentID)
	if err != nil {
		return 400, gin.H{"message": "Invalid student ID"}
	}

	cacheKey := "student:" + studentID
	studentData, err := redisClient.HGetAll(ctx, cacheKey).Result()

	if err != nil && err != redis.Nil {
		log.Printf("Redis error: %v", err)
	} else if len(studentData) > 0 {
		idStr := studentData["id"]
		name := studentData["name"]
		email := studentData["email"]

		if idStr != "" && name != "" && email != "" {
			id, parseErr := uuid.Parse(idStr)
			if parseErr == nil {
				student := database.Students{
					ID:    id,
					Name:  name,
					Email: email,
				}
				return 200, student
			}
			log.Printf("Failed to parse cached UUID: %v", parseErr)
		}
	}

	var student database.Students
	if result := db.Where("id = ?", studentUUID).First(&student); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 404, gin.H{"message": "Student not found"}
		}
		return 500, gin.H{"message": "Database error", "error": result.Error.Error()}
	}

	pipe := redisClient.Pipeline()
	pipe.HSet(ctx, cacheKey, map[string]interface{}{
		"id":    student.ID.String(),
		"name":  student.Name,
		"email": student.Email,
	})
	pipe.Expire(ctx, cacheKey, 30*time.Minute)
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("Redis pipeline error: %v", err)
	}


	return 200, student
}