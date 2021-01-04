package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var (
	MongoClient = GetClient()
)

func GetClient() *mongo.Client {
	uri := os.Getenv("MONGO_CONNECTION_STRING")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		// log error
		panic("Could not connect to database.Try after some time.")
	}
	return client
}

func GetDataBase() *mongo.Database {
	return GetClient().Database(os.Getenv("DATABASE"))
}

func GetCollection(collection string) *mongo.Collection {
	return GetDataBase().Collection(collection)
}
