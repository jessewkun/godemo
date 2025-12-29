// Package config 配置包，用于加载和解析业务配置
package config

import (
	"fmt"

	xconfig "github.com/jessewkun/gocommon/config"
	xcron "github.com/jessewkun/gocommon/cron"
	"github.com/jessewkun/gocommon/middleware"
	"github.com/spf13/viper"
)

// BusinessConfig 业务配置
type BusinessConfig struct {
	Cros  middleware.CrosConfig `mapstructure:"cros" json:"cros"` // 跨域配置
	Oss   OssConfig             `mapstructure:"oss" json:"oss"`   // oss 配置
	Crons []xcron.TaskConfig    `mapstructure:"crons" json:"crons"`
}

// Reload 重新加载 BusinessConfig 配置.
// business 模块的所有配置项都被认为是安全的，可以进行热更新.
func (c *BusinessConfig) Reload(v *viper.Viper) {
	if err := v.UnmarshalKey("business", c); err != nil {
		fmt.Printf("failed to reload business config: %v\n", err)
		return
	}
	fmt.Printf("business config reload success, config: %+v\n", c)
}

// OssConfig oss 配置
type OssConfig struct {
	Bucket         string `mapstructure:"bucket" json:"bucket"`
	PublicEndpoint string `mapstructure:"public_endpoint" json:"public_endpoint"` // 公网访问地址
	Endpoint       string `mapstructure:"endpoint" json:"endpoint"`               // 内网访问地址，上传最走这个
	RegionEndpoint string `mapstructure:"region_endpoint" json:"region_endpoint"` // 区域访问地址
	AccessKey      string `mapstructure:"access_key" json:"access_key"`
	SecretKey      string `mapstructure:"secret_key" json:"secret_key"`
	RoleArn        string `mapstructure:"role_arn" json:"role_arn"`
	Region         string `mapstructure:"region" json:"region"`
}

// BusinessCfg 业务配置，注册为全局变量，方便使用
var BusinessCfg = &BusinessConfig{}

func init() {
	xconfig.Register("business", BusinessCfg)
}

var (
	version   = "dev"
	commit    = "none"
	buildTime = "unknown"
)

// BuildInfo 构建信息
type BuildInfo struct {
	Version   string
	Commit    string
	BuildTime string
}

// GetBuildInfo 获取构建信息
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   version,
		Commit:    commit,
		BuildTime: buildTime,
	}
}
