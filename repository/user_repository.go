package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"timelyship.com/accounts/application"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

const UserCollection = "user"

func GetUserByGoogleID(googleID string) (*domain.User, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), application.IntConst.DBAccessMaxThreshold)
	defer cancel()
	filter := bson.D{{Key: "google_auth_info.id", Value: googleID}}
	result := domain.User{}
	error := GetCollection(UserCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return &result, nil
}

func GetUserByID(id primitive.ObjectID) (*domain.User, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), application.IntConst.DBAccessMaxThreshold)
	defer cancel()
	filter := bson.D{{Key: "_id", Value: id}}
	result := domain.User{}
	error := GetCollection(UserCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not query.", &error)
	}
	return &result, nil
}

func GetUserByEmailOrPhone(emailOrPhone string) (*domain.User, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), application.IntConst.DBAccessMaxThreshold)
	defer cancel()
	verifiedEmailFilter := getVerifiedEmailFilter(emailOrPhone)
	verifiedPhoneFilter := getVerifiedPhoneFilter(emailOrPhone)
	emailOrPhoneFilter := bson.D{{Key: "$or", Value: bson.A{verifiedEmailFilter, verifiedPhoneFilter}}}
	activeFilter := getActiveUserFilter()
	filter := bson.D{{Key: "$and", Value: bson.A{emailOrPhoneFilter, activeFilter}}}
	result := domain.User{}
	error := GetCollection(UserCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not query.", &error)
	}
	return &result, nil
}

func getActiveUserFilter() bson.D {
	return bson.D{{Key: "is_active", Value: true}}
}

func getVerifiedEmailFilter(emailOrPhone string) bson.D {
	return bson.D{{Key: "email", Value: emailOrPhone}}
}

func getVerifiedPhoneFilter(emailOrPhone string) bson.D {
	return bson.D{{Key: "phone", Value: emailOrPhone}}
}
