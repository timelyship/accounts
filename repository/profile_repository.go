package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"time"
	"timelyship.com/accounts/domain"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/utility"
)

var PhoneVerificationQueue = "phone_verification_queue"

type ProfileRepository struct {
	logger zap.Logger
}

func (r *ProfileRepository) Patch(id string, request []*request.ProfilePatchRequest) *utility.RestError {
	userID, parseHexErr := primitive.ObjectIDFromHex(id)
	if parseHexErr != nil {
		r.logger.Error("User id parse error", zap.Error(parseHexErr))
		return utility.NewInternalServerError("Could not parse userId", &parseHexErr)
	}
	patches := bson.A{}
	for _, req := range request {
		patches = append(patches, bson.D{
			{Key: "$set", Value: bson.D{{Key: req.Field, Value: req.Value}}},
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := GetCollection(UserCollection).UpdateOne(ctx, bson.M{"_id": userID}, patches)
	if err != nil {
		r.logger.Error("Update User", zap.Error(err))
		return utility.NewInternalServerError("Could not replace to database. Try after some time.", &err)
	}
	return nil
}

func (r *ProfileRepository) GetProfileById(id primitive.ObjectID) (*domain.User, *utility.RestError) {
	return GetUserByID(id)
}

func (r *ProfileRepository) ChangePhoneNumber(userID primitive.ObjectID, phone string) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	update := bson.M{
		"$set": bson.M{
			"phone":             phone,
			"is_phone_verified": false,
		}}
	updateResult, err := GetCollection(UserCollection).UpdateOne(ctx, bson.M{"_id": userID}, update)
	if updateResult.MatchedCount == 0 {
		rErrMsg := fmt.Sprintf("Match not found with key = userId,value=%s, %v", userID, updateResult)
		rErr := errors.New(rErrMsg)
		return utility.NewBadRequestError(rErrMsg, &rErr)
	}
	r.logger.Debug("updateResult", zap.Any("updateResult", updateResult))
	if err != nil {
		r.logger.Error("Update User", zap.Error(err))
		return utility.NewInternalServerError("Could not replace to database. Try after some time.", &err)
	}
	return nil
}

func (r *ProfileRepository) EnqueuePhoneVerification(userID primitive.ObjectID, verification *domain.PhoneVerification) *utility.RestError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	insertResult, err := GetCollection(PhoneVerificationQueue).InsertOne(ctx, verification)
	r.logger.Debug("insertResult", zap.Any("insertResult", insertResult))
	if err != nil {
		utility.NewInternalServerError("Queue phone verification failed", &err)
	}
	return nil
}

func ProvideProfileRepository(l zap.Logger) ProfileRepository {
	return ProfileRepository{
		logger: l,
	}
}
