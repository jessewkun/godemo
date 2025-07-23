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

	"godemo/internal/dto"
	"godemo/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jessewkun/gocommon/common"
	"github.com/jessewkun/gocommon/config"
	"github.com/jessewkun/gocommon/logger"

	_ "godemo/config"

	_ "github.com/jessewkun/gocommon/debug"
	_ "github.com/jessewkun/gocommon/http"
)

var (
	configFile string
	baseConfig *config.BaseConfig
)

func init() {
	flag.StringVar(&configFile, "c", "config.yml", "config file path")
	flag.Parse()
}

func loadConfig() error {
	var err error
	baseConfig, err = config.Init(configFile)
	if err != nil {
		return fmt.Errorf("load config file %s error: %w", configFile, err)
	}
	return nil
}

func startServer() *http.Server {
	gin.SetMode(baseConfig.Mode)
	r := gin.Default()

	// 注册自定义验证函数
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		dto.RegisterValidator(v)
	}

	srv := &http.Server{
		Addr:         baseConfig.Port,
		Handler:      router.InitRouter(r),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("Starting server on %s\n", baseConfig.Port)

	go func() {
		log(context.Background(), "MAIN", "Starting server on %s", baseConfig.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log(context.Background(), "MAIN", "Server startup failed: %v", err)
			os.Exit(1)
		}
	}()

	return srv
}

// gracefulExit 优雅退出
func gracefulExit(srv *http.Server) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-exit
	log(context.Background(), "MAIN", "Received signal: %v. Shutting down server...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log(context.Background(), "MAIN", "Server shutdown failed: %v", err)
		os.Exit(1)
	}

	log(context.Background(), "MAIN", "Server gracefully shutdown")
	fmt.Println("Server gracefully shutdown")
}

func main() {
	if err := loadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	srv := startServer()

	gracefulExit(srv)
}

func log(c context.Context, tag string, msg string, args ...interface{}) {
	if common.IsDebug() {
		logger.Info(c, tag, msg, args...)
	} else {
		logger.InfoWithAlarm(c, tag, msg, args...)
	}
}
