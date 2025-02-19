package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

// A Redis hash can store field-value pairs associated with a product, making it easy to cache product details.
func getProductByIDHash(db *gorm.DB, id uint) (Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)

	var product Product

	// check if the product is in the Redis hash
	result, err := client.HGetAll(ctx, cacheKey).Result()
	if err == nil && len(result) > 0 {
		product.ID = id
		product.Name = result["name"]
		product.Price, _ = strconv.Atoi(result["price"])

		return product, nil
	}

	// if it doesn't exist in the has, then check db
	if err := db.First(&product, id).Error; err != nil {
		// return an empty product
		return product, err
	}

	// we've now grabbed the product from the db
	// so insert it into the hash
	client.HMSet(ctx, cacheKey, map[string]any{
		"name":  product.Name,
		"price": product.Price,
	})

	// using TimeToLive invalidation we set an expiry
	client.Expire(ctx, cacheKey, 5*time.Minute)

	return product, nil
}

// Store recently access product IDs in a Redis list
func addToRecentProductsList(id uint) {
	// Insert all the specified values at the head of the list stored at key.
	client.LPush(ctx, "recent_products", id)

	// Trim an existing list so that it will contain only the specified range of elements specified.
	client.LTrim(ctx, "recent_products", 0, 9) // keep the 10 most recent products
}

// Write-through caching
// Write-through caching writes updates simultaneously to the cache and database, ensuring both stay in sync.
func createOrUpdateProductWriteThrough(db *gorm.DB, id uint, name string, price int) error {
	// update db:
	product := Product{ID: id, Name: name, Price: price}
	if err := db.Save(&product).Error; err != nil {
		return err
	}

	// update in the cache too - with the most recent data
	cacheKey := fmt.Sprintf("product:%d", id)
	client.HMSet(ctx, cacheKey, map[string]any{
		"name":  name,
		"price": price,
	})

	// set TLL expiry
	client.Expire(ctx, cacheKey, 5*time.Minute)

	return nil
}

// Manual invalidation
// Manual invalidation removes outdated entries from the cache explicitly, such as when data changes.
func invalidateProductCache(id uint) error {
	cacheKey := fmt.Sprintf("product:%d", id)
	_, err := client.Del(ctx, cacheKey).Result() // remove the cache entry
	return err
}

// Event-based
// Event-based invalidation can be used to clear cache entries in response to specific application events, such as a significant data update or a deletion.
func deleteProductEventBased(db *gorm.DB, id uint) error {
	// delete product from db
	if err := db.Delete(&Product{}, id).Error; err != nil {
		return err
	}

	// emit an event to clear the cache
	return invalidateProductCache(id)
}

// Cache-Aside (Lazy Loading):
/*
1. When your application needs to read data from the database, it checks the cache first to determine whether the data is available.

2. If the data is available (a cache hit), the cached data is returned, and the response is issued to the caller.

3. If the data isn’t available (a cache miss), the database is queried for the data. The cache is then populated with the data that is retrieved from the database, and the data is returned to the caller.
*/
func getRecentProducts(db *gorm.DB) ([]Product, error) {
	productIDs, err := client.LRange(ctx, "recent_products", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var products []Product
	for _, idStr := range productIDs {
		id, _ := strconv.Atoi(idStr)
		product, err := getProductByIDHash(db, uint(id))

		if err == nil {
			products = append(products, product)
		}
	}

	return products, nil
}

// Redis transaction for atomic updates - all or nothing
func updateProductWithTransaction(db *gorm.DB, id uint, name string, price int) error {
	cacheKey := fmt.Sprintf("product:%d", id)

	// Start a database transaction
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update the database first
	if err := tx.Model(&Product{}).Where("id = ?", id).Updates(Product{Name: name, Price: price}).Error; err != nil {
		tx.Rollback() // Rollback if DB update fails
		return err
	}

	// Commit the DB transaction before updating Redis
	if err := tx.Commit().Error; err != nil {
		return err
	}

	// Update Redis cache after DB commit
	_, err := client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(ctx, cacheKey, map[string]any{
			"name":  name,
			"price": price,
		})
		pipe.Expire(ctx, cacheKey, time.Minute)
		return nil
	})

	return err // Return Redis error if any
}

// START write-behind caching
// Write-behind caching strategy
// In a Write-behind caching strategy, updates are written to the cache first and asynchronously to the database later.
// func createOrUpdateProductWriteBehind(id uint, name string, price int) error {
// 	// Write data to Redis cache
// 	cacheKey := fmt.Sprintf("product:%d", id)
// 	client.HMSet(ctx, cacheKey, map[string]any{
// 		"name":  name,
// 		"price": price,
// 	})

// 	client.Expire(ctx, cacheKey, time.Minute) // set TTL

// 	// Track this operation in a queue for background processing
// 	// .RPush pushes to the tail of the list stored at the key
// 	client.RPush(ctx, "write-behind-queue", cacheKey)
// 	return nil
// }

// // Background Worker for Database Synchronization
// // A background worker will listen for cache entries in write-behind-queue and write them to the database.
// func writeBehindWorker(db *gorm.DB) {
// 	for {
// 		// pop from the queue
// 		cacheKey, err := client.LPop(ctx, "write-behind-queue").Result()
// 		if err == redis.Nil {
// 			time.Sleep(time.Second) // Sleep if queue is empty
// 			continue
// 		} else if err != nil {
// 			log.Println("Error fetching from queue:", err)
// 			continue
// 		}

// 		// there is something on the queue to process
// 		values, err := client.HGetAll(ctx, cacheKey).Result()
// 		if err != nil || len(values) == 0 {
// 			continue
// 		}

// 		// Write to the database
// 		idStr := strings.TrimPrefix(cacheKey, "product:")
// 		id, _ := strconv.Atoi(idStr)
// 		price, _ := strconv.Atoi(values["price"])

// 		product := Product{ID: uint(id), Name: values["name"], Price: price}
// 		if err := db.Save(&product).Error; err != nil {
// 			log.Println("Error saving to DB:", err)
// 			continue
// 		}

// 		// Optionally, delete the cache entry if the cache is temporary
// 		// client.Del(ctx, cacheKey)
// 	}
// }

// END write-behind caching
