//+build wireinject

package appwiring

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"timelyship.com/accounts/repository"
	"timelyship.com/accounts/service"
)

// // a repository.AccountRepository, z zap.Logger
func InitAccountService(logger zap.Logger) service.AccountService {
	wire.Build(repository.ProvideAccountRepository, service.ProvideAccountService)
	return service.AccountService{}
}

func InitAuthService(logger zap.Logger) service.AuthService {
	wire.Build(repository.ProvideAuthRepository, service.ProvideAuthService)
	return service.AuthService{}
}

func InitFbLoginService(logger zap.Logger) service.FbLoginService {
	wire.Build(repository.ProvideFbLoginRepository, service.ProvideFbLoginService)
	return service.FbLoginService{}
}

func InitGoogleLoginService(logger zap.Logger) service.GoogleLoginService {
	wire.Build(repository.ProvideGoogleLoginRepository, repository.ProvideAccountRepository, service.ProvideGoogleLoginService)
	return service.GoogleLoginService{}
}

func InitProfileService(logger zap.Logger) service.ProfileService {
	wire.Build(repository.ProvideProfileRepository, service.ProvideProfileService)
	return service.ProfileService{}
}
