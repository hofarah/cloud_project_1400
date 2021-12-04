package game

import (
	"cloudProject/models/game/dataModel"
	"cloudProject/models/game/dataSource/cacheDS"
	"cloudProject/models/game/dataSource/mysqlDS"
	"cloudProject/pkg/logger"
	"cloudProject/pkg/redis"
	"context"
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"sort"
)

type gameRepository struct {
	mysqlDS mysqlDS.MysqlDS
	cacheDS cacheDS.RedisDS
}

var gameNotFoundErr = sql.ErrNoRows

var Repo Repository

type Repository interface {
	GetByRank(spanCtx context.Context, rank int) (dataModel.GameSales, bool, string, error)
	GetByName(spanCtx context.Context, name string) ([]dataModel.GameSales, string, error)
}

func init() {
	Repo = &gameRepository{
		mysqlDS: mysqlDS.GetDataSource(),
		cacheDS: cacheDS.GetDataSource(),
	}
}

func (repo *gameRepository) GetByRank(spanCtx context.Context, rank int) (game dataModel.GameSales, exist bool, errStr string, err error) {
	traceID := logger.GetTraceIDFromContext(spanCtx)
	game, err = repo.cacheDS.GetFromCacheByRank(rank)
	if err != nil {
		if err != redis.NotFoundInCacheErr {
			zap.L().Error("get from cache err", zap.String("traceID", traceID), zap.Error(err))
		}
		game, err = repo.mysqlDS.GetByRank(spanCtx, rank)
		if err != nil {
			if err == gameNotFoundErr {
				return dataModel.GameSales{}, false, "", errors.New("gameNotFound")
			}
			zap.L().Error("get by rank err", zap.String("traceID", traceID), zap.Error(err))
			return dataModel.GameSales{}, false, "01", err
		}
		err = repo.cacheDS.SetToCache(game.Rank, game)
		if err != nil {
			zap.L().Error("set to cache err", zap.String("traceID", traceID), zap.Error(err))
		}
	}
	return game, true, "", nil
}
func (repo *gameRepository) GetByName(spanCtx context.Context, name string) ([]dataModel.GameSales, string, error) {
	traceID := logger.GetTraceIDFromContext(spanCtx)
	games, err := repo.mysqlDS.GetByName(spanCtx, name)
	if err != nil {
		zap.L().Error("get by name err", zap.String("traceID", traceID), zap.Error(err))
	}
	sort.Slice(games, func(i, j int) bool {
		return games[i].Rank > games[j].Rank
	})
	return games, "01", err
}
