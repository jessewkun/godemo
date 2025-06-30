package router

import (
	"net/http"
	"net/http/pprof"
	"time"

	"godemo/internal/wire"

	"github.com/gin-gonic/gin"
	"github.com/jessewkun/gocommon/db/mongodb"
	"github.com/jessewkun/gocommon/db/mysql"
	"github.com/jessewkun/gocommon/db/redis"
	"github.com/jessewkun/gocommon/middleware"
	"github.com/jessewkun/gocommon/response"
	"golang.org/x/time/rate"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 跨域配置
var crosConfig = middleware.CrosConfig{
	AllowedOrigins: map[string]bool{
		"http://xiedehao.cn":    true,
		"https://xiedehao.cn":   true,
		"http://localhost:5173": true, // 前端开发环境备用端口
		"http://127.0.0.1:5173": true, // 前端开发环境备用端口
	},
	AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "OPTIONS"},
	AllowHeaders: []string{"Content-Type, Authorization, Content-Length,Keep-Alive,credentials,Cache-Control,user,X-Requested-With,If-Modified-Since,Cache-Control,Pragma,Last-Modified,Accept,Accept-Encoding,Accept-Language,Connection,Host,Referer,User-Agent,Origin,Sec-Ch-Ua,Sec-Ch-Ua-Mobile,Sec-Ch-Ua-Platform,Sec-Fetch-Dest,Sec-Fetch-Mode,Sec-Fetch-Site"},
}

// pprof 限流配置 - 仅允许本地访问
var pprofRateLimitConfig = &middleware.RateLimiterConfig{
	GlobalLimiter:       nil, // 关闭全局限流
	IPLimiters:          make(map[string]*rate.Limiter),
	IPLastUsed:          make(map[string]time.Time),
	IPLimit:             0, // 对非白名单IP设置为0，即完全禁止
	IPBurst:             0,
	EnableIPLimit:       true,
	EnableLog:           true,
	CleanupInterval:     time.Minute * 10,
	IPExpiration:        time.Minute * 10,
	Whitelist:           []string{"127.0.0.1", "::1", "localhost"}, // 本地IP白名单
	WhitelistSkipGlobal: true,                                      // 白名单IP跳过全局限流
	WhitelistChecker:    nil,                                       // 使用内置白名单检查
}

// InitRouter 初始化路由
func InitRouter(r *gin.Engine) *gin.Engine {
	r.Use(middleware.Recovery(), middleware.Cros(crosConfig), middleware.Trace(), middleware.IOLog(nil))
	r.NoMethod(HandleNotFound)
	r.NoRoute(HandleNotFound)

	// 注册系统路由
	registerSystemRoutes(r)

	// 注册API路由
	registerAPIRoutes(r)

	return r
}

func HandleNotFound(c *gin.Context) {
	c.Status(http.StatusNotFound)
}

func registerSystemRoutes(r *gin.Engine) {
	// ping
	r.GET("/healthcheck/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// 组件探活
	r.GET("/healthcheck/active", func(c *gin.Context) {
		data := map[string]interface{}{
			"mysql":   mysql.HealthCheck(),
			"redis":   redis.HealthCheck(),
			"mongodb": mongodb.HealthCheck(),
		}
		c.JSON(http.StatusOK, response.SuccessResp(c, data))
	})

	// swagger
	if gin.Mode() == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// pprof 路由 - 仅允许本地访问
	pprofGroup := r.Group("/debug/pprof")
	pprofGroup.Use(middleware.RateLimiter(pprofRateLimitConfig))
	{
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
		pprofGroup.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
		pprofGroup.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
		pprofGroup.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
		pprofGroup.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
		pprofGroup.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
	}
}

func registerAPIRoutes(r *gin.Engine) {
	// 使用 wire 初始化依赖
	userHandler, err := wire.InitializeAPI()
	if err != nil {
		panic(err)
	}

	v1 := r.Group("/api/v1")
	{
		// 用户相关路由
		user := v1.Group("/users")
		{
			user.POST("", userHandler.Create) // 创建用户
			user.GET("", userHandler.List)    // 获取用户列表
		}
	}
}
