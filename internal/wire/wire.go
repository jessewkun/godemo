//go:build wireinject
// +build wireinject

package wire

import (
	"godemo/internal/handler"
	"godemo/internal/wire/provider"

	"github.com/google/wire"
)

type APIs struct {
	UserHandler *handler.UserHandler
}

func InitializeAPIs() (*APIs, error) {
	panic(wire.Build(
		// Infrastructure providers
		provider.ProvideMainDB,
		wire.Value(provider.MainDBNameValue),
		provider.ProvideMainCache,
		wire.Value(provider.MainCacheNameValue),

		// Aggregated provider sets
		RepositorySet,
		ServiceSet,
		HandlerSet,

		// The final struct to build
		wire.Struct(new(APIs), "*"),
	))
}
