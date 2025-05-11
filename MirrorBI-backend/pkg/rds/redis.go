package rds

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"mrbi/config"
	"mrbi/pkg/redlock"
)

var redisClient *redis.Client

func init() {
	var config = config.LoadConfig()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Rds.Host, config.Rds.Port),
		Password: fmt.Sprintf("%s", config.Rds.Password),
		DB:       0,
	})
	// 通过 Ping 命令测试连接
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis 连接失败:", err)
	} else {
		fmt.Println("Redis 连接成功！")
	}
	status := redisClient.Set(ctx, "LinkTest", "Success", 0).Err()
	if status != nil {
		log.Fatalf("Redis Set Error:", status)
	}
	//初始化redis分布式锁
	redlock.InitRedSync(redisClient)
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func IsNilErr(err error) bool {
	return err == redis.Nil
}
