package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

const TOKEN_COLLECTION = "token"

func GetTokenByRefreshToken(refreshToken string) (*domain.TokenDetails, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"refresh_token", refreshToken}}
	result := domain.TokenDetails{}
	error := GetCollection(TOKEN_COLLECTION).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not find.", &error)
	}
	return &result, nil
}

func UpdateToken(token *domain.TokenDetails) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"_id", token.ID}}
	updateResult := GetCollection(TOKEN_COLLECTION).FindOneAndReplace(ctx, filter, token)
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
	insertResult, error := GetCollection(TOKEN_COLLECTION).InsertOne(ctx, token)
	fmt.Printf("%v\n", insertResult)
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return nil
}
