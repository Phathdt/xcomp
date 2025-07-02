package database

import (
	"xcomp"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Config *xcomp.ConfigService `inject:"ConfigService"`
	client *redis.Client
}

func (rs *RedisService) GetServiceName() string {
	return "RedisClient"
}

func (rs *RedisService) GetClient() *redis.Client {
	return rs.client
}

func (rs *RedisService) Initialize() error {
	redisURL := rs.Config.GetString("redis.url", "redis://localhost:6379/0")

	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}

	rs.client = redis.NewClient(options)

	return nil
}

func (rs *RedisService) Close() error {
	if rs.client != nil {
		return rs.client.Close()
	}
	return nil
}
