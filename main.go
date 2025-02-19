package main

/*
https://dev.to/truongpx396/golang-restful-api-with-gin-gorm-redis-cache-2gia
*/

import (
	"strconv"

	"github.com/gin-gonic/gin"
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

	// used for write-behind caching - not the most consistent method
	// go writeBehindWorker(db) // start the background worker who listens to the queue

	router := gin.Default()

	router.POST("/product", func(c *gin.Context) {
		var req struct {
			ID    uint   `json:"id"`
			Name  string `json:"name"`
			Price int    `json:"price"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		err := createOrUpdateProductWriteThrough(db, req.ID, req.Name, req.Price)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "created/updated"})
	})

	router.DELETE("/product/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		err := deleteProductEventBased(db, uint(id))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "deleted"})
	})

	router.POST("/product/invalidate/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		err := invalidateProductCache(uint(id))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "cache invalidated"})
	})

	router.GET("/product/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		product, err := getProductByIDHash(db, uint(id))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		addToRecentProductsList(uint(id))
		c.JSON(200, product)
	})

	router.PUT("/product/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var req struct {
			Name  string `json:"name"`
			Price int    `json:"price"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		if err := updateProductWithTransaction(db, uint(id), req.Name, req.Price); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "updated"})
	})

	router.GET("/recent_products", func(c *gin.Context) {
		products, err := getRecentProducts(db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, products)
	})

	router.Run(":8080")
}
