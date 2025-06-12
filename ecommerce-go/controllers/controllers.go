package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Blaze5333/ecommerce_go/database"
	"github.com/Blaze5333/ecommerce_go/models"
	"github.com/Blaze5333/ecommerce_go/tokens"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// hash password,verify password
func HashPassword(password string) string {
	// Implementation for hashing the password
	return password // Placeholder, replace with actual hashing logic
}
func VerifyPassword(userpassword, givenpassword string) (bool, string) {
	// Implementation for verifying the password
	return userpassword == givenpassword, "Password verification failed" // Placeholder, replace with actual verification logic
}

// signup
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for user signup
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		UserCollection := database.UserData(database.Client, "Users")
		var count, err = UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user existence"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		var count1, err1 = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking phone existence"})
			return
		}
		if count1 > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone number already exists"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID() // Example user ID generation
		user.User_ID = user.ID.Hex()
		token, refresh_token, _ := tokens.GenerateToken(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refresh_token
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)

		_, err = UserCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user into database"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User signed up successfully"})
	}
}

// login
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for user login
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		UserCollection := database.UserData(database.Client, "Users")

		var foundUser models.User
		err := UserCollection.FindOne(ctx, bson.M{"email": &user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password is incorrect"})
			return
		}
		passwordCorrect, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if !passwordCorrect {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			fmt.Println("Password verification failed:", msg)
			return
		}
		token, refresh_token, _ := tokens.GenerateToken(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_ID)
		err1 := tokens.UpdateToken(token, refresh_token, foundUser.User_ID, UserCollection)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating token"})
			return
		}
		foundUser.Token = &token
		foundUser.Refresh_Token = &refresh_token
		foundUser.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		c.JSON(http.StatusFound, foundUser)
	}
}

//product view admin

func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for adding a product by admin
		c.JSON(200, gin.H{"message": "Product added successfully"})
	}
}

//search product

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for searching a product
		c.JSON(200, gin.H{"message": "Product search results"})
	}
}

//search product by query

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for searching a product by query
		c.JSON(200, gin.H{"message": "Product search results by query"})
	}
}
