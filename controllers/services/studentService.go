package services

import (
	"VincentLimarus/go-redis/controllers/helper"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllStudent(c *gin.Context) {
	code, result := helper.GetAllStudent()
	if code != 200 {
		c.JSON(code, gin.H{"message": "Failed to retrieve students", "error": result})
		return
	}
	c.JSON(200, gin.H{"message": "Success", "data": result})
}

func GetStudentByID(c *gin.Context) {
	studentID := c.Param("id")

	if _, err := uuid.Parse(studentID); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid UUID",
		})
		return
	}

	code, result := helper.GetStudent(studentID)
	if code != 200 {
		c.JSON(code, gin.H{"message": "Failed to retrieve student", "error": result})
		return
	}
	c.JSON(200, gin.H{"message": "Success", "data": result})
}

func BaseStudentServices(route *gin.RouterGroup) {
	route.GET("/students", GetAllStudent)
	route.GET("/student/:id", GetStudentByID)
}
