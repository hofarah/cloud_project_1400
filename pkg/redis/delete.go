package redis

import (
	"context"
	"go.uber.org/zap"
)

func Delete(key string) (err error) {
	result := getConnection().Del(context.Background(), key)
	if err = result.Err(); err != nil {
		zap.L().Error("delete cache err", zap.Error(err))
	}
	return err
}
