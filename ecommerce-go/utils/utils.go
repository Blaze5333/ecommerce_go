package utils

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func VerifyUserId(c *gin.Context) primitive.ObjectID {
	userId, exists := c.Get("user_id")
	if !exists || userId == nil || userId == "" {
		c.Header("Content-Type", "application/json")
		c.JSON(400, gin.H{"message": "User ID is missing"})
		c.Abort()
		return primitive.NilObjectID
	}
	usert_id, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.JSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return primitive.NilObjectID
	}
	return usert_id
}

func VerifyAddressId(c *gin.Context) primitive.ObjectID {
	addressId := c.Param("address_id")
	if addressId == "" || addressId == " " {
		c.Header("Content-Type", "application/json")
		c.JSON(400, gin.H{"message": "Address ID is missing"})
		c.Abort()
		return primitive.NilObjectID
	}
	address_id, err := primitive.ObjectIDFromHex(addressId)
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.JSON(500, gin.H{"error": err.Error(), "message": "Invalid Address ID"})
		c.Abort()
		return primitive.NilObjectID
	}
	return address_id
}
func VerifyProductId(c *gin.Context) primitive.ObjectID {
	productId := c.Param("product_id")
	if productId == "" || productId == " " {
		c.Header("Content-Type", "application/json")
		c.JSON(400, gin.H{"error": "Product ID is missing"})
		c.Abort()
		return primitive.NilObjectID
	}
	product_id, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.JSON(500, gin.H{"error": err.Error()})
		c.Abort()
		return primitive.NilObjectID
	}
	return product_id
}
