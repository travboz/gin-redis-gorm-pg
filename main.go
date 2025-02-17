package main

import (
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/travboz/gorm-redis-gin-api/internal/env"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dbc := dbConfig{
		addr:         env.GetString("DB_ADDR", "postgres://admin:adminpass@localhost/gredis?sslmode=disable"),
		maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	database, err := NewDB(
		dbc.addr,
		dbc.maxOpenConns,
		dbc.maxIdleConns,
		dbc.maxIdleTime,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()
	log.Println("database connection pool established")

}
