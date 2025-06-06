package router

import (
	"net/http"

	"godemo/internal/wire"

	"github.com/gin-gonic/gin"
	"github.com/jessewkun/gocommon/cache"
	"github.com/jessewkun/gocommon/db"
	"github.com/jessewkun/gocommon/middleware"
	"github.com/jessewkun/gocommon/response"

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
			"db":    db.HealthCheck(),
			"cache": cache.HealthCheck(),
		}
		c.JSON(http.StatusOK, response.SuccessResp(c, data))
	})

	// swagger
	if gin.Mode() == gin.DebugMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
