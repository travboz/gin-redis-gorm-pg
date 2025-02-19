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

	err := app.store.createOrUpdateProductWriteThrough(req.ID, req.Name, req.Price)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "created/updated"})
}

func (app *application) deleteProductByIdHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := app.store.deleteProductEventBased(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "deleted"})
}

func (app *application) invalidateProductInCacheHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := app.store.invalidateProductCache(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "cache invalidated"})
}

func (app *application) getProductByIdHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	product, err := app.store.getProductByIDHash(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	app.store.addToRecentProductsList(uint(id))
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
	if err := app.store.updateProductWithTransaction(uint(id), req.Name, req.Price); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "updated"})
}

func (app *application) getRecentProductsHandler(c *gin.Context) {
	products, err := app.store.getRecentProducts()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, products)
}
