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
