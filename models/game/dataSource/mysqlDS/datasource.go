package mysqlDS

import (
	"cloudProject/models/game/dataModel"
	"cloudProject/pkg/logger"
	mysqlDB "cloudProject/pkg/mysql"
	"context"
	"database/sql"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type mysqlDataSource struct {
	conn   *sql.DB
	tracer opentracing.Tracer
}

var mysqlDS MysqlDS

type MysqlDS interface {
	GetByName(spanCtx context.Context, name string) (games []dataModel.GameSales, err error)
	GetByRank(spanCtx context.Context, rank int) (game dataModel.GameSales, err error)
	GetBestOnPlatform(spanCtx context.Context, platform string, N int) (games []dataModel.GameSales, err error)
	GetBestOnYear(spanCtx context.Context, year, N int) (games []dataModel.GameSales, err error)
	GetBestOnGenre(spanCtx context.Context, genre string, N int) (games []dataModel.GameSales, err error)
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

	rows, err := mysqlDS.conn.Query("SELECT Rank,`Name`,Platform,`Year`,Genre,Publisher,NA_Sales,EU_Sales,JP_Sales,Other_Sales,Global_Sales FROM vgsales WHERE Name LIKE '%" + name + "%'")
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

	err = mysqlDS.conn.QueryRow("SELECT * FROM vgsales WHERE Rank=?", rank).Scan(
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
func (mysqlDS *mysqlDataSource) GetBestOnPlatform(spanCtx context.Context, platform string, N int) (games []dataModel.GameSales, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get game best on platform")
	defer logger.FinishSpan(dbSpan)

	rows, err := mysqlDS.conn.Query("SELECT * FROM vgsales WHERE Platform=? ORDER BY Rank ASC LIMIT ?", platform, N)
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
