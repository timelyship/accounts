package repository

import "go.uber.org/zap"

type AuthRepository struct {
	logger zap.Logger
}

func ProvideAuthRepository(l zap.Logger) AuthRepository {
	return AuthRepository{
		logger: l,
	}
}
