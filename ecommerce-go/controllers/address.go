package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Blaze5333/ecommerce_go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for adding an address
		c.JSON(200, gin.H{"message": "Add Address endpoint"})
	}

}
func GetAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for getting an address
		c.JSON(200, gin.H{"message": "Get Address endpoint"})
	}

}
func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for editing work address
		c.JSON(200, gin.H{"message": "Edit Work Address endpoint"})
	}
}
func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for editing home address
		c.JSON(200, gin.H{"message": "Edit Home Address endpoint"})
	}
}
func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for deleting an address
		user_id := c.Query("id")
		if user_id == "" || strings.TrimSpace(user_id) == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			c.Abort()
			return
		}
		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		c.JSON(200, gin.H{"message": "Delete Address endpoint"})
		updateResult, err := UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "Address deleted successfully", "result": updateResult})
	}
}
