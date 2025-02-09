package redis

import (
	"context"
	"log"

	"github.com/SuperMatch/config"
	Redis "github.com/redis/go-redis/v9"
)

var RedisClient *Redis.Client

func CreateRedisClient(config config.Config) error {
	RedisClient = Redis.NewClient(&Redis.Options{
		Addr:     config.RedisConfig.Host + ":" + config.RedisConfig.Port,
		Password: config.RedisConfig.Password,
		DB:       0,
	})

	if Ping := RedisClient.Ping(context.TODO()); Ping.Err() != nil {
		log.Println("error in pinging redis.")
		log.Fatal(Ping.Err().Error())
	}
	return nil
}
