package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/utility"
)

type AccountRepository struct {
	logger zap.Logger
}

func (r *AccountRepository) IsExistingEmail(email string) (bool, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// https://stackoverflow.com/questions/51179588/how-to-sort-and-limit-results-in-mongodb/51181206
	options := options.Find()
	options.SetLimit(1)
	filter := bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "email", Value: email}},
		// todo : fix this when you implement google auth flow
		bson.D{{Key: "google_auth_info.email", Value: email}},
		bson.D{{Key: "facebook_auth_info.email", Value: email}},
	}}}
	count, error := GetCollection(UserCollection).CountDocuments(ctx, filter)
	if error != nil {
		r.logger.Info("IsExistingEmail db lookup error", zap.Error(error))
		return false, utility.NewInternalServerError("Could not query.", &error)
	}
	return count > 0, nil
}

func (r *AccountRepository) SaveUser(user *domain.User) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, error := GetCollection(UserCollection).InsertOne(ctx, user)
	if error != nil {
		fmt.Println("db-error:", error)
		r.logger.Error("Save error", zap.Error(error))
		return utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return nil
}

func (r *AccountRepository) SaveVerificationSecret(vs *domain.VerificationSecret) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, error := GetCollection(VerificationSecretCollection).InsertOne(ctx, vs)
	if error != nil {
		r.logger.Error("Save Verification Secret", zap.Error(error))
		return utility.NewInternalServerError("Could not insert to database. Try after some time.", &error)
	}
	return nil
}

func (r *AccountRepository) GetVerificationSecret(secret string) (*domain.VerificationSecret, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "$and", Value: bson.A{
		bson.D{{Key: "secret", Value: secret}},
		bson.D{{Key: "valid_until", Value: bson.D{
			{Key: "$gte", Value: time.Now()},
		}}},
	}}}
	//js, _ := json.Marshal(filter)
	//fmt.Printf("mgo query: %v\n", string(js))
	result := domain.VerificationSecret{}
	error := GetCollection(VerificationSecretCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		r.logger.Error("Get Verification Secret error", zap.Error(error))
		return nil, utility.NewBadRequestError("Invalid secret", &error)
	}
	return &result, nil
}

func (r *AccountRepository) GetUserByEmail(email string) (*domain.User, *utility.RestError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "email", Value: email}}
	result := domain.User{}
	error := GetCollection(UserCollection).FindOne(ctx, filter).Decode(&result)
	if error != nil {
		r.logger.Error("Get User By Email", zap.Error(error))
		return nil, utility.NewInternalServerError("Could not query.", &error)
	}
	return &result, nil
}

func (r *AccountRepository) UpdateUser(user *domain.User) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{Key: "_id", Value: user.ID}}
	updateResult := GetCollection(UserCollection).FindOneAndReplace(ctx, filter, user)
	error := updateResult.Err()
	if error != nil {
		r.logger.Error("Update User", zap.Error(error))
		return utility.NewInternalServerError("Could not replace to database. Try after some time.", &error)
	}
	return nil
}

func ProvideAccountRepository(logger zap.Logger) AccountRepository {
	return AccountRepository{
		logger: logger,
	}
}
