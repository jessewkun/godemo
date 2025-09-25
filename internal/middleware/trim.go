package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// JSON 的递归 trim
func trimJSON(data interface{}) interface{} {
	switch v := data.(type) {
	case string:
		return strings.TrimSpace(v)
	case []interface{}:
		for i := range v {
			v[i] = trimJSON(v[i])
		}
		return v
	case map[string]interface{}:
		for k, val := range v {
			v[k] = trimJSON(val)
		}
		return v
	default:
		return v
	}
}

// TrimMiddleware 中间件，用于trim请求的body、query、form数据
//
// gin 的 shouldbindjson 会自动验证参数，所以需要在验证前全局 trim，之后再进行验证
func TrimMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ct := c.ContentType()

		// 处理 JSON
		if strings.Contains(ct, "application/json") {
			body, err := io.ReadAll(c.Request.Body)
			if err == nil && len(body) > 0 {
				var tmp interface{}
				if err := json.Unmarshal(body, &tmp); err == nil {
					tmp = trimJSON(tmp)
					newBody, _ := json.Marshal(tmp)
					c.Request.Body = io.NopCloser(bytes.NewReader(newBody))
					c.Request.ContentLength = int64(len(newBody))
				} else {
					// JSON 解析失败，恢复原 body
					c.Request.Body = io.NopCloser(bytes.NewReader(body))
				}
			}
		}

		// 处理 form / query
		if err := c.Request.ParseForm(); err == nil {
			for k, vals := range c.Request.Form {
				for i := range vals {
					vals[i] = strings.TrimSpace(vals[i])
				}
				c.Request.Form[k] = vals
			}
			c.Request.PostForm = c.Request.Form
		}

		c.Next()
	}
}
