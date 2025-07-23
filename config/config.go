// Package config 配置包，用于加载和解析业务配置
package config

import (
	"fmt"

	xconfig "github.com/jessewkun/gocommon/config"
	"github.com/spf13/viper"
)

// BusinessConfig 业务配置
type BusinessConfig struct {
	Token               Token               `mapstructure:"token" json:"token"`                               // 登录态 token 加解密配置
	PersonalInformation PersonalInformation `mapstructure:"personal_information" json:"personal_information"` // 个人信息加解密配置
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

// Token 登录态 token 加解密配置
type Token struct {
	Key string `mapstructure:"key" json:"key"` // 加密密钥
	Iv  string `mapstructure:"iv" json:"iv"`   // 加密向量
}

// PersonalInformation 个人信息加解密配置
type PersonalInformation struct {
	Key string `mapstructure:"key" json:"key"` // 加密密钥
	Iv  string `mapstructure:"iv" json:"iv"`   // 加密向量
}

// BusinessCfg 业务配置，注册为全局变量，方便使用
var BusinessCfg = &BusinessConfig{}

func init() {
	xconfig.Register("business", BusinessCfg)
}
