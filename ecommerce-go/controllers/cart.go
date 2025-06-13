package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Blaze5333/ecommerce_go/database"
	"github.com/Blaze5333/ecommerce_go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func InitializeCollections(userCollection, prodCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}

}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("productId")
		if productQueryId == "" {
			log.Println("Product ID is required")
			c.AbortWithError(http.StatusBadRequest, errors.New("product ID is required"))
			return
		}
		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("User ID is required")
			c.AbortWithError(http.StatusBadRequest, errors.New("user ID is required"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println("Invalid Product ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("invalid Product ID"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryId)
		if err != nil {
			log.Println("Error adding product to cart:", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully"})
	}
}
func (app *Application) RemoveCartItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("productId")
		if productQueryId == "" {
			log.Println("Product ID is required")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is required"))
			return
		}
		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("User ID is required")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is required"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println("Invalid Product ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("invalid product id"))
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryId)
		if err != nil {
			log.Println("Error removing product from cart:", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product removed from cart successfully"})
	}
}
func (app *Application) GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userId")
		if user_id == "" {
			log.Println("User ID is required")
			c.Header("Content-Type", "application/json")
			c.AbortWithError(http.StatusBadRequest, errors.New("user ID is required"))
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println("Invalid User ID")
			c.Header("Content-Type", "application/json")
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid User ID"))
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var filledCart models.User
		err = app.userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledCart)
		if err != nil {
			log.Println("Error fetching cart items:", err)
			c.Header("Content-Type", "application/json")
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("error fetching cart items"))
			return
		}
		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
		pipeline := mongo.Pipeline{filter_match, unwind, group}
		pointcursor, err := app.userCollection.Aggregate(ctx, pipeline)
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

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Panic("User ID is missing")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is required"))
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := database.BuyItemFromCart(ctx, app.prodCollection, app.userCollection, userQueryId)
		if err != nil {
			log.Println("Error processing cart purchase:", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Cart purchase successful"})

	}
}
func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("productId")
		if productQueryId == "" {
			log.Println("Product ID is required")
			c.AbortWithError(http.StatusBadRequest, errors.New("product ID is required"))
			return
		}
		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("User ID is required")
			c.AbortWithError(http.StatusBadRequest, errors.New("user ID is required"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println("Invalid Product ID")
			c.AbortWithError(http.StatusBadRequest, errors.New("invalid Product ID"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, userQueryId)
		if err != nil {
			log.Println("Error processing instant buy:", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Instant buy successful"})
	}
}
