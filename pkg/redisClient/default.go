/*
 * @Description:
 * @Author: gphper
 * @Date: 2021-08-20 20:46:25
 */
package redisClient

import (
	"cmsApp/configs"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	rdb *redis.Client
)

func Init() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:         configs.App.Redis.Addr,
		Password:     configs.App.Redis.Password, // 没有密码，默认值
		DB:           configs.App.Redis.Db,       // 默认DB 0
		PoolSize:     50,                         // 连接池连接数量
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolTimeout:  30 * time.Second,
	})
	ctx := context.Background()
	err := rdb.Ping(ctx).Err()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		panic(err)
	}
	return nil
}

// 获取redis客户端
func GetRedisClient() *redis.Client {
	return rdb
}
