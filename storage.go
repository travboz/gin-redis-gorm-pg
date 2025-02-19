package main

import (
	"gorm.io/gorm"
)

type Storage interface {
	getProduct(id uint) (Product, error)
	createProduct(p Product) error
	deleteProduct(id uint) error
	updateProduct(id uint, name string, price int) error
}

type PGStorage struct {
	db *gorm.DB
}

func NewPGStorage(db *gorm.DB) *PGStorage {
	return &PGStorage{
		db: db,
	}
}

func (pg *PGStorage) getProduct(id uint) (Product, error) {
	// if it doesn't exist in the has, then check db
	var product Product

	if err := pg.db.First(&product, id).Error; err != nil {
		// return an empty product
		return product, err
	}

	return product, nil
}

func (pg *PGStorage) createProduct(p Product) error {
	// update db:
	if err := pg.db.Save(&p).Error; err != nil {
		return err
	}

	return nil
}

func (pg *PGStorage) deleteProduct(id uint) error {
	// delete product from db
	if err := pg.db.Delete(&Product{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (pg *PGStorage) updateProduct(id uint, name string, price int) error {
	// Start a database transaction
	tx := pg.db.Begin()
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

	return nil
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
