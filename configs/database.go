package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB    *gorm.DB
var Redis *redis.Client

func ConnectToDB() error {
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, password, dbname, port)
    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })

    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }

    log.Println("Connected to database")
    return nil
}

func ConnectToRedis() error {
    host := os.Getenv("REDIS_HOST")
    port := os.Getenv("REDIS_PORT")
    password := os.Getenv("REDIS_PASSWORD")

    Redis = redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", host, port),
        Password: password,
        DB:       0,
    })

    ctx := context.Background()
    _, err := Redis.Ping(ctx).Result()
    if err != nil {
        return fmt.Errorf("failed to connect to Redis: %w", err)
    }

    log.Println("Connected to Redis")
    return nil
}

func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database is not initialized")
	}
	return DB
}

func GetRedis() *redis.Client {
	if Redis == nil {
		log.Fatal("Redis is not initialized")
	}
	return Redis
}

