package mysqlDS

import (
	"cloudProject/models/game/dataModel"
	"cloudProject/pkg/logger"
	mysqlDB "cloudProject/pkg/mysql"
	"context"
	"database/sql"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"strings"
)

type mysqlDataSource struct {
	conn   *sql.DB
	tracer opentracing.Tracer
}

var mysqlDS MysqlDS

type MysqlDS interface {
	GetByName(spanCtx context.Context, name string) (games []dataModel.GameSales, err error)
	GetByRank(spanCtx context.Context, rank int) (game dataModel.GameSales, err error)
	GetProducerOnYears(spanCtx context.Context, publisher string, from, to int) (data []dataModel.GameSales, err error)
	GetBestOnPlatform(spanCtx context.Context, platform string, N int) (games []dataModel.GameSales, err error)
	GetBestOnYear(spanCtx context.Context, year, N int) (games []dataModel.GameSales, err error)
	GetBestOnYearAndPlatform(spanCtx context.Context, platform string, year, N int) (games []dataModel.GameSales, err error)
	GetEuropeMoreThanNorthAmerica(spanCtx context.Context) (games []dataModel.GameSales, err error)
	GetBestOnGenre(spanCtx context.Context, genre string, N int) (games []dataModel.GameSales, err error)
	GetOnGenres(spanCtx context.Context, from, to int) (games []dataModel.GameSales, err error)
	GetOnYears(spanCtx context.Context, from, to int) (games []dataModel.GameSales, err error)
	GetGameSells(spanCtx context.Context, game string) (g dataModel.GameSales, err error)
}

func init() {
	mysqlDS = &mysqlDataSource{
		conn:   mysqlDB.GetConnection(),
		tracer: logger.InitJaeger("mysql"),
	}
}

func GetDataSource() MysqlDS {
	return mysqlDS
}

func (mysqlDS *mysqlDataSource) GetByName(spanCtx context.Context, name string) (games []dataModel.GameSales, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get game by name query")
	defer logger.FinishSpan(dbSpan)

	rows, err := mysqlDS.conn.Query("SELECT * FROM vgsales WHERE Name LIKE '%" + name + "%'")
	if err != nil {
		zap.L().Error("select games by name err", zap.String("traceID", traceID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	games = make([]dataModel.GameSales, 0)
	for rows.Next() {
		game := dataModel.GameSales{}
		err = rows.Scan(
			&game.Rank,
			&game.Name,
			&game.Platform,
			&game.Year,
			&game.Genre,
			&game.Publisher,
			&game.NASales,
			&game.EUSales,
			&game.JPSales,
			&game.OtherSales,
			&game.GlobalSales,
		)
		if err != nil {
			logger.JaegerErrorLog(dbSpan, err)
			zap.L().Error("scan err", zap.String("traceID", traceID), zap.Error(err))
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}
func (mysqlDS *mysqlDataSource) GetByRank(spanCtx context.Context, rank int) (game dataModel.GameSales, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get game by rank query")
	defer logger.FinishSpan(dbSpan)

	err = mysqlDS.conn.QueryRow("select * from vgsales where `Rank`= ?", int64(rank)).Scan(
		&game.Rank,
		&game.Name,
		&game.Platform,
		&game.Year,
		&game.Genre,
		&game.Publisher,
		&game.NASales,
		&game.EUSales,
		&game.JPSales,
		&game.OtherSales,
		&game.GlobalSales,
	)
	if err != nil {
		logger.JaegerErrorLog(dbSpan, err)
		zap.L().Error("select games by rank err", zap.String("traceID", traceID), zap.Error(err))
	}
	return game, err
}

func (mysqlDS *mysqlDataSource) GetProducerOnYears(spanCtx context.Context, publisher string, from, to int) (games []dataModel.GameSales, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get game by rank query")
	defer logger.FinishSpan(dbSpan)

	rows, err := mysqlDS.conn.Query("SELECT SUM(Global_Sales),`Year` FROM vgsales WHERE Publisher=? AND `Year` >= ? AND `Year` <= ? Group BY `Year`", publisher, from, to)
	if err != nil {
		logger.JaegerErrorLog(dbSpan, err)
		zap.L().Error("select games by producers year err", zap.String("traceID", traceID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	games = make([]dataModel.GameSales, 0)
	for rows.Next() {
		game := dataModel.GameSales{}
		err = rows.Scan(
			&game.GlobalSales,
			&game.Year)
		games = append(games, game)
	}
	return games, err
}

func (mysqlDS *mysqlDataSource) GetGameSells(spanCtx context.Context, game string) (g dataModel.GameSales, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get game by rank query")
	defer logger.FinishSpan(dbSpan)

	err = mysqlDS.conn.QueryRow("SELECT SUM(Global_Sales),SUM(NA_Sales),SUM(EU_Sales),SUM(JP_Sales),SUM(Other_Sales) FROM vgsales WHERE `Name`=?", game).Scan(
		&g.GlobalSales,
		&g.NASales,
		&g.EUSales,
		&g.JPSales,
		&g.OtherSales)
	if err != nil {
		logger.JaegerErrorLog(dbSpan, err)
		zap.L().Error("select game sales err", zap.String("traceID", traceID), zap.Error(err))
		return dataModel.GameSales{}, err
	}
	return g, err
}
func (mysqlDS *mysqlDataSource) GetBestOnPlatform(spanCtx context.Context, platform string, N int) (games []dataModel.GameSales, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get game best on platform")
	defer logger.FinishSpan(dbSpan)

	rows, err := mysqlDS.conn.Query("SELECT * FROM vgsales WHERE lower(Platform)=? ORDER BY `Rank` LIMIT ?", strings.ToLower(platform), N)
	if err != nil {
		logger.JaegerErrorLog(dbSpan, err)
		zap.L().Error("select games by platform err", zap.String("traceID", traceID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	games = make([]dataModel.GameSales, 0)
	for rows.Next() {
		game := dataModel.GameSales{}
		err = rows.Scan(
			&game.Rank,
			&game.Name,
			&game.Platform,
			&game.Year,
			&game.Genre,
			&game.Publisher,
			&game.NASales,
			&game.EUSales,
			&game.JPSales,
			&game.OtherSales,
			&game.GlobalSales,
		)
		if err != nil {
			logger.JaegerErrorLog(dbSpan, err)
			zap.L().Error("scan err", zap.String("traceID", traceID), zap.Error(err))
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}
func (mysqlDS *mysqlDataSource) GetEuropeMoreThanNorthAmerica(spanCtx context.Context) (games []dataModel.GameSales, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get game EU > NA")
	defer logger.FinishSpan(dbSpan)

	rows, err := mysqlDS.conn.Query("SELECT * FROM vgsales WHERE NA_Sales < EU_Sales")
	if err != nil {
		logger.JaegerErrorLog(dbSpan, err)
		zap.L().Error("select games EU > NA err", zap.String("traceID", traceID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	games = make([]dataModel.GameSales, 0)
	for rows.Next() {
		game := dataModel.GameSales{}
		err = rows.Scan(
			&game.Rank,
			&game.Name,
			&game.Platform,
			&game.Year,
			&game.Genre,
			&game.Publisher,
			&game.NASales,
			&game.EUSales,
			&game.JPSales,
			&game.OtherSales,
			&game.GlobalSales,
		)
		if err != nil {
			logger.JaegerErrorLog(dbSpan, err)
			zap.L().Error("scan err", zap.String("traceID", traceID), zap.Error(err))
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}
func (mysqlDS *mysqlDataSource) GetBestOnYear(spanCtx context.Context, year, N int) (games []dataModel.GameSales, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get game best on year")
	defer logger.FinishSpan(dbSpan)

	rows, err := mysqlDS.conn.Query("SELECT * FROM vgsales WHERE `Year`=? ORDER BY Rank ASC LIMIT ?", year, N)
	if err != nil {
		logger.JaegerErrorLog(dbSpan, err)
		zap.L().Error("select games by platform err", zap.String("traceID", traceID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	games = make([]dataModel.GameSales, 0)
	for rows.Next() {
		game := dataModel.GameSales{}
		err = rows.Scan(
			&game.Rank,
			&game.Name,
			&game.Platform,
			&game.Year,
			&game.Genre,
			&game.Publisher,
			&game.NASales,
			&game.EUSales,
			&game.JPSales,
			&game.OtherSales,
			&game.GlobalSales,
		)
		if err != nil {
			logger.JaegerErrorLog(dbSpan, err)
			zap.L().Error("scan err", zap.String("traceID", traceID), zap.Error(err))
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}
func (mysqlDS *mysqlDataSource) GetBestOnYearAndPlatform(spanCtx context.Context, platform string, year, N int) (games []dataModel.GameSales, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get game best on year")
	defer logger.FinishSpan(dbSpan)

	rows, err := mysqlDS.conn.Query("SELECT * FROM vgsales WHERE `platform`=? AND `Year`=?   ORDER BY Rank ASC LIMIT ?", platform, year, N)
	if err != nil {
		logger.JaegerErrorLog(dbSpan, err)
		zap.L().Error("select games by platform err", zap.String("traceID", traceID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	games = make([]dataModel.GameSales, 0)
	for rows.Next() {
		game := dataModel.GameSales{}
		err = rows.Scan(
			&game.Rank,
			&game.Name,
			&game.Platform,
			&game.Year,
			&game.Genre,
			&game.Publisher,
			&game.NASales,
			&game.EUSales,
			&game.JPSales,
			&game.OtherSales,
			&game.GlobalSales,
		)
		if err != nil {
			logger.JaegerErrorLog(dbSpan, err)
			zap.L().Error("scan err", zap.String("traceID", traceID), zap.Error(err))
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}
func (mysqlDS *mysqlDataSource) GetBestOnGenre(spanCtx context.Context, genre string, N int) (games []dataModel.GameSales, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get game best on genre")
	defer logger.FinishSpan(dbSpan)

	rows, err := mysqlDS.conn.Query("SELECT * FROM vgsales WHERE Genre=? ORDER BY Rank ASC LIMIT ?", genre, N)
	if err != nil {
		logger.JaegerErrorLog(dbSpan, err)
		zap.L().Error("select games by genre err", zap.String("traceID", traceID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	games = make([]dataModel.GameSales, 0)
	for rows.Next() {
		game := dataModel.GameSales{}
		err = rows.Scan(
			&game.Rank,
			&game.Name,
			&game.Platform,
			&game.Year,
			&game.Genre,
			&game.Publisher,
			&game.NASales,
			&game.EUSales,
			&game.JPSales,
			&game.OtherSales,
			&game.GlobalSales,
		)
		if err != nil {
			logger.JaegerErrorLog(dbSpan, err)
			zap.L().Error("scan err", zap.String("traceID", traceID), zap.Error(err))
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}
func (mysqlDS *mysqlDataSource) GetOnGenres(spanCtx context.Context, from, to int) (games []dataModel.GameSales, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get game best on genre")
	defer logger.FinishSpan(dbSpan)

	rows, err := mysqlDS.conn.Query("SELECT SUM(Global_Sales),Genre FROM vgsales WHERE `Year` >= ? AND `Year` <= ? Group BY Genre", from, to)
	if err != nil {
		logger.JaegerErrorLog(dbSpan, err)
		zap.L().Error("select games by genre err", zap.String("traceID", traceID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	games = make([]dataModel.GameSales, 0)
	for rows.Next() {
		game := dataModel.GameSales{}
		err = rows.Scan(
			&game.GlobalSales,
			&game.Genre,
		)
		if err != nil {
			logger.JaegerErrorLog(dbSpan, err)
			zap.L().Error("scan err", zap.String("traceID", traceID), zap.Error(err))
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

func (mysqlDS *mysqlDataSource) GetOnYears(spanCtx context.Context, from, to int) (games []dataModel.GameSales, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get game best on genre")
	defer logger.FinishSpan(dbSpan)

	rows, err := mysqlDS.conn.Query("SELECT SUM(Global_Sales),`Year` FROM vgsales  WHERE `Year` >= ? AND `Year` <= ? GROUP BY `Year`", from, to)
	if err != nil {
		logger.JaegerErrorLog(dbSpan, err)
		zap.L().Error("select games by year err", zap.String("traceID", traceID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	games = make([]dataModel.GameSales, 0)
	for rows.Next() {
		game := dataModel.GameSales{}
		err = rows.Scan(
			&game.GlobalSales,
			&game.Year,
		)
		if err != nil {
			logger.JaegerErrorLog(dbSpan, err)
			zap.L().Error("scan err", zap.String("traceID", traceID), zap.Error(err))
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}
