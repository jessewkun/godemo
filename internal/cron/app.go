package cron

import (
	"context"

	"godemo/config"

	xcron "github.com/jessewkun/gocommon/cron"
	"github.com/jessewkun/gocommon/logger"
)

// App is the cron application.
type App struct {
	Manager *xcron.Manager
}

func NewApp(
	manager *xcron.Manager,
	cfg *config.BusinessConfig,

	// 如果有其他任务，请在这里添加
	demoTask *DemoTask,
) *App {
	taskRegistry := map[string]xcron.Task{
		// 新增一个任务时，请在这里把它加入到 map 中来把配置和具体的任务关联起来
		"demo": demoTask,
	}

	// 遍历配置文件，将配置与任务实现结合，并注册到管理器中
	for _, taskConfig := range cfg.Crons {
		taskImpl, ok := taskRegistry[taskConfig.Key]
		if !ok {
			logger.Warn(context.Background(), "CRON", "Task %s not found in registry, skipping", taskConfig.Key)
			continue
		}

		configurableTask := xcron.NewConfigurableTask(taskImpl, taskConfig)
		manager.RegisterTask(configurableTask)
	}

	return &App{Manager: manager}
}

// Start starts the cron manager.
func (a *App) Start(ctx context.Context) error {
	return a.Manager.Start(ctx)
}

// Stop stops the cron manager gracefully.
func (a *App) Stop(ctx context.Context) {
	a.Manager.Stop(ctx)
}

// RunTask manually runs a specific task by name
func (a *App) RunTask(ctx context.Context, taskName string) error {
	return a.Manager.RunTask(ctx, taskName)
}
