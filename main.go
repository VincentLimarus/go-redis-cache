package main

import (
	"VincentLimarus/go-redis/configs"
	"VincentLimarus/go-redis/models/database"
	"VincentLimarus/go-redis/routers"
)

func init() {
	configs.LoadEnvironment()
    configs.ConnectToDB()
    configs.ConnectToRedis()
    // configs.Faker()
}

func main() {
    db := configs.GetDB()
    db.AutoMigrate(&database.Students{})

    r := routers.SetupRouter()
    r.Run(":3000")
}
