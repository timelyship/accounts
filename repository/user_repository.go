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

const USER_COLLECTION = "user"

func GetUserByGoogleId(googleId string) (*domain.User, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"google_auth_info.id", googleId}}
	result := domain.User{}
	error := GetCollection(USER_COLLECTION).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return &result, nil
}

func GetUserByEmail(email string) (*domain.User, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"$or", bson.A{
		bson.D{{"primary_email", email}},
		bson.D{{"google_auth_info.email", email}},
		bson.D{{"facebook_auth_info.email", email}},
	}}}
	result := domain.User{}
	error := GetCollection(USER_COLLECTION).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not query.", &error)
	}
	return &result, nil
}

func GetUserById(id primitive.ObjectID) (*domain.User, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{"_id", id}}
	result := domain.User{}
	error := GetCollection(USER_COLLECTION).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not query.", &error)
	}
	return &result, nil
}

func GetUserByEmailOrPhone(email, phone string) (*domain.User, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	verifiedEmailFilter := getVerifiedEmailFilter(email)
	verifiedPhoneFilter := getVerifiedPhoneFilter(phone)
	filter := bson.D{{"$or", bson.A{verifiedEmailFilter, verifiedPhoneFilter}}}
	result := domain.User{}
	error := GetCollection(USER_COLLECTION).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewInternalServerError("Could not query.", &error)
	}
	return &result, nil
}

func getVerifiedEmailFilter(email string) bson.D {
	return bson.D{{"$and", bson.A{
		bson.D{{"primary_email", email}},
		bson.D{{"is_primary_email_verified", true}},
	}}}
}

func getVerifiedPhoneFilter(phone string) bson.D {
	return bson.D{{"$and", bson.A{
		bson.D{{"phone_numbers.number", phone}},
		bson.D{{"phone_numbers.is_verified", true}},
	}}}
}

func SaveUser(user *domain.User) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertResult, error := GetCollection(USER_COLLECTION).InsertOne(ctx, user)
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
	filter := bson.D{{"_id", user.Id}}
	updateResult := GetCollection(USER_COLLECTION).FindOneAndReplace(ctx, filter, user)
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
	count, error := GetCollection(USER_COLLECTION).CountDocuments(ctx, filter)
	if error != nil {
		fmt.Println("db-error:", error)
		return false, utility.NewInternalServerError("Could not query.", &error)
	}
	return count > 0, nil
}
