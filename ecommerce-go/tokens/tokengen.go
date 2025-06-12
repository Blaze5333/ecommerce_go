package tokens

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateToken(email, firstName, lastName, userID string) (string, string, error) {
	// Implementation for generating a token
	// This is a placeholder; replace with actual token generation logic
	return "generated_token_for_" + userID, "generated_refresh_token_for_" + userID, nil
}
func UpdateToken(token, refreshToken, userID string, userCollection *mongo.Collection) error {
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
	_, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
