package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Blaze5333/ecommerce_go/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func initializeCollections(userCollection, prodCollection *mongo.Collection) *Application {
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
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product ID is required"))
			return
		}
		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("User ID is required")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is required"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println("Invalid Product ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid Product ID"))
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
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product ID is required"))
			return
		}
		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("User ID is required")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is required"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println("Invalid Product ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid Product ID"))
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
func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {

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
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product ID is required"))
			return
		}
		userQueryId := c.Query("userId")
		if userQueryId == "" {
			log.Println("User ID is required")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User ID is required"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println("Invalid Product ID")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Invalid Product ID"))
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
