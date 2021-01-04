package repository

import (
	"context"
	"fmt"
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

//func
