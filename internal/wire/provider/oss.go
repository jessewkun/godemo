// Package provider 提供 OSS 客户端实例
package provider

import (
	"fmt"
	"godemo/config"

	"github.com/jessewkun/gocommon/oss"
)

// OssClient 包装类型，便于 Wire 识别
type OssClient struct{ *oss.Oss }

// ProvideOssClient 提供 OSS 客户端实例
func ProvideOssClient() OssClient {
	client, err := oss.NewOssSimple(
		config.BusinessCfg.Oss.Endpoint,
		config.BusinessCfg.Oss.AccessKey,
		config.BusinessCfg.Oss.SecretKey,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create oss client: %w", err))
	}
	return OssClient{client}
}
