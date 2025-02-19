package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	getProductByIDHash(id uint) (Product, error)
	addToRecentProductsList(id uint)
	createOrUpdateProductWriteThrough(id uint, name string, price int) error
	invalidateProductCache(id uint) error
	deleteProductEventBased(id uint) error
	getRecentProducts() ([]Product, error)
	updateProductWithTransaction(id uint, name string, price int) error
}

type RedisCache struct {
	dbStore Storage
	client  *redis.Client
}

func NewRedisCache(db Storage, addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisCache{
		dbStore: db,
		client:  client,
	}
}

// A Redis hash can store field-value pairs associated with a product, making it easy to cache product details.
func (r *RedisCache) getProductByIDHash(id uint) (Product, error) {
	ctx := context.Background()

	cacheKey := fmt.Sprintf("product:%d", id)

	var product Product

	// check if the product is in the Redis hash
	result, err := r.client.HGetAll(ctx, cacheKey).Result()
	if err == nil && len(result) > 0 {
		product.ID = id
		product.Name = result["name"]
		product.Price, _ = strconv.Atoi(result["price"])

		return product, nil
	}

	// if it doesn't exist in the has, then check db
	exists_prod, err := r.dbStore.getProduct(id)
	if err != nil {
		// return an empty product
		return product, err
	}

	// we've now grabbed the product from the db
	// so insert it into the hash
	r.client.HMSet(ctx, cacheKey, map[string]any{
		"name":  exists_prod.Name,
		"price": exists_prod.Price,
	})

	// using TimeToLive invalidation we set an expiry
	r.client.Expire(ctx, cacheKey, 5*time.Minute)

	return exists_prod, nil
}

// Store recently access product IDs in a Redis list
func (r *RedisCache) addToRecentProductsList(id uint) {
	ctx := context.Background()

	// Insert all the specified values at the head of the list stored at key.
	r.client.LPush(ctx, "recent_products", id)

	// Trim an existing list so that it will contain only the specified range of elements specified.
	r.client.LTrim(ctx, "recent_products", 0, 9) // keep the 10 most recent products
}

// Write-through caching
// Write-through caching writes updates simultaneously to the cache and database, ensuring both stay in sync.
func (r *RedisCache) createOrUpdateProductWriteThrough(id uint, name string, price int) error {
	ctx := context.Background()

	// update db:
	product := Product{ID: id, Name: name, Price: price}

	if err := r.dbStore.createProduct(product); err != nil {
		return err
	}

	// update in the cache too - with the most recent data
	cacheKey := fmt.Sprintf("product:%d", id)
	r.client.HMSet(ctx, cacheKey, map[string]any{
		"name":  name,
		"price": price,
	})

	// set TLL expiry
	r.client.Expire(ctx, cacheKey, 5*time.Minute)

	return nil
}

// Manual invalidation
// Manual invalidation removes outdated entries from the cache explicitly, such as when data changes.
func (r *RedisCache) invalidateProductCache(id uint) error {
	ctx := context.Background()

	cacheKey := fmt.Sprintf("product:%d", id)
	_, err := r.client.Del(ctx, cacheKey).Result() // remove the cache entry
	return err
}

// Event-based
// Event-based invalidation can be used to clear cache entries in response to specific application events, such as a significant data update or a deletion.
func (r *RedisCache) deleteProductEventBased(id uint) error {
	// delete product from db
	if err := r.dbStore.deleteProduct(id); err != nil {
		return err
	}
	// emit an event to clear the cache
	return r.invalidateProductCache(id)
}

// Cache-Aside (Lazy Loading):
/*
1. When your application needs to read data from the database, it checks the cache first to determine whether the data is available.

2. If the data is available (a cache hit), the cached data is returned, and the response is issued to the caller.

3. If the data isn’t available (a cache miss), the database is queried for the data. The cache is then populated with the data that is retrieved from the database, and the data is returned to the caller.
*/
func (r *RedisCache) getRecentProducts() ([]Product, error) {
	ctx := context.Background()

	productIDs, err := r.client.LRange(ctx, "recent_products", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var products []Product
	for _, idStr := range productIDs {
		id, _ := strconv.Atoi(idStr)
		product, err := r.getProductByIDHash(uint(id))

		if err == nil {
			products = append(products, product)
		}
	}

	return products, nil
}

// Redis transaction for atomic updates - all or nothing
func (r *RedisCache) updateProductWithTransaction(id uint, name string, price int) error {
	ctx := context.Background()

	cacheKey := fmt.Sprintf("product:%d", id)

	// perform transaction on db
	if err := r.dbStore.updateProduct(id, name, price); err != nil {
		return err
	}

	// Update Redis cache after DB commit
	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(ctx, cacheKey, map[string]any{
			"name":  name,
			"price": price,
		})
		pipe.Expire(ctx, cacheKey, time.Minute)
		return nil
	})

	return err // Return Redis error if any
}
