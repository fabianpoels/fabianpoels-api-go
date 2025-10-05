package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

type Service struct {
	C *gin.Context
}

func cacheClient() redis.Client {
	if redisClient == nil {
		CacheConnect()
	}
	return *redisClient
}

func CacheConnect() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	log.Printf("Connecting to redis client at: %s", fmt.Sprintf("%s:%s", host, port))
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println(err)
		log.Fatal("⛒ Connection Failed to Cache")
		log.Fatal(err)
	}
	defer cancel()

	log.Println("Connected to cache")

	redisClient = client
}

func (service Service) Set(key string, val interface{}, age time.Duration) error {
	return cacheClient().Set(service.C, key, val, age).Err()
}

func (service Service) Get(key string) (string, error) {
	return cacheClient().Get(service.C, key).Result()
}

func (service Service) GetByes(key string) ([]byte, error) {
	return cacheClient().Get(service.C, key).Bytes()
}

func (service Service) Del(key string) (int64, error) {
	return cacheClient().Del(service.C, key).Result()
}
