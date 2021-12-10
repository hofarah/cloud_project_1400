package redis

import (
	"cloudProject/pkg/cast"
	"context"
	redisClient "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"os"
	"time"
)

var connection *redisClient.Client

func setClient() {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		//panic("REDIS_URL not found")
	}
	zap.L().Info("REDIS_URL at env " + url)

	poolSize, _ := cast.ToInt(os.Getenv("REDIS_MAXCON"))
	if poolSize == 0 {
		poolSize = 500
	}

	minCon, _ := cast.ToInt(os.Getenv("REDIS_MINCON"))
	if minCon == 0 {
		minCon = 100
	}

	connection = redisClient.NewClient(&redisClient.Options{
		Addr:         url,
		Password:     "",
		DB:           0, // use default db
		PoolSize:     poolSize,
		MinIdleConns: minCon,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := connection.Ping(ctx).Result()
	if err != nil {
		zap.L().Error("redis ping err", zap.Error(err))
		//panic(err)
	}
}

func getConnection() *redisClient.Client {
	if connection == nil {
		setClient()
	}
	return connection
}
