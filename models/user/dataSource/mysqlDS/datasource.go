package mysqlDS

import (
	"cloudProject/models/user/dataModel"
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
	Insert(spanCtx context.Context, user dataModel.User) error
	GetByUserName(spanCtx context.Context, username string) (user dataModel.User, err error)
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

func (mysqlDS *mysqlDataSource) Insert(spanCtx context.Context, user dataModel.User) error {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "insert user query")
	defer logger.FinishSpan(dbSpan)

	_, err := mysqlDS.conn.Exec("INSERT INTO `users` (username,secret) VALUES (?,?)", user.UserName, user.Secret)
	if err != nil {
		logger.JaegerErrorLog(dbSpan, err)
		zap.L().Error("insert user err", zap.String("traceID", traceID), zap.Error(err), zap.Any("insertField", user))
	}
	return err
}
func (mysqlDS *mysqlDataSource) GetByUserName(spanCtx context.Context, username string) (user dataModel.User, err error) {
	dbSpan, traceID := logger.StartSpan(spanCtx, mysqlDS.tracer, "get user by username query")
	defer logger.FinishSpan(dbSpan)

	err = mysqlDS.conn.QueryRow("SELECT ID,username,secret FROM users WHERE username=?", username).Scan(
		&user.ID,
		&user.UserName,
		&user.Secret,
	)
	if err != nil {
		zap.L().Error("select user by username err", zap.String("traceID", traceID), zap.Error(err))
	}
	return user, err
}
