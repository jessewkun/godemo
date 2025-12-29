package wire

import (
	repository "godemo/internal/repository"

	"github.com/google/wire"
)

// RepositorySet aggregates all repository provider sets.
var RepositorySet = wire.NewSet(
	repository.ProviderSet,
)
