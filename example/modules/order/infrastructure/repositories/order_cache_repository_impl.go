package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"example/modules/order/domain/entities"
	"example/modules/order/domain/interfaces"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type OrderCacheRepositoryImpl struct {
	RedisClient *redis.Client `inject:"RedisClient"`
}

func (r *OrderCacheRepositoryImpl) GetServiceName() string {
	return "OrderCacheRepository"
}

func (r *OrderCacheRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	key := fmt.Sprintf("order:%s", id.String())
	val, err := r.RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get order from cache: %w", err)
	}

	var order entities.Order
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return &order, nil
}

func (r *OrderCacheRepositoryImpl) Set(ctx context.Context, order *entities.Order, expiration time.Duration) error {
	key := fmt.Sprintf("order:%s", order.ID.String())

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	if err := r.RedisClient.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set order in cache: %w", err)
	}

	return nil
}

func (r *OrderCacheRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	key := fmt.Sprintf("order:%s", id.String())
	if err := r.RedisClient.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete order from cache: %w", err)
	}
	return nil
}

func (r *OrderCacheRepositoryImpl) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*entities.Order, error) {
	key := fmt.Sprintf("orders:customer:%s", customerID.String())
	val, err := r.RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get customer orders from cache: %w", err)
	}

	var orders []*entities.Order
	if err := json.Unmarshal([]byte(val), &orders); err != nil {
		return nil, fmt.Errorf("failed to unmarshal customer orders: %w", err)
	}

	return orders, nil
}

func (r *OrderCacheRepositoryImpl) SetByCustomerID(ctx context.Context, customerID uuid.UUID, orders []*entities.Order, expiration time.Duration) error {
	key := fmt.Sprintf("orders:customer:%s", customerID.String())

	data, err := json.Marshal(orders)
	if err != nil {
		return fmt.Errorf("failed to marshal customer orders: %w", err)
	}

	if err := r.RedisClient.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set customer orders in cache: %w", err)
	}

	return nil
}

func (r *OrderCacheRepositoryImpl) DeleteByCustomerID(ctx context.Context, customerID uuid.UUID) error {
	key := fmt.Sprintf("orders:customer:%s", customerID.String())
	if err := r.RedisClient.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete customer orders from cache: %w", err)
	}
	return nil
}

func (r *OrderCacheRepositoryImpl) Clear(ctx context.Context) error {
	iter := r.RedisClient.Scan(ctx, 0, "order:*", 0).Iterator()
	var keysToDelete []string

	for iter.Next(ctx) {
		keysToDelete = append(keysToDelete, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan cache keys: %w", err)
	}

	if len(keysToDelete) > 0 {
		if err := r.RedisClient.Del(ctx, keysToDelete...).Err(); err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}
	}

	return nil
}

var _ interfaces.OrderCacheRepository = (*OrderCacheRepositoryImpl)(nil)
