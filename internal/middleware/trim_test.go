package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTrimJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "字符串trim",
			input:    "  hello world  ",
			expected: "hello world",
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "只有空格",
			input:    "   ",
			expected: "",
		},
		{
			name:     "数字类型",
			input:    123,
			expected: 123,
		},
		{
			name:     "布尔类型",
			input:    true,
			expected: true,
		},
		{
			name:     "nil值",
			input:    nil,
			expected: nil,
		},
		{
			name: "字符串数组",
			input: []interface{}{
				"  hello  ",
				"  world  ",
				"  test  ",
			},
			expected: []interface{}{
				"hello",
				"world",
				"test",
			},
		},
		{
			name: "嵌套对象",
			input: map[string]interface{}{
				"name":   "  John Doe  ",
				"email":  "  john@example.com  ",
				"age":    30,
				"active": true,
				"tags":   []interface{}{"  tag1  ", "  tag2  "},
				"address": map[string]interface{}{
					"street": "  123 Main St  ",
					"city":   "  New York  ",
				},
			},
			expected: map[string]interface{}{
				"name":   "John Doe",
				"email":  "john@example.com",
				"age":    30,
				"active": true,
				"tags":   []interface{}{"tag1", "tag2"},
				"address": map[string]interface{}{
					"street": "123 Main St",
					"city":   "New York",
				},
			},
		},
		{
			name: "混合类型数组",
			input: []interface{}{
				"  string  ",
				123,
				true,
				map[string]interface{}{
					"key": "  value  ",
				},
			},
			expected: []interface{}{
				"string",
				123,
				true,
				map[string]interface{}{
					"key": "value",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimJSON(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTrimMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		method        string
		contentType   string
		body          string
		queryParams   map[string]string
		formData      map[string]string
		expectedBody  string
		expectedQuery map[string]string
		expectedForm  map[string]string
	}{
		{
			name:        "JSON请求trim",
			method:      "POST",
			contentType: "application/json",
			body: `{
				"name": "  John Doe  ",
				"email": "  john@example.com  ",
				"tags": ["  tag1  ", "  tag2  "],
				"address": {
					"street": "  123 Main St  ",
					"city": "  New York  "
				}
			}`,
			expectedBody: `{"address":{"city":"New York","street":"123 Main St"},"email":"john@example.com","name":"John Doe","tags":["tag1","tag2"]}`,
		},
		{
			name:        "JSON请求空值",
			method:      "POST",
			contentType: "application/json",
			body: `{
				"name": "  ",
				"email": "",
				"description": null
			}`,
			expectedBody: `{"description":null,"email":"","name":""}`,
		},
		{
			name:         "JSON请求无效JSON",
			method:       "POST",
			contentType:  "application/json",
			body:         `{"name": "  John  ", "email": "john@example.com"`,
			expectedBody: `{"name": "  John  ", "email": "john@example.com"`,
		},
		{
			name:        "Query参数trim",
			method:      "GET",
			contentType: "application/json",
			queryParams: map[string]string{
				"name":  "  John Doe  ",
				"email": "  john@example.com  ",
				"age":   "  30  ",
			},
			expectedQuery: map[string]string{
				"name":  "John Doe",
				"email": "john@example.com",
				"age":   "30",
			},
		},
		{
			name:        "Form数据trim",
			method:      "POST",
			contentType: "application/x-www-form-urlencoded",
			formData: map[string]string{
				"name":  "  John Doe  ",
				"email": "  john@example.com  ",
				"age":   "  30  ",
			},
			expectedForm: map[string]string{
				"name":  "John Doe",
				"email": "john@example.com",
				"age":   "30",
			},
		},
		{
			name:         "空body",
			method:       "POST",
			contentType:  "application/json",
			body:         "",
			expectedBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试路由
			router := gin.New()
			router.Use(TrimMiddleware())

			// 创建测试处理器来捕获处理后的数据
			var capturedBody string
			var capturedQuery map[string]string
			var capturedForm map[string]string

			router.Any("/test", func(c *gin.Context) {
				// 读取处理后的body
				if c.Request.Body != nil {
					body, _ := io.ReadAll(c.Request.Body)
					capturedBody = string(body)
				}

				// 获取处理后的query参数（从Form中获取，因为中间件会处理）
				capturedQuery = make(map[string]string)
				for k, v := range c.Request.Form {
					if len(v) > 0 {
						capturedQuery[k] = v[0]
					}
				}

				// 获取处理后的form数据
				capturedForm = make(map[string]string)
				for k, v := range c.Request.PostForm {
					if len(v) > 0 {
						capturedForm[k] = v[0]
					}
				}

				c.JSON(200, gin.H{"status": "ok"})
			})

			// 构建请求
			var req *http.Request
			if tt.method == "GET" {
				req = httptest.NewRequest(tt.method, "/test", nil)
				// 添加query参数
				q := req.URL.Query()
				for k, v := range tt.queryParams {
					q.Add(k, v)
				}
				req.URL.RawQuery = q.Encode()
			} else {
				if tt.contentType == "application/x-www-form-urlencoded" {
					// Form数据
					formData := make([]string, 0)
					for k, v := range tt.formData {
						formData = append(formData, k+"="+v)
					}
					req = httptest.NewRequest(tt.method, "/test", strings.NewReader(strings.Join(formData, "&")))
					req.Header.Set("Content-Type", tt.contentType)
				} else {
					// JSON数据
					req = httptest.NewRequest(tt.method, "/test", strings.NewReader(tt.body))
					req.Header.Set("Content-Type", tt.contentType)
				}
			}

			// 执行请求
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证结果
			assert.Equal(t, 200, w.Code)

			// 验证JSON body
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, capturedBody)
			}

			// 验证query参数
			if tt.expectedQuery != nil {
				assert.Equal(t, tt.expectedQuery, capturedQuery)
			}

			// 验证form数据
			if tt.expectedForm != nil {
				assert.Equal(t, tt.expectedForm, capturedForm)
			}
		})
	}
}

func TestTrimMiddlewareEdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("空Content-Type", func(t *testing.T) {
		router := gin.New()
		router.Use(TrimMiddleware())

		router.POST("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("POST", "/test", strings.NewReader(`{"name": "  test  "}`))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("非JSON Content-Type", func(t *testing.T) {
		router := gin.New()
		router.Use(TrimMiddleware())

		router.POST("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("POST", "/test", strings.NewReader(`{"name": "  test  "}`))
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("ParseForm错误", func(t *testing.T) {
		router := gin.New()
		router.Use(TrimMiddleware())

		router.POST("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// 创建一个会导致ParseForm失败的请求
		req := httptest.NewRequest("POST", "/test", strings.NewReader("invalid form data"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.ContentLength = -1 // 这会导致ParseForm失败
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})
}

// 基准测试
func BenchmarkTrimJSON(b *testing.B) {
	data := map[string]interface{}{
		"name":  "  John Doe  ",
		"email": "  john@example.com  ",
		"tags":  []interface{}{"  tag1  ", "  tag2  "},
		"address": map[string]interface{}{
			"street": "  123 Main St  ",
			"city":   "  New York  ",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trimJSON(data)
	}
}

func BenchmarkTrimMiddleware(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TrimMiddleware())

	router.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	body := `{
		"name": "  John Doe  ",
		"email": "  john@example.com  ",
		"tags": ["  tag1  ", "  tag2  "],
		"address": {
			"street": "  123 Main St  ",
			"city": "  New York  "
		}
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/test", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
