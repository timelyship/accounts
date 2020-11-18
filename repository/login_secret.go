package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)
const LOGIN_SECRET_COLLECTION = "LOGIN_SECRET_COLLECTION"

func SaveLoginState(loginState *domain.LoginState) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertResult, error := GetCollection(LOGIN_SECRET_COLLECTION).InsertOne(ctx, loginState)
	fmt.Printf("%v\n", insertResult)
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return nil
}

func GetEncKeyByState(state string) (string,*utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"state", state}}
	result := domain.LoginState{}
	error := GetCollection(LOGIN_SECRET_COLLECTION).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return "", utility.NewInternalServerError("Could not query.", &error)
	}
	return result.Key, nil
}