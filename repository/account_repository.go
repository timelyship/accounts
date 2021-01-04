package repository

import (
	"go.uber.org/zap"
)

type AccountRepository struct {
	logger zap.Logger
}

func ProvideAccountRepository(logger zap.Logger) AccountRepository {
	return AccountRepository{
		logger: logger,
	}
}
