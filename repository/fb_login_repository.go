package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

const FbStateCollection = "fb_state"

type FbLoginRepository struct {
	logger zap.Logger
}

func ProvideFbLoginRepository(logger zap.Logger) FbLoginRepository {
	return FbLoginRepository{
		logger: logger,
	}
}

func GetByFBState(state string) (*domain.FBState, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "state", Value: state}}
	result := domain.FBState{}
	error := GetCollection(FbStateCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return &result, nil
}

func SaveFBState(fbState *domain.FBState) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertResult, error := GetCollection(FbStateCollection).InsertOne(ctx, fbState)
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	fmt.Println("Successfully inserted", insertResult)
	return nil
}
