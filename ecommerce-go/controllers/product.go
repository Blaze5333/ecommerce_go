package controllers

import (
	"context"
	"time"

	"github.com/Blaze5333/ecommerce_go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddProductInDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.ProductUser
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(400, gin.H{"message": "Invalid input data", "error": err.Error()})
			return
		}
		product.Product_ID = primitive.NewObjectID()
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		_, err := ProductCollection.InsertOne(ctx, product)
		if err != nil {
			c.JSON(500, gin.H{"message": "Failed to add product", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Product added successfully", "product": product})
	}
}
