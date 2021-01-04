package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AccountRepository struct {
	collection mongo.Collection
	logger     zap.Logger
}

func ProvideAccountRepository(mongoCollection mongo.Collection, logger zap.Logger) AccountRepository {
	return AccountRepository{
		collection: mongoCollection,
		logger:     logger,
	}
}
