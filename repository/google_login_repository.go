package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

const GoogleStateCollection = "google_state"

func GetByGoogleState(state string) (*domain.GoogleState, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "state", Value: state}}
	result := domain.GoogleState{}
	error := GetCollection(GoogleStateCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return &result, nil
}

func SaveGoogleState(googleState *domain.GoogleState) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertResult, error := GetCollection(GoogleStateCollection).InsertOne(ctx, googleState)
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	fmt.Println("Successfully inserted", insertResult)
	return nil
}
