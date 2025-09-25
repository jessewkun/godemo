// Package cron 提供定时任务的具体实现
package cron

import (
	"context"
	"godemo/internal/repository"
	"time"

	"github.com/jessewkun/gocommon/cron"
)

// DemoTask
type DemoTask struct {
	cron.BaseTask
	userRepo repository.UserRepository
}

// NewDemoTask
func NewDemoTask(userRepo repository.UserRepository) *DemoTask {
	return &DemoTask{
		BaseTask: cron.BaseTask{
			TaskName:    "demo_task",
			TaskDesc:    "demo",
			TaskEnabled: true,
			TaskSpec:    "0 * * * * *",    // 每分钟执行一次
			TaskTimeout: 55 * time.Second, // 55秒超时，提前退出
		},
		userRepo: userRepo,
	}
}

// BeforeRun 任务执行前的准备工作
func (t *DemoTask) BeforeRun(ctx context.Context) error {
	return nil
}

func (t *DemoTask) Run(ctx context.Context) error {
	return nil
}

// AfterRun 任务执行后的准备工作
func (t *DemoTask) AfterRun(ctx context.Context) error {
	return nil
}
