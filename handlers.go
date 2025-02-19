package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (app *application) createProductHandler(c *gin.Context) {
	var req struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Price int    `json:"price"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	err := app.cache.createOrUpdateProductWriteThrough(req.ID, req.Name, req.Price)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "created/updated"})
}

func (app *application) deleteProductByIdHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := app.cache.deleteProductEventBased(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "deleted"})
}

func (app *application) invalidateProductInCacheHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := app.cache.invalidateProductCache(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "cache invalidated"})
}

func (app *application) getProductByIdHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := app.cache.getProductByIDHash(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	app.cache.addToRecentProductsList(uint(id))
	c.JSON(200, product)
}

func (app *application) updateProductByIdHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	if err := app.cache.updateProductWithTransaction(uint(id), req.Name, req.Price); err != nil {
		switch err {
		case ErrProductNotFound:
			c.JSON(400, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(200, gin.H{"status": "updated"})
}

func (app *application) getRecentProductsHandler(c *gin.Context) {
	products, err := app.cache.getRecentProducts()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, products)
}
