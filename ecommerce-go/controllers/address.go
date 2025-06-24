package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Blaze5333/ecommerce_go/models"
	"github.com/Blaze5333/ecommerce_go/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		usert_id := utils.VerifyUserId(c)
		if usert_id == primitive.NilObjectID {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		match_filter := bson.D{{Key: "$match", Value: filter}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		pipeline := mongo.Pipeline{match_filter, unwind, group}
		cursor, err := UserCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		defer cursor.Close(ctx)
		var countResult []bson.M
		if err = cursor.All(ctx, &countResult); err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		var size int32
		for _, json := range countResult {
			size = json["count"].(int32)
		}
		if size < 2 {

			var address models.Address
			address.Address_ID = primitive.NewObjectID()
			if err := c.ShouldBindJSON(&address); err != nil {
				c.Header("Content-Type", "application/json")
				c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			update := bson.D{{Key: "$push", Value: bson.D{{Key: "address", Value: address}}}}
			_, err = UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				c.Header("Content-Type", "application/json")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				c.Abort()
				return
			}
			c.JSON(200, gin.H{"message": "Address added successfully", "result": address})

		} else {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, gin.H{"error": "You can only add two addresses"})
			c.Abort()
			return
		}

	}

}
func GetAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for getting an address
		usert_id := utils.VerifyUserId(c)
		if usert_id == primitive.NilObjectID {
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		var user models.User
		err := UserCollection.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Error fetching user address"})
			c.Abort()
			return
		}
		match_filter := bson.D{{Key: "$match", Value: filter}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}
		pipeline := mongo.Pipeline{match_filter, unwind, group}
		cursor, err := UserCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Internal Server Error"})
			c.Abort()
			return
		}
		defer cursor.Close(ctx)
		var countResult []bson.M
		if err = cursor.All(ctx, &countResult); err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "Error decoding address count"})
			c.Abort()
			return
		}
		var size int32
		for _, json := range countResult {
			size = json["count"].(int32)
		}

		c.JSON(http.StatusFound, gin.H{"total_address": size, "address": user.Address_Details, "message": "Fetched Address successfully"})
	}

}
func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for editing work address
		user_id := utils.VerifyUserId(c)
		if user_id == primitive.NilObjectID {
			return
		}
		address_id := utils.VerifyAddressId(c)
		if address_id == primitive.NilObjectID {
			return
		}
		var address models.Address
		address.Address_ID = address_id
		if err := c.BindJSON(&address); err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if address.Address_Line_1 == nil || address.Phone == nil || address.House == nil || address.City == nil || address.Pincode == nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: user_id}, primitive.E{Key: "address._id", Value: address_id}}
		update := bson.M{"$set": bson.M{"address.$.address_line_1": address.Address_Line_1,
			"address.$.address_line_2": address.Address_Line_2,
			"address.$.phone":          address.Phone,
			"address.$.house":          address.House,
			"address.$.city":           address.City,
			"address.$.pincode":        address.Pincode}}
		updateResult, err := UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
		if updateResult.MatchedCount == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			c.Abort()
			return
		}
		if updateResult.ModifiedCount == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotModified, gin.H{"message": "No changes made to the address"})
			c.Abort()
			return
		}

		c.JSON(200, gin.H{"message": "Edit Work Address endpoint", "result": address})
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
		usert_id := utils.VerifyUserId(c)
		if usert_id == primitive.NilObjectID {
			return
		}
		addresses := make([]models.Address, 0)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}

		updateResult, err := UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
		if updateResult.MatchedCount == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		if updateResult.ModifiedCount == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotModified, gin.H{"message": "No changes made to the address"})
			c.Abort()
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "Address deleted successfully"})
	}
}
func DeleteAddressById() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for deleting an address
		usert_id := utils.VerifyUserId(c)
		if usert_id == primitive.NilObjectID {
			return
		}
		addresst_id := utils.VerifyAddressId(c)
		if addresst_id == primitive.NilObjectID {
			return
		}
		// addresses := make([]models.Address, 0)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$pull", Value: bson.D{primitive.E{Key: "address", Value: bson.D{primitive.E{Key: "_id", Value: addresst_id}}}}}}

		updateResult, err := UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
		if updateResult.MatchedCount == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		if updateResult.ModifiedCount == 0 {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotModified, gin.H{"message": "No changes made to the address"})
			c.Abort()
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "Address deleted successfully"})
	}
}
