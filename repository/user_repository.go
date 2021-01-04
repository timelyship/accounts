package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

const UserCollection = "user"

func GetUserByGoogleID(googleID string) (*domain.User, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	verifiedEmailFilter := getVerifiedEmailFilter(emailOrPhone)
	verifiedPhoneFilter := getVerifiedPhoneFilter(emailOrPhone)
	filter := bson.D{{Key: "$or", Value: bson.A{verifiedEmailFilter, verifiedPhoneFilter}}}
	result := domain.User{}
	error := GetCollection(UserCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not query.", &error)
	}
	return &result, nil
}

func getVerifiedEmailFilter(emailOrPhone string) bson.D {
	return bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "primary_email", Value: emailOrPhone}},
		bson.D{{Key: "is_primary_email_verified", Value: true}},
	}}}
}

func getVerifiedPhoneFilter(emailOrPhone string) bson.D {
	return bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "phone_numbers.number", Value: emailOrPhone}},
		bson.D{{Key: "phone_numbers.is_verified", Value: true}},
	}}}
}
