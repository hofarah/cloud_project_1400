package redis

import (
	"context"
	redisClient "github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

func Get(key string, data interface{}) error {
	bytes, err := getConnection().Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redisClient.Nil {
			return NotFoundInCacheErr
		}
		zap.L().Error("get from cache err", zap.Error(err))
		return err
	}
	err = jsoniter.Unmarshal(bytes, data)
	if err != nil {
		zap.L().Error("err when unmarshal", zap.Error(err))
	}
	return err
}
