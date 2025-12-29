package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"godemo/internal/app"
	"godemo/internal/dto"
	"godemo/internal/router"
	"godemo/internal/wire"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	_ "godemo/config"

	_ "github.com/jessewkun/gocommon/debug"
	_ "github.com/jessewkun/gocommon/http"
)

// apiServer wraps the HTTP server and its dependencies to implement the app.Server interface.
type apiServer struct {
	srv  *http.Server
	apis *wire.APIs
}

// Start begins listening for HTTP requests.
func (s *apiServer) Start(ctx context.Context) error {
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server startup failed: %w", err)
	}
	return nil
}

// Stop gracefully shuts down the HTTP server and closes associated resources.
func (s *apiServer) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}
	return nil
}

// newAPIServer encapsulates the creation and configuration of the HTTP server.
func newAPIServer(opts *app.Options, apis *wire.APIs) (*apiServer, error) {
	gin.SetMode(opts.BaseConfig.Mode)
	r := gin.New()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		dto.RegisterValidator(v)
	}

	srv := &http.Server{
		Addr:         opts.BaseConfig.Port,
		Handler:      router.InitRouter(r, apis),
		ReadTimeout:  5 * 60 * time.Second,
		WriteTimeout: 5 * 60 * time.Second,
		IdleTimeout:  300 * time.Second,
	}

	return &apiServer{
		srv:  srv,
		apis: apis,
	}, nil
}

// 主函数
// 这个函数是整个应用的入口，它负责创建应用程序实例，初始化API服务器，并启动应用程序。
// 这里的错误直接输出，没有进入日志，是因为可以在 github action 中看到错误或者手动运行时看到错误，方便排查问题。
func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "config.yml", "config file path")
	flag.Parse()

	application, err := app.NewApp("godemo-api", configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create app: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	apis, err := wire.InitializeAPIs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize APIs: %v\n", err)
		os.Exit(1)
	}

	apiSrv, err := newAPIServer(application.Options(), apis)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create API server: %v\n", err)
		os.Exit(1)
	}
	application.AddServer(apiSrv)
	if err := application.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Application run failed: %v\n", err)
		os.Exit(1)
	}
}
