package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

const VerificationSecretCollection = "verification_secret"

func SaveVerificationSecret(vs *domain.VerificationSecret) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertResult, error := GetCollection(VerificationSecretCollection).InsertOne(ctx, vs)
	fmt.Printf("%v\n", insertResult)
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return nil
}

func GetVerificationSecret(secret string) (*domain.VerificationSecret, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "secret", Value: secret}},
		bson.D{{Key: "valid_until", Value: bson.D{
			{Key: "$gte", Value: time.Now()},
		}}},
	}}}
	js, _ := json.Marshal(filter)
	fmt.Printf("mgo query: %v\n", string(js))
	result := domain.VerificationSecret{}
	error := GetCollection(VerificationSecretCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		fmt.Println("db-error:", error)
		return nil, utility.NewBadRequestError("Invalid secret", &error)
	}
	return &result, nil
}
