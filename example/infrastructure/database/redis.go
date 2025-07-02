package database

import (
	"fmt"

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
	host := "localhost"
	if h := rs.Config.Get("redis.host"); h != nil {
		host = h.(string)
	}

	port := 6379
	if p := rs.Config.Get("redis.port"); p != nil {
		port = p.(int)
	}

	password := ""
	if pwd := rs.Config.Get("redis.password"); pwd != nil {
		password = pwd.(string)
	}

	db := 0
	if d := rs.Config.Get("redis.db"); d != nil {
		db = d.(int)
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	rs.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return nil
}

func (rs *RedisService) Close() error {
	if rs.client != nil {
		return rs.client.Close()
	}
	return nil
}
