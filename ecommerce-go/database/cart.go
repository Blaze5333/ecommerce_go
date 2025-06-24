package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Blaze5333/ecommerce_go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrorProductNotFound     = errors.New("product not found")
	ErrorUserIdIsNotValid    = errors.New("user ID is not valid")
	ErrorCantDecodeProducts  = errors.New("can't decode products")
	ErrorCantUpdateCart      = errors.New("can't update cart")
	ErrorCantRemoveCartItem  = errors.New("can't remove cart item")
	ErrorCantGetCartItems    = errors.New("can't get cart items")
	ErrorCantBuyItemFromCart = errors.New("can't buy item from cart")
	ErrorCantInstantBuyItem  = errors.New("can't instant buy item")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID, userId primitive.ObjectID) error {
	searchFromDb := prodCollection.FindOne(ctx, bson.M{"_id": productID})

	if searchFromDb == nil {
		log.Println("Product not found in database")
		return ErrorProductNotFound
	}
	var product models.ProductUser
	err := searchFromDb.Decode(&product)

	if err != nil {
		log.Println("Error decoding product:", err.Error())
		return ErrorCantDecodeProducts
	}
	var filter = bson.D{{Key: "_id", Value: userId}}
	var update = bson.D{{Key: "$push", Value: bson.D{{Key: "user_cart", Value: product}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error updating user cart:", err)
		return ErrorCantUpdateCart
	}

	return nil // Placeholder, replace with actual logic

}
func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID, userId primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "user_cart", Value: bson.D{{Key: "_id", Value: productID}}}}}}
	_, err := userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println("Error removing item from cart:", err)
		return ErrorCantRemoveCartItem
	}
	return nil // Placeholder, replace with actual logic

}
func BuyItemFromCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, userId primitive.ObjectID) error {

	defer ctx.Done()
	filter := bson.D{{Key: "_id", Value: userId}}
	var user models.User
	err := userCollection.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		log.Println("Error fetching user cart:", err)
		return ErrorCantGetCartItems
	}
	filter_match := bson.D{{Key: "$match", Value: filter}}
	unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$user_cart"}}}}
	group := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{{Key: "$sum", Value: "$user_cart.price"}}}}}}
	pipeline := mongo.Pipeline{filter_match, unwind, group}
	cursor, err := userCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Error aggregating user cart:", err)
		return ErrorCantGetCartItems
	}
	defer cursor.Close(ctx)
	var filledCartTotal []bson.M
	if err = cursor.All(ctx, &filledCartTotal); err != nil {
		log.Println("Error decoding cart total:", err)
		return ErrorCantGetCartItems
	}
	if len(filledCartTotal) == 0 {
		log.Println("Cart is empty")
		return ErrorCantGetCartItems
	}
	var totalCost int64
	for _, item := range filledCartTotal {
		totalCost = item["total"].(int64)
	}
	var u = uint64(totalCost)

	if totalCost <= 0 {
		log.Println("Cart is empty")
		return ErrorCantGetCartItems
	}
	// Proceed with the buying logic
	var Order models.Order
	Order.UserId = userId
	Order.Order_ID = primitive.NewObjectID()

	Order.Order_Cart = user.UserCart

	Order.Price = &u
	Order.Payment_Method.COD = true
	Order.Ordered_At = time.Now()
	Order.Discount = nil // Assuming no discount for simplicity
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "order_status", Value: Order}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error updating user order status:", err)
		return ErrorCantBuyItemFromCart
	}

	// Clear the user's cart after successful purchase
	clearCartUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: "user_cart", Value: []models.ProductUser{}}}}}
	_, err = userCollection.UpdateOne(ctx, filter, clearCartUpdate)
	if err != nil {
		log.Println("Error clearing user cart after purchase:", err)
		return ErrorCantBuyItemFromCart
	}

	return nil // Placeholder, replace with actual logic

}
func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userId primitive.ObjectID) error {

	searchFromDb := prodCollection.FindOne(ctx, bson.M{"_id": productID})
	if searchFromDb == nil {
		log.Println("Product not found in database")
		return ErrorProductNotFound
	}
	var product models.ProductUser

	err := searchFromDb.Decode(&product)
	if err != nil {
		log.Println("Error decoding product:", err)
		return ErrorCantDecodeProducts
	}
	var Order models.Order
	Order.UserId = userId
	Order.Order_ID = primitive.NewObjectID()
	Order.Order_Cart = []models.ProductUser{product}
	var totalCost uint64 = 0
	totalCost += product.Price
	Order.Price = &totalCost
	Order.Payment_Method.COD = true
	Order.Ordered_At = time.Now()
	Order.Discount = nil // Assuming no discount for simplicity
	filter := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "order_status", Value: Order}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error updating user order status:", err)
		return ErrorCantBuyItemFromCart
	}
	return nil // Placeholder, replace with actual logic

}
