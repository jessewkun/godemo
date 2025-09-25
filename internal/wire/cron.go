//go:build wireinject
// +build wireinject

package wire

import (
	"godemo/internal/cron"
	"godemo/internal/repository"
	"godemo/internal/wire/provider"

	"github.com/google/wire"
	commonCron "github.com/jessewkun/gocommon/cron"
)

// InitializeCronApp initializes the cron application.
func InitializeCronApp() (*cron.App, func(), error) {
	panic(wire.Build(
		// DB
		wire.Value(provider.MainDBNameValue),
		provider.ProvideMainDB,

		// Repository
		repository.NewUserRepository,

		// Cron Manager
		commonCron.NewManager,

		// Tasks
		cron.NewDemoTask,

		// App
		cron.NewApp,
	))
}
