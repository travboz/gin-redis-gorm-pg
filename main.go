package main

import (
	_ "github.com/lib/pq"
)

func main() {
	db := initDB()
	store := NewPGStorage(db)
	redisCache := NewRedisCache(store, "localhost:6379")

	app := application{
		cache: redisCache,
	}

	router := setupRouter(&app)
	router.Run(":8080")
}
