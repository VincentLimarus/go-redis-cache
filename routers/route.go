package routers

import (
	"VincentLimarus/go-redis/controllers/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    router.NoRoute(func(c *gin.Context) {
        c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Page not found"})
    })

    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, "Router is working")
    })

    route := router.Group("/api/v1")
    services.BaseStudentServices(route)

    return router
}
