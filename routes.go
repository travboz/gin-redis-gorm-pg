package main

import "github.com/gin-gonic/gin"

func setupRouter(app *application) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")
	{
		products := v1.Group("/products")
		{
			products.POST("/", app.createProductHandler)

			products.GET("/:id", app.getProductByIdHandler)
			products.PUT("/:id", app.updateProductByIdHandler)
			products.DELETE("/:id", app.deleteProductByIdHandler)

			products.POST("/invalidate/:id", app.invalidateProductInCacheHandler)
			products.GET("/recent", app.getRecentProductsHandler)
		}
	}
	return router
}
