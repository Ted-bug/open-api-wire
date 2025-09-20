package redis

import (
	"github.com/go-redis/redis"
)

type RedisClient struct {
	*redis.Client
}

type Config struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Prefix   string `mapstructure:"prefix"` // 项目前缀
}

func NewRedisClient(c Config) (*RedisClient, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})

	// 测试连接
	_, err := r.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &RedisClient{
		Client: r,
	}, nil
}
