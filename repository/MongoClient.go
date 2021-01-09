package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var mongoClient *mongo.Client

func GetClient() *mongo.Client {
	return mongoClient
}
func DisconnectMongoClient() {
	if mongoClient != nil {
		log.Println("Disconnecting mongo connection")
	} else {
		log.Println("Can not disconnect mongoClient is null")
	}
}

func GetDataBase() *mongo.Database {
	return GetClient().Database(os.Getenv("DATABASE"))
}

func GetCollection(collection string) *mongo.Collection {
	return GetDataBase().Collection(collection)
}

func InitClient() {
	var err error
	uri := os.Getenv("MONGO_CONNECTION_STRING")
	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Could not connect to database.Try after some time.")
	}
}
