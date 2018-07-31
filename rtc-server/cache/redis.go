package cache

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/gomodule/redigo/redis"
)

type redisConfig struct {
	Address  string
	Password string
	DB       int
}

type redisClient struct {
	pool *redis.Pool
}

func newRedisClient(c *redisConfig) (*redisClient, error) {
	client := &redisClient{}
	client.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", c.Address) },
	}
	return client, nil
}

func (r *redisClient) Get(key string, container interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()
	objStr, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return errors.NewServerError(http.StatusInternalServerError, fmt.Sprintf("could not set the value for key %s: %v", key, err))
	}
	if err = r.unmarshal(objStr, container); err != nil {
		return err
	}
	return nil
}

func (r *redisClient) Set(key string, val interface{}) error {
	conn := r.pool.Get()
	defer conn.Close()
	valStr, err := r.marshal(val)
	if err != nil {
		return err
	}
	val, err = conn.Do("SET", key, valStr, "NX")
	if err != nil {
		return errors.NewServerError(http.StatusInternalServerError, fmt.Sprintf("could not set the value for key %s: %v", key, err))
	}
	if val == nil {
		return errors.NewServerError(http.StatusConflict, fmt.Sprintf("key %s eixsts already", key))
	}
	return nil
}

func (r *redisClient) Remove(key string) error {
	conn := r.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("DEL", key); err != nil {
		return errors.NewServerError(http.StatusInternalServerError, fmt.Sprintf("error occurred when deleting a key %s: %v", key, err))
	}
	return nil
}

func (r *redisClient) marshal(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", errors.NewServerError(http.StatusInternalServerError, fmt.Sprintf("could not marshal a value %v: %v", v, err))
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (r *redisClient) unmarshal(v string, container interface{}) error {
	objB, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return errors.NewServerError(http.StatusInternalServerError, fmt.Sprintf("could not unmarshal the received value"))
	}
	if err = json.Unmarshal(objB, container); err != nil {
		return errors.NewServerError(http.StatusInternalServerError, fmt.Sprintf("could not unmarshal the fetched object: %v", err))
	}
	return nil
}
