package wire

import (
	service "godemo/internal/service"

	"github.com/google/wire"
)

// ServiceSet aggregates all service provider sets.
var ServiceSet = wire.NewSet(
	service.ProviderSet,
)
