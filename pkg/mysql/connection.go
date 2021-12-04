package mysqlDB

import (
	"cloudProject/pkg/cast"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"os"
)

var connection *sql.DB

func setMysqlClient() {
	var err error
	URI := os.Getenv("MYSQL")
	if URI == "" {
		panic("MYSQL not found")
	}
	zap.L().Info("mysql url at env ", zap.Any("url:", URI))
	USER := os.Getenv("MYSQL_USER")
	PASSWORD := os.Getenv("MYSQL_PASSWORD")
	DB := os.Getenv("MYSQL_DB")
	if USER == "" || PASSWORD == "" || DB == "" {
		zap.L().Error("mysql param required")
		panic("some MySQL param (MYSQL_USER_MYSQL_PASSWORD_MYSQL_DB) are missed")
	}
	MAXCONN, _ := cast.ToInt(os.Getenv("MYSQL_MAXCONN"))
	if MAXCONN == 0 {
		MAXCONN = 10000
	}
	MINCONN, _ := cast.ToInt(os.Getenv("MYSQL_MINCONN"))
	if MINCONN == 0 {
		MINCONN = 1000
	}
	connection, err = sql.Open("mysql", USER+":"+PASSWORD+"@tcp("+URI+")/"+DB)
	if err != nil {
		zap.L().Error("mysql connection error", zap.Error(err))
		panic(err)
	}
	connection.SetMaxIdleConns(MINCONN)
	connection.SetMaxOpenConns(MAXCONN)
	connection.SetConnMaxLifetime(0)
}

func GetConnection() *sql.DB {
	if connection == nil {
		setMysqlClient()
	}
	return connection
}
