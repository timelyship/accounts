package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

const LoginSecretCollection = "login_secret"

func SaveLoginState(loginState *domain.LoginState) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertResult, error := GetCollection(LoginSecretCollection).InsertOne(ctx, loginState)
	fmt.Printf("%v\n", insertResult)
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return nil
}
func UpdateLoginState(loginState *domain.LoginState) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "state", Value: loginState.State}}
	updateResult := GetCollection(LoginSecretCollection).FindOneAndReplace(ctx, filter, loginState)
	error := updateResult.Err()
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not replace to database. Try after some time.", &error)
	}
	return nil
}
func GetLoginState(state string) (*domain.LoginState, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "state", Value: state}}
	result := domain.LoginState{}
	error := GetCollection(LoginSecretCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not query.", &error)
	}
	return &result, nil
}
