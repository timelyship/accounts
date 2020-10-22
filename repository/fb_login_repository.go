package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

const FB_STATE_COLLECTION = "fb_state"

func GetByFBState(state string) (*domain.FBState, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"state", state}}
	result := domain.FBState{}
	error := GetCollection(FB_STATE_COLLECTION).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not insert to database. Try after some time.")
	}
	return &result, nil
}

func SaveFBState(fbState *domain.FBState) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertResult, error := GetCollection(FB_STATE_COLLECTION).InsertOne(ctx, fbState)
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not insert to database. Try after some time.")
	}
	fmt.Println("Successfully inserted", insertResult)
	return nil
}
