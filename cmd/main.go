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

	"github.com/gin-gonic/gin"
	"github.com/jessewkun/gocommon/config"
	"github.com/jessewkun/gocommon/logger"

	_ "github.com/jessewkun/gocommon/debug"
)

var configFile string
var baseConfig *config.BaseConfig

func init() {
	flag.StringVar(&configFile, "c", "config.yml", "config file path")
	flag.Parse()

	var err error
	baseConfig, err = config.Init(configFile)
	if err != nil {
		panic(fmt.Errorf("load config file %s error: %s\n", configFile, err))
	}
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
