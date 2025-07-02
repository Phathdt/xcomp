package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"example/modules/customer/domain/entities"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CustomerCacheRepositoryImpl struct {
	RedisClient *redis.Client `inject:"RedisClient"`
}

func (r *CustomerCacheRepositoryImpl) GetServiceName() string {
	return "CustomerCacheRepositoryImpl"
}

func (r *CustomerCacheRepositoryImpl) Set(ctx context.Context, key string, customer *entities.Customer, ttl time.Duration) error {
	data, err := json.Marshal(customer)
	if err != nil {
		return err
	}

	return r.RedisClient.Set(ctx, key, data, ttl).Err()
}

func (r *CustomerCacheRepositoryImpl) Get(ctx context.Context, key string) (*entities.Customer, error) {
	data, err := r.RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var customer entities.Customer
	if err := json.Unmarshal([]byte(data), &customer); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerCacheRepositoryImpl) Delete(ctx context.Context, key string) error {
	return r.RedisClient.Del(ctx, key).Err()
}

func (r *CustomerCacheRepositoryImpl) GetCustomerCacheKey(id uuid.UUID) string {
	return fmt.Sprintf("customer:id:%s", id.String())
}

func (r *CustomerCacheRepositoryImpl) GetCustomerUsernameCacheKey(username string) string {
	return fmt.Sprintf("customer:username:%s", username)
}

func (r *CustomerCacheRepositoryImpl) GetCustomerEmailCacheKey(email string) string {
	return fmt.Sprintf("customer:email:%s", email)
}
