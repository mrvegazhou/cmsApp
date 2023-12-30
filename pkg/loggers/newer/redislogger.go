package newer

import (
	"cmsApp/pkg/redisClient"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisLogger struct {
	Client *redis.Client
	Path   string
}

func NewRedisLogger(path string) *RedisLogger {

	return &RedisLogger{
		Client: redisClient.GetRedisClient(),
		Path:   path,
	}
}

func (logger *RedisLogger) Info(msg string, info map[string]string) {
	info["level"] = "info"
	info["msg"] = msg
	info["ts"] = time.Now().String()

	str, _ := json.Marshal(info)
	time := time.Now().Format("20060102")
	err := logger.Client.LPush(context.Background(), "logs:"+time+":"+logger.Path+":info", string(str)).Err()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key does not exists")
		}
	}
}

func (logger *RedisLogger) Error(msg string, info map[string]string) {
	info["level"] = "error"
	info["msg"] = msg
	info["ts"] = time.Now().String()

	str, _ := json.Marshal(info)
	time := time.Now().Format("20060102")

	logger.Client.LPush(context.Background(), "logs:"+time+":"+logger.Path+":error", string(str))
}
