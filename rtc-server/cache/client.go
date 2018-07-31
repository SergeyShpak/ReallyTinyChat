package cache

import (
	"net/http"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/config"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
)

type Client interface {
	Get(key string, container interface{}) error
	Set(key string, val interface{}) error
	Remove(key string) error
}

func NewClient(config *config.Cache) (Client, error) {
	if config.Redis != nil {
		redisConfig := &redisConfig{
			Address:  config.Redis.Address,
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
		}
		return newRedisClient(redisConfig)
	}
	return nil, errors.NewServerError(http.StatusInternalServerError, "no cache configuration found")
}
