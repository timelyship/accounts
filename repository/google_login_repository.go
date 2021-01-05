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

const GoogleStateCollection = "google_state"

type GoogleLoginRepository struct {
	logger zap.Logger
}

func ProvideGoogleLoginRepository(logger zap.Logger) GoogleLoginRepository {
	return GoogleLoginRepository{
		logger: logger,
	}
}

func (r GoogleLoginRepository) SaveGoogleState(googleState *domain.GoogleState) *utility.RestError {
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

func (r GoogleLoginRepository) GetByGoogleState(state string) (*domain.GoogleState, *utility.RestError) {
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
