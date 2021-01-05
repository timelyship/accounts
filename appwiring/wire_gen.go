// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package appwiring

import (
	"go.uber.org/zap"
	"timelyship.com/accounts/repository"
	"timelyship.com/accounts/service"
)

// Injectors from wire.go:

func InitAccountService(logger zap.Logger) service.AccountService {
	accountRepository := repository.ProvideAccountRepository(logger)
	accountService := service.ProvideAccountService(accountRepository, logger)
	return accountService
}

func InitAuthService(logger zap.Logger) service.AuthService {
	authRepository := repository.ProvideAuthRepository(logger)
	authService := service.ProvideAuthService(authRepository, logger)
	return authService
}

func InitFbLoginService(logger zap.Logger) service.FbLoginService {
	fbLoginRepository := repository.ProvideFbLoginRepository(logger)
	fbLoginService := service.ProvideFbLoginService(fbLoginRepository, logger)
	return fbLoginService
}

func InitGoogleLoginService(logger zap.Logger) service.GoogleLoginService {
	accountRepository := repository.ProvideAccountRepository(logger)
	googleLoginRepository := repository.ProvideGoogleLoginRepository(logger)
	googleLoginService := service.ProvideGoogleLoginService(accountRepository, googleLoginRepository, logger)
	return googleLoginService
}
