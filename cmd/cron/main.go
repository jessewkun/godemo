package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"godemo/internal/cron"
	"godemo/internal/wire"

	"github.com/jessewkun/gocommon/common"
	xconfig "github.com/jessewkun/gocommon/config"
	"github.com/jessewkun/gocommon/logger"

	_ "godemo/config"

	_ "github.com/jessewkun/gocommon/debug"
)

var (
	configFile string
	baseConfig *xconfig.BaseConfig
	taskName   string
)

func init() {
	flag.StringVar(&configFile, "c", "config.yml", "config file path")
	flag.StringVar(&taskName, "t", "", "task name")
	flag.Parse()
}

func loadConfig() error {
	var err error
	baseConfig, err = xconfig.Init(configFile)
	if err != nil {
		return fmt.Errorf("load config file %s error: %w", configFile, err)
	}
	return nil
}

func gracefulShutdown(ctx context.Context, app *cron.App) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ssig := <-exit
	log(ctx, "CRON", "Received signal: %v. Shutting down cron manager...", ssig)

	app.Stop()

	log(ctx, "CRON", "Cron manager gracefully shutdown")
	fmt.Println("Cron manager gracefully shutdown")
}

func main() {
	if err := loadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}
	baseConfig.AppName = "godemo-cron"

	app, cleanup, err := wire.InitializeCronApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize dependencies: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	ctx := context.Background()

	// 如果指定了任务名称，手动执行该任务
	if taskName != "" {
		fmt.Printf("Start task: %s\n", taskName)

		err := app.RunTask(ctx, taskName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running task %s: %v\n", taskName, err)
			os.Exit(1)
		}

		fmt.Printf("Task %s executed successfully.\n", taskName)
		return
	}

	// 正常启动定时任务调度器
	if err := app.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start cron manager: %v\n", err)
		os.Exit(1)
	}

	log(ctx, "CRON", "Cron app started successfully")

	gracefulShutdown(ctx, app)
}

func log(c context.Context, tag string, msg string, args ...interface{}) {
	if common.IsDebug() {
		logger.Info(c, tag, msg, args...)
	} else {
		logger.InfoWithAlarm(c, tag, msg, args...)
	}
}
