package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func DBSet() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://mustafa:Mufaddal*53@cluster0.mqghj5s.mongodb.net/"))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		return nil
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		return nil
	}
	log.Println("Connected to MongoDB successfully")
	return client
}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommerce_go").Collection(collectionName)
	return collection
}
func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommerce_go").Collection(collectionName)
	return collection
}
