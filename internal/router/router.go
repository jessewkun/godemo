package router

import (
	"net/http"

	"godemo/config"
	godemoMiddleware "godemo/internal/middleware"
	"godemo/internal/wire"

	"github.com/gin-gonic/gin"
	"github.com/jessewkun/gocommon/db/mongodb"
	"github.com/jessewkun/gocommon/db/mysql"
	"github.com/jessewkun/gocommon/db/redis"
	"github.com/jessewkun/gocommon/middleware"
	"github.com/jessewkun/gocommon/response"
)

// InitRouter 初始化路由
func InitRouter(r *gin.Engine, apis *wire.APIs) *gin.Engine {
	r.Use(middleware.Trace(), godemoMiddleware.TrimMiddleware(), middleware.IOLog(nil), middleware.Recovery(), middleware.Prometheus(), middleware.Cros(config.BusinessCfg.Cros))
	r.NoMethod(HandleNotFound)
	r.NoRoute(HandleNotFound)

	// 注册系统路由
	registerSystemRoutes(r)

	// 注册API路由
	registerAPIRoutes(r, apis)

	return r
}

func HandleNotFound(c *gin.Context) {
	c.Status(http.StatusNotFound)
}

func registerSystemRoutes(r *gin.Engine) {
	// 组件探活
	r.GET("/health/check", func(c *gin.Context) {
		data := map[string]interface{}{
			"mysql":   mysql.HealthCheck(),
			"redis":   redis.HealthCheck(),
			"mongodb": mongodb.HealthCheck(),
		}
		response.Success(c, data)
	})
}

func registerAPIRoutes(r *gin.Engine, apis *wire.APIs) {
	v1 := r.Group("/api/v1")
	{
		// 用户相关路由
		user := v1.Group("/users")
		{
			user.POST("", apis.UserHandler.Create) // 创建用户
			user.GET("", apis.UserHandler.List)    // 获取用户列表
		}
	}
}
