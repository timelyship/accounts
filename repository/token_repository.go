package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

const TokenCollection = "token"

func GetTokenByRefreshToken(refreshToken string) (*domain.TokenDetails, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "refresh_token", Value: refreshToken}}
	result := domain.TokenDetails{}
	error := GetCollection(TokenCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not find.", &error)
	}
	return &result, nil
}

func UpdateToken(token *domain.TokenDetails) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "_id", Value: token.ID}}
	updateResult := GetCollection(TokenCollection).FindOneAndReplace(ctx, filter, token)
	error := updateResult.Err()
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not replace to database. Try after some time.", &error)
	}
	return nil
}

func SaveToken(token *domain.TokenDetails) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertResult, error := GetCollection(TokenCollection).InsertOne(ctx, token)
	fmt.Printf("%v\n", insertResult)
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return nil
}
