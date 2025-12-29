package wire

import (
	"godemo/internal/handler"

	"github.com/google/wire"
)

// HandlerSet aggregates all handler provider sets.
var HandlerSet = wire.NewSet(
	handler.ProviderSet,
)
