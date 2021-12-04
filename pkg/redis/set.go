package redis

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"time"
)

func Set(key string, data interface{}, expireTime ...time.Duration) error {
	bytes, err := jsoniter.Marshal(data)
	if err != nil {
		zap.L().Error("marshal err", zap.Error(err))
		return err
	}
	exp := DefaultExpTime
	if len(expireTime) > 0 {
		exp = expireTime[0]
	}
	result := getConnection().Set(context.Background(), key, bytes, exp)
	if err = result.Err(); err != nil {
		zap.L().Error("set to redis err", zap.Error(err))
	}
	return err
}
