package cron

import "github.com/google/wire"

// ProviderSet is a provider set for cron tasks.
var ProviderSet = wire.NewSet(
	NewDemoTask,
	// 如果有其他任务，请在这里添加
)
