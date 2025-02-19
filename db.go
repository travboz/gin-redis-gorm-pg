package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
