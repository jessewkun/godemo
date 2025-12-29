package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"godemo/internal/app"
	"godemo/internal/cron"
	"godemo/internal/wire"

	_ "godemo/config"

	_ "github.com/jessewkun/gocommon/debug"
)

// cronServer adapts the cron.App to the app.Server interface.
type cronServer struct {
	cronApp *cron.App
}

// Start begins the cron scheduler.
func (s *cronServer) Start(ctx context.Context) error {
	return s.cronApp.Start(ctx)
}

// Stop halts the cron scheduler.
func (s *cronServer) Stop(ctx context.Context) error {
	s.cronApp.Stop(ctx)
	return nil
}

// 主函数
// 这个函数是整个应用的入口，它负责创建应用程序实例，初始化API服务器，并启动应用程序。
// 这里的错误直接输出，没有进入日志，是因为可以在 github action 中看到错误或者手动运行时看到错误，方便排查问题。
func main() {
	var configFile string
	var taskName string

	flag.StringVar(&configFile, "c", "config.yml", "config file path")
	flag.StringVar(&taskName, "t", "", "task name to run manually")
	flag.Parse()

	application, err := app.NewApp("godemo-cron", configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create app: %v\n", err)
		os.Exit(1)
	}

	cronApp, cleanup, err := wire.InitializeCron()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize dependencies: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	ctx := context.Background()

	// 如果指定了任务名称，则运行指定的任务
	if taskName != "" {
		fmt.Fprintf(os.Stdout, "Starting single task: %s\n", taskName)
		err := cronApp.RunTask(ctx, taskName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running task %s: %v\n", taskName, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Task %s executed successfully.\n", taskName)
		os.Exit(0)
	}

	cronSrv := &cronServer{cronApp: cronApp}
	application.AddServer(cronSrv)

	if err := application.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Application run failed: %v\n", err)
		os.Exit(1)
	}
}
