package tokens

import (
	"context"
	"log"
	"os"

	// "os"
	"time"

	"github.com/Blaze5333/ecommerce_go/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}

var UserData *mongo.Collection = database.UserData(database.Client, "Users")

var SECRET_KEY string

func init() {
	// Load the secret key from environment variable or use a default value
	// SECRET_KEY = os.Getenv("SECRET_KEY")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	SECRET_KEY = os.Getenv("SECRET_KEY")
}

func GenerateToken(email, firstName, lastName, userID string) (string, string, error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Uid:       userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24).Unix(), // Token expires in 72 hours
			Issuer:    "EcommerceApp",
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 168).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println("Error signing refresh token:", err)
		return "", "", err
	}
	return token, refreshToken, nil
}
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		msg = err.Error()
		return nil, msg
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return nil, msg
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "the token is expired"
		return nil, msg
	} // Extend the expiration time
	return claims, msg
}
func UpdateToken(token, refreshToken, userID string) error {
	// Implementation for updating the token
	// This is a placeholder; replace with actual token update logic
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"user_id": userID}

	update := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"updated_at":    time.Now(),
		},
	}
	_, err := UserData.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
