package main

import (
	"VincentLimarus/go-redis/configs"
	"VincentLimarus/go-redis/models/database"
	"VincentLimarus/go-redis/routers"
	"log"
)

func init() {
	configs.LoadEnvironment()
    if err := configs.ConnectToDB(); err != nil {
        log.Fatalf("Error: %v", err)
    }
    
    if err := configs.ConnectToRedis(); err != nil {
        log.Fatalf("Error: %v", err)
    }
    
    // configs.Faker() -> go run . -> uncomment this -> go run .
}

func main() {
    db := configs.GetDB()
    db.AutoMigrate(&database.Students{})

    r := routers.SetupRouter()
    r.Run(":3000")
}
