package redis

import (
	"errors"
	"time"
)

const DefaultExpTime = time.Hour * 6

var NotFoundInCacheErr = errors.New("key dose not exist")

func IsNotFoundInCacheErr(err error) bool {
	return errors.Is(err, NotFoundInCacheErr)
}
