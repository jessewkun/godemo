package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"godemo/internal/router"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/jessewkun/gocommon/alarm"
	"github.com/jessewkun/gocommon/cache"
	"github.com/jessewkun/gocommon/config"
	"github.com/jessewkun/gocommon/db"
	"github.com/jessewkun/gocommon/logger"
	"github.com/spf13/viper"
)

var configFile string
var baseConfig *config.BaseConfig

func init() {
	flag.StringVar(&configFile, "c", "config.yml", "config file path")
	flag.Parse()

	var err error
	baseConfig, err = config.LoadConfig(configFile)
	if err != nil {
		panic(fmt.Errorf("load config file %s error: %s\n", configFile, err))
	}

	// 监控配置文件是否变化
	// viper 会自动监控配置文件的变化，当配置文件发生变化时，viper 会自动更新配置信息
	// 但是 viper 不会自动更新结构体，所以需要手动更新结构体
	// 这里只是为了 debug config 可以被动态更新，其他情况下不建议使用
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
	})

	// 初始化各个组件
	err = initComponents()
	if err != nil {
		panic(fmt.Errorf("init components error: %s", err))
	}
}

// initComponents 初始化各个组件
func initComponents() error {
	err := logger.InitLogger(baseConfig.Log)
	if err != nil {
		return fmt.Errorf("init logger error: %s", err)
	}
	err = db.InitMysql(baseConfig.Mysql)
	if err != nil {
		return fmt.Errorf("init mysql error: %s", err)
	}
	err = cache.InitRedis(baseConfig.Redis)
	if err != nil {
		return fmt.Errorf("init redis error: %s", err)
	}
	err = alarm.InitBark(baseConfig.Alarm)
	if err != nil {
		return fmt.Errorf("init alarm error: %s", err)
	}
	return nil
}

func main() {
	gin.SetMode(baseConfig.Mode)
	r := gin.Default()

	srv := &http.Server{
		Addr:    baseConfig.Port,
		Handler: router.InitRouter(r),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("Server startup failed: %v", err))
		}
	}()

	gracefulExit(srv)
}

// gracefulExit 优雅退出
func gracefulExit(srv *http.Server) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-exit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.Info(ctx, "MAIN", "Received signal: %v. Shutting down server...", sig)

	if err := srv.Shutdown(ctx); err != nil {
		logger.ErrorWithMsg(ctx, "MAIN", "Server shutdown failed: %v", err)
	}
	logger.Info(ctx, "MAIN", "Server gracefully shutdown")
	fmt.Println("Server gracefully shutdown")
}
