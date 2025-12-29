// Package cron 提供定时任务的具体实现
package cron

import (
	"context"

	repo "godemo/internal/repository"

	xcron "github.com/jessewkun/gocommon/cron"
)

// DemoTask holds dependencies for the task.
type DemoTask struct {
	xcron.BaseTask
	repo repo.UserRepository
}

func NewDemoTask(repo repo.UserRepository) *DemoTask {
	return &DemoTask{
		BaseTask: xcron.BaseTask{},
		repo:     repo,
	}
}

// BeforeRun 任务执行前的准备工作
func (t *DemoTask) BeforeRun(ctx context.Context) error {
	return nil
}

func (t *DemoTask) Run(ctx context.Context) error {
	return nil
}

// AfterRun 任务执行后的清理工作
func (t *DemoTask) AfterRun(ctx context.Context) error {
	return nil
}
