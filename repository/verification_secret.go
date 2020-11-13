package repository

import (
	"context"
	"fmt"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

const VERIFICATION_SECRET_COLLECTION = "verification_secret"

func SaveVerificationSecret(vs *domain.VerificationSecret) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertResult, error := GetCollection(VERIFICATION_SECRET_COLLECTION).InsertOne(ctx, vs)
	fmt.Printf("%v\n", insertResult)
	if error != nil {
		fmt.Println("db-error:", error)
		return utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return nil
}
