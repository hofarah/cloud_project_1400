package cacheDS

import (
	"cloudProject/models/game/dataModel"
	"cloudProject/pkg/cast"
	"cloudProject/pkg/redis"
	"time"
)

type redisDataSource struct {
	DefaultCachePrefix   string
	GetByRankCachePrefix string
	DefaultExpireTime    time.Duration
}

var cacheDS RedisDS

type RedisDS interface {
	GetFromCacheByName(name string) (model dataModel.GameSales, err error)
	GetFromCacheByRank(rank int) (model dataModel.GameSales, err error)
	SetToCache(key interface{}, model dataModel.GameSales) error
}

func init() {
	cacheDS = &redisDataSource{
		DefaultCachePrefix:   "game_",
		GetByRankCachePrefix: "game_rank_",
		DefaultExpireTime:    time.Hour * 10,
	}
}

func GetDataSource() RedisDS {
	return cacheDS
}

func (cache *redisDataSource) GetFromCacheByName(name string) (model dataModel.GameSales, err error) {
	err = redis.Get(cache.prepareKey(name), &model)
	return model, err
}
func (cache *redisDataSource) GetFromCacheByRank(rank int) (model dataModel.GameSales, err error) {
	err = redis.Get(cache.prepareKey(rank), &model)
	return model, err
}
func (cache *redisDataSource) SetToCache(key interface{}, model dataModel.GameSales) error {
	return redis.Set(cache.prepareKey(key), model, time.Hour*12)
}
func (cache *redisDataSource) prepareKey(key interface{}) string {
	cacheKey, _ := cast.ToString(key)
	return cache.DefaultCachePrefix + cacheKey
}
