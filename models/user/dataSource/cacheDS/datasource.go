package cacheDS

import (
	"cloudProject/models/user/dataModel"
	"cloudProject/pkg/redis"
	"time"
)

type redisDataSource struct {
	DefaultCachePrefix string
	DefaultExpireTime  time.Duration
}

var cacheDS RedisDS

type RedisDS interface {
	GetFromCacheByName(username string) (model dataModel.User, err error)
	SetToCache(model dataModel.User) error
}

func init() {
	cacheDS = &redisDataSource{
		DefaultCachePrefix: "user_",
		DefaultExpireTime:  time.Hour * 10,
	}
}

func GetDataSource() RedisDS {
	return cacheDS
}

func (cache *redisDataSource) GetFromCacheByName(username string) (model dataModel.User, err error) {
	err = redis.Get(cache.prepareKey(username), &model)
	return model, err
}
func (cache *redisDataSource) SetToCache(model dataModel.User) error {
	return redis.Set(cache.prepareKey(model.UserName), model, time.Hour*12)
}
func (cache *redisDataSource) prepareKey(username string) string {
	return cache.DefaultCachePrefix + username
}
