package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"time"
	"timelyship.com/accounts/dto/request"
	"timelyship.com/accounts/utility"
)

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

func ProvideProfileRepository(l zap.Logger) ProfileRepository {
	return ProfileRepository{
		logger: l,
	}
}
