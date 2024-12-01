package main

import (
	"context"
	"log"

	// "time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var OrderCollection *mongo.Collection

func InitOrderDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB for Order Service: %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB for Order Service: %v", err)
	}

	OrderCollection = client.Database("microservices_orders").Collection("orders")
	log.Println("Connected to MongoDB: microservices_orders database")
}
