package main

/*
https://dev.to/truongpx396/golang-restful-api-with-gin-gorm-redis-cache-2gia
*/

import (
	_ "github.com/lib/pq"
)

// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// }

// func main() {
// 	dbc := dbConfig{
// 		addr:         env.GetString("DB_ADDR", "postgres://admin:adminpass@localhost/gredis?sslmode=disable"),
// 		maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
// 		maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
// 		maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
// 	}

// 	database, err := NewDB(
// 		dbc.addr,
// 		dbc.maxOpenConns,
// 		dbc.maxIdleConns,
// 		dbc.maxIdleTime,
// 	)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer database.Close()
// 	log.Println("database connection pool established")

// }

func main() {
	db := initDB()

	app := application{
		store: NewGormPgStorage(db),
	}

	// used for write-behind caching - not the most consistent method
	// go writeBehindWorker(db) // start the background worker who listens to the queue

	router := setupRouter(&app)
	router.Run(":8080")
}
