package cron

import (
	"context"

	"github.com/jessewkun/gocommon/cron"
)

// App is the cron application.
type App struct {
	Manager *cron.Manager
}

// NewApp creates a new cron App.
// wire will inject all tasks into this constructor.
func NewApp(manager *cron.Manager, task *DemoTask) *App {
	// Register all tasks with the manager
	manager.RegisterTask(task)
	// If you have more tasks, add them here, and also in the NewApp signature.
	// manager.RegisterTask(task2)

	return &App{Manager: manager}
}

// Start starts the cron manager.
func (a *App) Start(ctx context.Context) error {
	return a.Manager.Start(ctx)
}

// Stop stops the cron manager gracefully.
func (a *App) Stop() {
	a.Manager.Stop()
}

// RunTask manually runs a specific task by name
func (a *App) RunTask(ctx context.Context, taskName string) error {
	return a.Manager.RunTask(ctx, taskName)
}
