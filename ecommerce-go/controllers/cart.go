package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Blaze5333/ecommerce_go/database"
	"github.com/Blaze5333/ecommerce_go/models"
	"github.com/Blaze5333/ecommerce_go/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := utils.VerifyProductId(c)
		if productID == primitive.NilObjectID {
			return
		}
		userQueryId := utils.VerifyUserId(c)
		if userQueryId == primitive.NilObjectID {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := database.AddProductToCart(ctx, ProductCollection, UserCollection, productID, userQueryId)
		if err != nil {
			log.Println("Error adding product to cart:", err)
			c.Header("Content-Type", "application/json")
			// Return a JSON error response
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add product to cart", "error": err.Error()})
			// Abort the request with an error

			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully"})
	}
}
func RemoveCartItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productId := utils.VerifyProductId(c)
		if productId == primitive.NilObjectID {
			return
		}
		userQueryId := utils.VerifyUserId(c)
		if userQueryId == primitive.NilObjectID {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := database.RemoveCartItem(ctx, ProductCollection, UserCollection, productId, userQueryId)
		if err != nil {
			log.Println("Error removing product from cart:", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product removed from cart successfully"})
	}
}
func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		usert_id := utils.VerifyUserId(c)
		if usert_id == primitive.NilObjectID {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var filledCart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledCart)
		if err != nil {
			log.Println("Error fetching cart items:", err)
			c.Header("Content-Type", "application/json")
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error fetching cart items"))
			return
		}
		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$user_cart"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$user_cart.price"}}}}}}
		pipeline := mongo.Pipeline{filter_match, unwind, group}
		pointcursor, err := UserCollection.Aggregate(ctx, pipeline)
		if err != nil {
			log.Println("Error aggregating cart items:", err)
			c.Header("Content-Type", "application/json")
			c.AbortWithError(http.StatusInternalServerError, errors.New("error aggregating cart items"))
			return
		}
		defer pointcursor.Close(ctx)
		var filledCartTotal []bson.M
		if err = pointcursor.All(ctx, &filledCartTotal); err != nil {
			log.Println("Error decoding cart total:", err)
			c.Header("Content-Type", "application/json")
			c.AbortWithError(http.StatusInternalServerError, errors.New("error decoding cart total"))
			return
		}
		var results []gin.H
		for _, json := range filledCartTotal {
			results = append(results, gin.H{
				"total":    json["total"],
				"usercart": filledCart.UserCart,
			})
		}
		c.IndentedJSON(http.StatusOK, results)
	}

}

func BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := utils.VerifyUserId(c)
		if userQueryId == primitive.NilObjectID {
			return
		}
		// Check if userQueryId is empty
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := database.BuyItemFromCart(ctx, ProductCollection, UserCollection, userQueryId)
		if err != nil {
			log.Println("Error processing cart purchase:", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Cart purchase successful"})

	}
}
func InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := utils.VerifyProductId(c)
		if productID == primitive.NilObjectID {
			return
		}
		userQueryId := utils.VerifyUserId(c)
		if userQueryId == primitive.NilObjectID {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := database.InstantBuyer(ctx, ProductCollection, UserCollection, productID, userQueryId)
		if err != nil {
			log.Println("Error processing instant buy:", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Instant buy successful"})
	}
}
