package game

import (
	"cloudProject/models/user/dataSource/cacheDS"
	"cloudProject/models/user/dataSource/mysqlDS"
	"database/sql"
)

type userRepository struct {
	mysqlDS mysqlDS.MysqlDS
	cacheDS cacheDS.RedisDS
}

var gameNotFoundErr = sql.ErrNoRows

var Repo Repository

type Repository interface {
}

func init() {
	Repo = &userRepository{
		mysqlDS: mysqlDS.GetDataSource(),
		cacheDS: cacheDS.GetDataSource(),
	}
}
