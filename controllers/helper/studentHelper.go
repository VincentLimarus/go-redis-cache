package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
		var students []database.Students
		if jsonErr := json.Unmarshal([]byte(cachedStudents), &students); jsonErr == nil {
			return 200, students
		}
	} else if err != redis.Nil {
		return 500, gin.H{"message": "Redis error", "error": err.Error()}
	}

	var students []database.Students
	if result := db.Find(&students); result.Error != nil {
		return 500, result.Error
	}

	studentData, _ := json.Marshal(students)

	pipe := redisClient.Pipeline()
	pipe.Set(ctx, cacheKey, studentData, 30*time.Minute)
	pipe.Exec(ctx)

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
		fmt.Println("Redis error:", err) 
	} else if len(studentData) > 0 {
		if studentData["id"] != "" && studentData["name"] != "" && studentData["email"] != "" {
			student := database.Students{
				ID:    uuid.MustParse(studentData["id"]),
				Name:  studentData["name"],
				Email: studentData["email"],
			}
			return 200, student
		}
	}

	var student database.Students
	if result := db.Where("id = ?", studentUUID).First(&student); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 404, gin.H{"message": "Student not found"}
		}
		return 500, gin.H{"message": "Database error", "error": result.Error.Error()}
	}

	if redisClient != nil {
		pipe := redisClient.Pipeline()
		pipe.HMSet(ctx, cacheKey, map[string]interface{}{
			"id":    student.ID.String(),
			"name":  student.Name,
			"email": student.Email,
		})
		pipe.Expire(ctx, cacheKey, 30*time.Minute)
		_, _ = pipe.Exec(ctx)
	}

	return 200, student
}
