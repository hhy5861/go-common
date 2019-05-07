package common

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

type (
	RedisConfig struct {
		Host              string        `yaml:"host"`
		Port              int           `yaml:"port"`
		Db                int           `yaml:"db"`
		Password          string        `yaml:"password"`
		MaxIdle           int           `yaml:"maxIdle"`
		MaxActive         int           `yaml:"maxActive"`
		IdleTimeout       time.Duration `yaml:"idleTimeout"`
		ConnectTimeout    time.Duration `yaml:"connectTimeout"`
		ReadTimeout       time.Duration `yaml:"readTimeout"`
		WriteTimeout      time.Duration `yaml:"writeTimeout"`
		DefaultExpiration time.Duration `json:"defaultExpiration"`
	}

	RedisStore struct {
		Pool              *redis.Pool
		defaultExpiration time.Duration
		Serialize         ISerialize
	}

	ISerialize interface {
		Deserialize(byt []byte, ptr interface{}) error

		Serialize(value interface{}) ([]byte, error)
	}
)

var (
	redisStore   *RedisStore
	ErrCacheMiss = errors.New("cache: key not found.")
	ErrNotStored = errors.New("cache: not stored.")
)

const (
	DEFAULT = time.Duration(0)
	FOREVER = time.Duration(-1)
)

func GetRedisStore() *RedisStore {
	return redisStore
}

func getRedisConn() redis.Conn {
	return redisStore.Pool.Get()
}

func NewRedisStore(config *RedisConfig, serialize ...ISerialize) *RedisStore {
	var iSerialize ISerialize
	if len(serialize) > 0 {
		iSerialize = serialize[0]
	}

	redisStore = &RedisStore{
		Pool: &redis.Pool{
			MaxIdle:     config.MaxIdle,
			MaxActive:   config.MaxActive,
			IdleTimeout: config.IdleTimeout * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial(
					"tcp", fmt.Sprintf("%s:%d", config.Host, config.Port),
					redis.DialConnectTimeout(time.Duration(config.ConnectTimeout*time.Second)),
					redis.DialReadTimeout(time.Duration(config.ReadTimeout*time.Second)),
					redis.DialWriteTimeout(time.Duration(config.WriteTimeout*time.Second)),
				)

				if err != nil {
					return nil, err
				}

				if len(config.Password) > 0 {
					if _, err := c.Do("AUTH", config.Password); err != nil {
						c.Close()
						return nil, err
					}
				} else {
					if _, err := c.Do("PING"); err != nil {
						c.Close()
						return nil, err
					}
				}

				return c, err
			},

			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if _, err := c.Do("PING"); err != nil {
					return err
				}
				return nil
			},
		},

		defaultExpiration: time.Duration(config.DefaultExpiration * time.Second),
		Serialize:         iSerialize,
	}

	return redisStore
}

func (c *RedisStore) Set(key string, value interface{}, expires time.Duration) error {
	return c.Invoke(getRedisConn().Do, key, value, expires)
}

func (c *RedisStore) Add(key string, value interface{}, expires time.Duration) error {
	conn := getRedisConn()
	if exists(conn, key) {
		return ErrNotStored
	}

	return c.Invoke(conn.Do, key, value, expires)
}

func (c *RedisStore) Replace(key string, value interface{}, expires time.Duration) error {
	conn := getRedisConn()
	if !exists(conn, key) {
		return ErrNotStored
	}

	err := c.Invoke(conn.Do, key, value, expires)
	if value == nil {
		return ErrNotStored
	} else {
		return err
	}
}

func (c *RedisStore) Get(key string, ptrValue interface{}) error {
	conn := getRedisConn()
	defer conn.Close()

	raw, err := conn.Do("GET", key)
	if raw == nil {
		return ErrCacheMiss
	}

	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}

	return c.Serialize.Deserialize(item, ptrValue)
}

func exists(conn redis.Conn, key string) bool {
	retval, _ := redis.Bool(conn.Do("EXISTS", key))

	return retval
}

func (c *RedisStore) Delete(key string) error {
	conn := getRedisConn()
	defer conn.Close()

	if !exists(conn, key) {
		return ErrCacheMiss
	}

	_, err := conn.Do("DEL", key)
	return err
}

func (c *RedisStore) Increment(key string, delta uint64) (uint64, error) {
	conn := getRedisConn()
	defer conn.Close()

	val, err := conn.Do("GET", key)
	if val == nil {
		return 0, ErrCacheMiss
	}

	if err == nil {
		currentVal, err := redis.Int64(val, nil)
		if err != nil {
			return 0, err
		}

		var sum = currentVal + int64(delta)
		_, err = conn.Do("SET", key, sum)
		if err != nil {
			return 0, err
		}

		return uint64(sum), nil
	} else {
		return 0, err
	}
}

func (c *RedisStore) Decrement(key string, delta uint64) (newValue uint64, err error) {
	conn := getRedisConn()
	defer conn.Close()

	if !exists(conn, key) {
		return 0, ErrCacheMiss
	}

	currentVal, err := redis.Int64(conn.Do("GET", key))
	if err == nil && delta > uint64(currentVal) {
		tempint, err := redis.Int64(conn.Do("DECRBY", key, currentVal))
		return uint64(tempint), err
	}

	tempint, err := redis.Int64(conn.Do("DECRBY", key, delta))
	return uint64(tempint), err
}

func (c *RedisStore) Flush() error {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("FLUSHALL")
	return err
}

func (c *RedisStore) Invoke(f func(string, ...interface{}) (interface{}, error),
	key string, value interface{}, expires time.Duration) error {

	switch expires {
	case DEFAULT:
		expires = c.defaultExpiration
	case FOREVER:
		expires = time.Duration(0)
	}

	bady, err := c.Serialize.Serialize(value)
	if err != nil {
		return err
	}

	conn := getRedisConn()
	defer conn.Close()

	if expires > 0 {
		_, err := f("SETEX", key, int32(expires/time.Second), bady)
		return err
	} else {
		_, err := f("SET", key, bady)
		return err
	}
}

func (c *RedisStore) ExecCommand(command string, args ...interface{}) (interface{}, error) {
	conn := getRedisConn()
	defer conn.Close()

	return conn.Do(command, args...)
}
