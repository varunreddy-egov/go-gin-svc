package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// GET with path param
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"message": "User info",
			"id":      id,
		})
	})

	// GET with query param
	r.GET("/search", func(c *gin.Context) {
		query := c.Query("q")
		c.JSON(http.StatusOK, gin.H{
			"message": "Search results",
			"query":   query,
		})
	})

	// GET with both path and query params
	r.GET("/order/:orderId/item", func(c *gin.Context) {
		orderId := c.Param("orderId")
		itemId := c.Query("itemId")
		c.JSON(http.StatusOK, gin.H{
			"message": "Order item details",
			"orderId": orderId,
			"itemId":  itemId,
		})
	})

	// GET with multiple query params
	r.GET("/filter", func(c *gin.Context) {
		category := c.Query("category")
		price := c.Query("price")
		c.JSON(http.StatusOK, gin.H{
			"message":  "Filter results",
			"category": category,
			"price":    price,
		})
	})

	// GET with multiple path params
	r.GET("/company/:companyId/employee/:employeeId", func(c *gin.Context) {
		companyId := c.Param("companyId")
		employeeId := c.Param("employeeId")
		c.JSON(http.StatusOK, gin.H{
			"message":    "Employee info",
			"companyId":  companyId,
			"employeeId": employeeId,
		})
	})

	r.Run(":8081")
}
