package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

const UserCollection = "user"

func GetUserByGoogleID(googleId string) (*domain.User, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "google_auth_info.id", Value: googleId}}
	result := domain.User{}
	error := GetCollection(UserCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return &result, nil
}

func GetUserByEmail(email string) (*domain.User, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "primary_email", Value: email}},
		bson.D{{Key: "google_auth_info.email", Value: email}},
		bson.D{{Key: "facebook_auth_info.email", Value: email}},
	}}}
	result := domain.User{}
	error := GetCollection(UserCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not query.", &error)
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
		bson.D{{"primary_email", emailOrPhone}},
		bson.D{{"is_primary_email_verified", true}},
	}}}
}

func getVerifiedPhoneFilter(emailOrPhone string) bson.D {
	return bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "phone_numbers.number", Value: emailOrPhone}},
		bson.D{{Key: "phone_numbers.is_verified", Value: true}},
	}}}
}

func SaveUser(user *domain.User) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertResult, error := GetCollection(UserCollection).InsertOne(ctx, user)
	fmt.Printf("%v\n", insertResult)
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return nil
}

func UpdateUser(user *domain.User) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"_id", user.ID}}
	updateResult := GetCollection(UserCollection).FindOneAndReplace(ctx, filter, user)
	error := updateResult.Err()
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not replace to database. Try after some time.", &error)
	}
	return nil
}

func IsExistingEmail(email string) (bool, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// https://stackoverflow.com/questions/51179588/how-to-sort-and-limit-results-in-mongodb/51181206
	options := options.Find()
	options.SetLimit(1)
	filter := bson.D{{"$or", bson.A{
		bson.D{{"primary_email", email}},
		bson.D{{"google_auth_info.email", email}},
		bson.D{{"facebook_auth_info.email", email}},
	}}}
	count, error := GetCollection(UserCollection).CountDocuments(ctx, filter)
	if error != nil {
		fmt.Println("db-error:", error)
		return false, utility.NewInternalServerError("Could not query.", &error)
	}
	return count > 0, nil
}
