package app

import (
	"context"
	"godemo/config"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jessewkun/gocommon/common"
	"github.com/jessewkun/gocommon/logger"
)

// App is the central structure of an application. It manages the lifecycle of servers.
type App struct {
	name    string
	options *Options
	servers []Server
	mu      sync.Mutex
	cancel  func()
}

// NewApp creates a new App instance. It initializes options and the logger.
// The appName is used for logging purposes. configFile is the path to the configuration file.
func NewApp(appName string, configFile string) (*App, error) {
	opts, err := NewOptions(configFile)
	if err != nil {
		return nil, err
	}
	opts.BaseConfig.AppName = appName

	app := &App{
		name:    appName,
		options: opts,
	}
	return app, nil
}

// AddServer registers a server to be managed by the application.
func (a *App) AddServer(s ...Server) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.servers = append(a.servers, s...)
}

// Run starts the application, including all registered servers, and waits for a
// shutdown signal. It honors the lifecycle of the provided parent context.
func (a *App) Run(ctx context.Context) error {
	appCtx, cancel := context.WithCancel(ctx)
	a.cancel = cancel

	a.startServers(appCtx)
	a.waitForShutdown(appCtx)
	a.stopServers()

	return nil
}

// startServers starts all registered servers in separate goroutines.
func (a *App) startServers(ctx context.Context) {
	for _, s := range a.servers {
		srv := s
		go func() {
			if err := srv.Start(ctx); err != nil {
				a.log(ctx, "APP_START_ERROR", "Application '%s' failed to start server: %v", a.name, err)
				a.cancel()
			}
		}()
	}
	buildInfo := config.GetBuildInfo()
	a.log(ctx, "APP_START", "Application '%s' started successfully with %d server(s). Build Version: %s, Commit: %s, Build time: %s", a.name, len(a.servers), buildInfo.Version, buildInfo.Commit, buildInfo.BuildTime)
}

// waitForShutdown blocks until a shutdown signal is received or the application context is cancelled.
func (a *App) waitForShutdown(ctx context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-quit:
		a.log(ctx, "APP_SHUTDOWN", "Received signal: %v. Shutting down...", sig)
	case <-ctx.Done():
		a.log(ctx, "APP_SHUTDOWN", "Context cancelled. Shutting down...")
	}
}

// stopServers gracefully shuts down all registered servers.
func (a *App) stopServers() {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	var wg sync.WaitGroup
	for _, s := range a.servers {
		wg.Add(1)
		srv := s
		go func() {
			defer wg.Done()
			if err := srv.Stop(shutdownCtx); err != nil {
				a.log(shutdownCtx, "SERVER_STOP_ERROR", "Failed to stop server gracefully: %v", err)
			}
		}()
	}
	wg.Wait()

	a.log(context.Background(), "APP_SHUTDOWN", "Application '%s' gracefully shut down.", a.name)
}

// Log is a helper function for consistent logging.
// 这个函数是为了在非 debug 模式下，能够自动感知到服务器的启动和关闭，避免出现服务发布但是实际没有启动或启动失败，开发人员还不知道的情况。
func (a *App) log(c context.Context, tag string, msg string, args ...interface{}) {
	if common.IsDebug() {
		logger.Info(c, tag, msg, args...)
	} else {
		logger.InfoWithAlarm(c, tag, msg, args...)
	}
}

// Options returns the application's configured options.
func (a *App) Options() *Options {
	return a.options
}
