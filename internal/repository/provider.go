package repository

import "github.com/google/wire"

// ProviderSet is a Wire provider set that provides all repositories for the tiku module.
var ProviderSet = wire.NewSet(
	NewUserRepository,
)
