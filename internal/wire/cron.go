//go:build wireinject
// +build wireinject

package wire

import (
	"godemo/config"
	"godemo/internal/cron"
	"godemo/internal/wire/provider"

	repository "godemo/internal/repository"

	xcron "github.com/jessewkun/gocommon/cron"

	"github.com/google/wire"
)

// providerBusinessConfig a provider for the business config.
func providerBusinessConfig() *config.BusinessConfig {
	return config.BusinessCfg
}

// InitializeCron initializes the cron application.
func InitializeCron() (*cron.App, func(), error) {
	wire.Build(
		// 基础服务提供者
		providerBusinessConfig,

		// Infrastructure providers
		provider.ProvideMainDB,
		wire.Value(provider.MainDBNameValue),

		// Cron 框架和所有任务的构造函数
		xcron.NewManager,
		cron.ProviderSet,

		repository.ProviderSet,
		// 最终的应用组装者
		cron.NewApp,
	)
	return &cron.App{}, nil, nil
}
