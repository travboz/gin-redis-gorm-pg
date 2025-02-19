package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ctx    = context.Background()
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
)

// func NewDB(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
// 	db, err := sql.Open("postgres", addr)

// 	if err != nil {
// 		return nil, err
// 	}

// 	db.SetMaxOpenConns(maxOpenConns)

// 	duration, err := time.ParseDuration(maxIdleTime)
// 	if err != nil {
// 		return nil, err
// 	}

// 	db.SetConnMaxIdleTime(duration)
// 	db.SetMaxIdleConns(maxIdleConns)

// 	// if it takes more than 5 seconds to connect to the db, we timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if err = db.PingContext(ctx); err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }

func initDB() *gorm.DB {
	// dsn := "host=localhost user=admin password=adminpass dbname=gredis port=5432 sslmode=disable TimeZone=Australia/Sydney"
	dsn := "postgres://admin:adminpass@localhost/gredis?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Product{})
	db.Create(&Product{Name: "Product1", Price: 100})

	return db
}

func initRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return client

}
