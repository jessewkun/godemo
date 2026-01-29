// Package config Prompt 加载与管理
package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

type Prompt struct {
	Key     string // 业务唯一标识，如 apps.html_code.v1
	Content string // 原始 Prompt 内容（模板）
}

// Render 渲染 Prompt（模板变量缺失直接报错）
func (p *Prompt) Render(data any) (string, error) {
	tpl, err := template.
		New(p.Key).
		Option("missingkey=error").
		Parse(p.Content)
	if err != nil {
		return "", fmt.Errorf("解析 prompt 模板失败 [%s]: %w", p.Key, err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("渲染 prompt 失败 [%s]: %w", p.Key, err)
	}
	return buf.String(), nil
}

//
// ========================
// Prompt Loader
// ========================
//

type promptCacheItem struct {
	prompt *Prompt
}

var (
	promptRoot = "config/prompt"

	promptCache      = make(map[string]*promptCacheItem)
	promptCacheMutex sync.RWMutex
)

// LoadPrompt 加载 Prompt
// key 格式：模块.类型.版本
// 示例：apps.demo.v1
func LoadPrompt(key string) (*Prompt, error) {
	// 1. 读缓存
	promptCacheMutex.RLock()
	if item, ok := promptCache[key]; ok {
		promptCacheMutex.RUnlock()
		return item.prompt, nil
	}
	promptCacheMutex.RUnlock()

	// 2. 解析 key
	module, name, version, err := parsePromptKey(key)
	if err != nil {
		return nil, err
	}

	// 3. 构建路径
	fileName := fmt.Sprintf("%s_%s.txt", name, version)
	filePath := filepath.Join(promptRoot, module, fileName)

	// 4. 读取文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取 prompt 文件失败 [%s]: %w", filePath, err)
	}

	prompt := &Prompt{
		Key:     key,
		Content: strings.TrimSpace(string(content)),
	}

	// 5. 写缓存
	promptCacheMutex.Lock()
	promptCache[key] = &promptCacheItem{prompt: prompt}
	promptCacheMutex.Unlock()

	return prompt, nil
}

// ReloadPrompt 清除指定 Prompt 缓存
func ReloadPrompt(key string) {
	promptCacheMutex.Lock()
	delete(promptCache, key)
	promptCacheMutex.Unlock()
}

// ClearPromptCache 清空所有 Prompt 缓存
func ClearPromptCache() {
	promptCacheMutex.Lock()
	promptCache = make(map[string]*promptCacheItem)
	promptCacheMutex.Unlock()
}

//
// ========================
// Prompt 组合
// ========================
//

// JoinPrompts 组合多个 Prompt
// 常用于：system + task + output
func JoinPrompts(prompts ...*Prompt) *Prompt {
	var contents []string
	var keys []string

	for _, p := range prompts {
		if p == nil {
			continue
		}
		contents = append(contents, p.Content)
		keys = append(keys, p.Key)
	}

	return &Prompt{
		Key:     strings.Join(keys, "+"),
		Content: strings.Join(contents, "\n\n"),
	}
}

// parsePromptKey
// apps.demo.v1
// → module=apps, name=demo, version=v1
func parsePromptKey(key string) (module, name, version string, err error) {
	parts := strings.Split(key, ".")
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf(
			"prompt key 格式错误，应为 模块.类型.版本，例如 apps.html_code.v1，当前: %s",
			key,
		)
	}

	module = parts[0]
	version = parts[len(parts)-1]
	name = strings.Join(parts[1:len(parts)-1], "_")

	return
}
