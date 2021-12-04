package user

import (
	"cloudProject/models/user/dataModel"
	"cloudProject/models/user/dataSource/cacheDS"
	"cloudProject/models/user/dataSource/mysqlDS"
	"cloudProject/pkg/jwt"
	"cloudProject/pkg/logger"
	"cloudProject/pkg/redis"
	"cloudProject/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"go.uber.org/zap"
)

type userRepository struct {
	mysqlDS mysqlDS.MysqlDS
	cacheDS cacheDS.RedisDS
}

var userNotFoundErr = sql.ErrNoRows

var Repo Repository

type Repository interface {
	SignUPUser(ctx context.Context, username string) (dataModel.User, string, string, error)
	GetUserByUserName(ctx context.Context, username string) (dataModel.User, bool, string, error)
}

func init() {
	Repo = &userRepository{
		mysqlDS: mysqlDS.GetDataSource(),
		cacheDS: cacheDS.GetDataSource(),
	}
}

func (repo *userRepository) SignUPUser(ctx context.Context, username string) (dataModel.User, string, string, error) {
	traceID := logger.GetTraceIDFromContext(ctx)
	_, exist, _, err := repo.GetUserByUserName(ctx, username)
	if err != nil {
		zap.L().Error("get user by username err", zap.String("traceID", traceID), zap.Any("username", username), zap.Error(err))
		return dataModel.User{}, "01", "", err
	}
	if exist {
		return dataModel.User{}, "02", "", errors.New("repetitiousUserName")
	}
	token, err := jwt.CreateToken(traceID)
	if err != nil {
		zap.L().Error("creatToken err", zap.String("traceID", traceID), zap.Error(err))
		return dataModel.User{}, "03", "", err
	}
	secret := utils.MD5(username + utils.RandString(10))
	user := dataModel.User{
		UserName: username,
		Secret:   secret,
	}
	err = repo.mysqlDS.Insert(ctx, user)
	if err != nil {
		zap.L().Error("insert User err", zap.String("traceID", traceID), zap.Error(err))
	}
	return user, "04", token, err
}
func (repo *userRepository) GetUserByUserName(ctx context.Context, username string) (user dataModel.User, exist bool, errStr string, err error) {
	traceID := logger.GetTraceIDFromContext(ctx)
	user, err = repo.cacheDS.GetFromCacheByName(username)
	if err != nil {
		if err != redis.NotFoundInCacheErr {
			zap.L().Error("get from cache err", zap.String("traceID", traceID), zap.Error(err))
		}
		user, err = repo.mysqlDS.GetByUserName(ctx, username)
		if err != nil {
			if err == userNotFoundErr {
				return user, false, "", nil
			}
			zap.L().Error("get user from db err", zap.String("traceID", traceID), zap.Error(err))
			return dataModel.User{}, false, "01", err
		}
		err = repo.cacheDS.SetToCache(user)
		if err != nil {
			zap.L().Error("set to cache err", zap.String("traceID", traceID), zap.Error(err))
		}
	}
	return user, true, "", nil
}
