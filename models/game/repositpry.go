package game

import (
	"cloudProject/models/game/dataModel"
	"cloudProject/models/game/dataSource/cacheDS"
	"cloudProject/models/game/dataSource/mysqlDS"
	"cloudProject/pkg/cast"
	"cloudProject/pkg/logger"
	"cloudProject/pkg/redis"
	"context"
	"database/sql"
	"errors"
	"github.com/wcharczuk/go-chart/v2"
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
	GetBestOnPlatform(spanCtx context.Context, platform string, N int) ([]dataModel.GameSales, string, error)
	GetBestOnYear(spanCtx context.Context, year, N int) (games []dataModel.GameSales, errStr string, err error)
	GetBestOnYearAndPlatform(spanCtx context.Context, platform string, year, N int) (games []dataModel.GameSales, errStr string, err error)
	GetGenreBetweenYears(spanCtx context.Context, from, to int) (data []chart.Value, errStr string, err error)
	GetBestOnGenre(spanCtx context.Context, genre string, N int) (games []dataModel.GameSales, errStr string, err error)
	GetProducerOnYears(spanCtx context.Context, p1, p2 string, from, to int) (data map[string][]chart.Value, errStr string, err error)
	GetEuropeMoreThanNorthAmerica(spanCtx context.Context) (games []dataModel.GameSales, errStr string, err error)
	GetGamesSell(spanCtx context.Context, game1, game2 string) (data map[string][]chart.Value, errStr string, err error)
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
func (repo *gameRepository) GetBestOnPlatform(spanCtx context.Context, platform string, N int) (games []dataModel.GameSales, errStr string, err error) {
	traceID := logger.GetTraceIDFromContext(spanCtx)
	games, err = repo.mysqlDS.GetBestOnPlatform(spanCtx, platform, N)
	if err != nil {
		zap.L().Error("get by platform err", zap.String("traceID", traceID), zap.Error(err))
	}
	return games, "01", err
}

func (repo *gameRepository) GetBestOnYear(spanCtx context.Context, year, N int) (games []dataModel.GameSales, errStr string, err error) {
	traceID := logger.GetTraceIDFromContext(spanCtx)
	games, err = repo.mysqlDS.GetBestOnYear(spanCtx, year, N)
	if err != nil {
		zap.L().Error("get by year err", zap.String("traceID", traceID), zap.Error(err))
	}
	return games, "01", err
}
func (repo *gameRepository) GetBestOnYearAndPlatform(spanCtx context.Context, platform string, year, N int) (games []dataModel.GameSales, errStr string, err error) {
	traceID := logger.GetTraceIDFromContext(spanCtx)
	games, err = repo.mysqlDS.GetBestOnYearAndPlatform(spanCtx, platform, year, N)
	if err != nil {
		zap.L().Error("get by year and platform err", zap.String("traceID", traceID), zap.Error(err))
	}
	return games, "01", err
}
func (repo *gameRepository) GetBestOnGenre(spanCtx context.Context, genre string, N int) (games []dataModel.GameSales, errStr string, err error) {
	traceID := logger.GetTraceIDFromContext(spanCtx)
	games, err = repo.mysqlDS.GetBestOnGenre(spanCtx, genre, N)
	if err != nil {
		zap.L().Error("get by genre err", zap.String("traceID", traceID), zap.Error(err))
	}
	return games, "01", err
}
func (repo *gameRepository) GetGenreBetweenYears(spanCtx context.Context, from, to int) (genres []chart.Value, errStr string, err error) {
	traceID := logger.GetTraceIDFromContext(spanCtx)

	games, err := repo.mysqlDS.GetOnGenres(spanCtx, from, to)
	if err != nil {
		zap.L().Error("get by genre err", zap.String("traceID", traceID), zap.Error(err))
	}
	for _, game := range games {
		genres = append(genres, chart.Value{Label: game.Genre, Value: game.GlobalSales})
	}
	return genres, "01", err
}

func (repo *gameRepository) GetEuropeMoreThanNorthAmerica(spanCtx context.Context) (games []dataModel.GameSales, errStr string, err error) {
	traceID := logger.GetTraceIDFromContext(spanCtx)
	games, err = repo.mysqlDS.GetEuropeMoreThanNorthAmerica(spanCtx)
	if err != nil {
		zap.L().Error("get by genre err", zap.String("traceID", traceID), zap.Error(err))
	}
	return games, "01", err
}

func (repo *gameRepository) GetProducerOnYears(spanCtx context.Context, p1, p2 string, from, to int) (map[string][]chart.Value, string, error) {
	traceID := logger.GetTraceIDFromContext(spanCtx)
	var data1, data2 []chart.Value
	p1data, err := repo.mysqlDS.GetProducerOnYears(spanCtx, p1, from, to)
	if err != nil {
		zap.L().Error("get producer in years err", zap.String("traceID", traceID), zap.Error(err))
		return nil, "01", err
	}
	sort.SliceStable(p1data, func(i, j int) bool {
		return p1data[i].Year < p1data[j].Year
	})
	p2data, err := repo.mysqlDS.GetProducerOnYears(spanCtx, p2, from, to)
	if err != nil {
		zap.L().Error("get producer in years err", zap.String("traceID", traceID), zap.Error(err))
		return nil, "02", err
	}
	sort.SliceStable(p2data, func(i, j int) bool {
		return p2data[i].Year < p2data[j].Year
	})
	for _, d := range p1data {
		year, _ := cast.ToString(d.Year)
		data1 = append(data1, chart.Value{Label: year, Value: d.GlobalSales})
	}
	for _, d := range p2data {
		year, _ := cast.ToString(d.Year)
		data2 = append(data2, chart.Value{Label: year, Value: d.GlobalSales})
	}
	var data = make(map[string][]chart.Value)
	data[p1] = data1
	data[p2] = data2
	return data, "01", err
}

func (repo *gameRepository) GetGamesSell(spanCtx context.Context, game1, game2 string) (map[string][]chart.Value, string, error) {
	traceID := logger.GetTraceIDFromContext(spanCtx)
	var data1, data2 []chart.Value
	g1data, err := repo.mysqlDS.GetGameSells(spanCtx, game1)
	if err != nil {
		zap.L().Error("get game sales err", zap.String("traceID", traceID), zap.Error(err))
		return nil, "01", err
	}
	g2data, err := repo.mysqlDS.GetGameSells(spanCtx, game2)
	if err != nil {
		zap.L().Error("get game sales err", zap.String("traceID", traceID), zap.Error(err))
		return nil, "02", err
	}
	var data = make(map[string][]chart.Value)
	data1 = append(data1, chart.Value{Label: "Global_Sales", Value: g1data.GlobalSales}, chart.Value{Label: "NA_Sales", Value: g1data.NASales},
		chart.Value{Label: "EU_Sales", Value: g1data.EUSales}, chart.Value{Label: "JP_Sales", Value: g1data.JPSales},
		chart.Value{Label: "Other_Sales", Value: g1data.OtherSales})

	data2 = append(data2, chart.Value{Label: "Global_Sales", Value: g2data.GlobalSales}, chart.Value{Label: "NA_Sales", Value: g2data.NASales},
		chart.Value{Label: "EU_Sales", Value: g2data.EUSales}, chart.Value{Label: "JP_Sales", Value: g2data.JPSales},
		chart.Value{Label: "Other_Sales", Value: g2data.OtherSales})

	data[game1] = data1
	data[game2] = data2
	return data, "01", err
}
