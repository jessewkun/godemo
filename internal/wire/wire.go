//go:build wireinject
// +build wireinject

package wire

import (
	"godemo/internal/handler"
	"godemo/internal/repository"
	"godemo/internal/service"
	"godemo/internal/wire/provider"

	"github.com/google/wire"
)

// InitializeAPI 完整依赖注入
func InitializeAPI() (*handler.UserHandler, error) {
	wire.Build(
		// DB
		wire.Value(provider.MainDBNameValue),
		provider.ProvideMainDB,
		// Cache
		wire.Value(provider.MainCacheNameValue),
		provider.ProvideMainCache,
		// Repository
		repository.NewUserRepository,
		// Service
		service.NewUserService,
		// Handler
		handler.NewUserHandler,
	)
	return nil, nil
}
