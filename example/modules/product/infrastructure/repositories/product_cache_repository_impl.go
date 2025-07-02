package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"example/modules/product/domain/entities"
	"example/modules/product/domain/interfaces"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ProductCacheRepositoryImpl struct {
	RedisClient *redis.Client `inject:"RedisClient"`
}

func (r *ProductCacheRepositoryImpl) GetServiceName() string {
	return "ProductCacheRepository"
}

func (r *ProductCacheRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	if r.RedisClient == nil {
		log.Printf("Redis client is nil, skipping cache get")
		return nil, nil
	}

	key := r.getProductKey(id)
	log.Printf("Attempting to get product from cache with key: %s", key)

	val, err := r.RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Product not found in cache: %s", key)
			return nil, nil
		}
		log.Printf("Error getting product from cache: %v", err)
		return nil, fmt.Errorf("failed to get product from cache: %w", err)
	}

	log.Printf("Found product in cache: %s", key)
	var product entities.Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		log.Printf("Error unmarshaling product from cache: %v", err)
		return nil, fmt.Errorf("failed to unmarshal product from cache: %w", err)
	}

	return &product, nil
}

func (r *ProductCacheRepositoryImpl) Set(ctx context.Context, product *entities.Product, ttl time.Duration) error {
	key := r.getProductKey(product.ID)
	productJSON, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product for cache: %w", err)
	}

	if err := r.RedisClient.Set(ctx, key, productJSON, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set product in cache: %w", err)
	}

	return nil
}

func (r *ProductCacheRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	key := r.getProductKey(id)
	if err := r.RedisClient.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete product from cache: %w", err)
	}
	return nil
}

func (r *ProductCacheRepositoryImpl) Clear(ctx context.Context) error {
	iter := r.RedisClient.Scan(ctx, 0, "product:*", 0).Iterator()
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

func (r *ProductCacheRepositoryImpl) getProductKey(id uuid.UUID) string {
	return fmt.Sprintf("product:%s", id.String())
}

var _ interfaces.ProductCacheRepository = (*ProductCacheRepositoryImpl)(nil)
