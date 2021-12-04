package cacheDS

import (
	"cloudProject/models/game/dataModel"
	"cloudProject/pkg/redis"
	"time"
)

type redisDataSource struct {
	DefaultCachePrefix string
	DefaultExpireTime  time.Duration
}

var cacheDS RedisDS

type RedisDS interface {
	GetFromCacheByName(name string) (model dataModel.GameSales, err error)
	SetToCache(model dataModel.GameSales) error
}

func init() {
	cacheDS = &redisDataSource{
		DefaultCachePrefix: "game_",
		DefaultExpireTime:  time.Hour * 10,
	}
}

func GetDataSource() RedisDS {
	return cacheDS
}

func (cache *redisDataSource) GetFromCacheByName(name string) (model dataModel.GameSales, err error) {
	err = redis.Get(cache.prepareKey(name), &model)
	return model, err
}
func (cache *redisDataSource) SetToCache(model dataModel.GameSales) error {
	return redis.Set(cache.prepareKey(model.Name), model, time.Hour*12)
}
func (cache *redisDataSource) prepareKey(name string) string {
	return cache.DefaultCachePrefix + name
}
