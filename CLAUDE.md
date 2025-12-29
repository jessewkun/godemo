# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

后端服务，采用 Go 语言开发。

**技术栈**：Go 1.24.10 + Gin v1.10.1 + MySQL + Redis + Wire v0.6.0

**系统组件**：
- `godemo` - 主应用，提供 HTTP API 服务
- `godemo-cron` - 定时任务调度器

**详细的命令、部署流程等信息请参考 [README.md](./README.md)**

## 架构设计

### 分层架构

```
Request → Router → Middleware → Handler → Service → Repository/Model → Database
```

**各层职责**：
- **Router** (`internal/router/`): 路由定义和中间件应用
- **Middleware** (`internal/middleware/`): 认证、角色权限、参数处理
- **Handler** (`internal/handler/`): HTTP 请求处理和参数验证
- **Service** (`internal/service/`): 核心业务逻辑
- **Repository** (`internal/repository/`): 缓存层
- **Model** (`internal/model/`): 数据库模型（GORM）

### 依赖注入（Wire）

**关键文件**：
- `internal/wire/wire.go` - 依赖定义（⚠️ **修改后必须运行 `make wire`**）
- `internal/wire/wire_gen.go` - 自动生成，不要手动修改
- `internal/wire/provider/` - 基础设施提供者（DB、Redis、OSS 等）

## 关键约定

### 配置管理

- **环境配置**：`config/{debug|test|release}.toml`
- **运行时配置**：`config/config.toml`（由 Makefile 从环境配置复制）
- **环境切换**：通过 `ENV` 参数，如 `make run ENV=test`

### 日志记录

使用 `github.com/jessewkun/gocommon/logger`：
- `logger.Info(ctx, tag, msg, args...)` - 普通日志
- `logger.InfoWithAlarm(ctx, tag, msg, args...)` - 生产环境带告警

`main.go` 中的 `log()` 函数根据环境自动选择日志级别。

### 系统路由

由 `gocommon/router` 提供：
- `/healthcheck/ping` - 健康检查
- `/health/check` - 组件状态（MySQL、Redis）
- `/metrics` - Prometheus 指标
- `/debug/pprof` - 性能分析

## 开发工作流

### 添加新 API 端点

1. 在 `internal/handler/` 创建或修改 Handler
2. 在 `internal/service/` 实现业务逻辑
3. 在 `internal/router/router.go` 注册路由
4. 如需新依赖，更新 `internal/wire/wire.go` 并运行 `make wire`
5. 测试：`make test` → `make build && make run ENV=debug`

### 添加定时任务

1. 在 `internal/cron/` 创建任务实现
2. 在 `internal/wire/cron.go` 注册任务
3. 运行 `make wire` 生成依赖注入代码
4. 手动测试：`make build-cron && make run-cron-task TASK=<task_name> ENV=debug`
5. 调度测试：`make run-cron ENV=debug`

### 修改数据库模型

1. 修改 `internal/model/` 中的模型定义
2. 确保 GORM 标签正确（`gorm:"column:xxx"`）
3. 如涉及新表或字段，需要数据库迁移
4. 运行测试验证

### 数据库迁移脚本

**存储规则**：
- **目录结构**：`db/migrations/{模块}/{版本号}/`
  - 按业务模块分类：`user/`（用户模块）、`tiku/`（题库模块）等
  - 按版本号组织：`1.0.48/`、`1.0.49/` 等
  - 示例：`db/migrations/user/1.0.48/`

**文件命名规范**：`{对象类型}_{操作}_{对象名}.sql`
- `table_create_` - 创建表（如 `table_create_app_versions.sql`）
- `table_alter_` - 修改表结构
- `index_create_` - 创建索引
- `trigger_create_` - 创建触发器（如 `trigger_create_trg_app_versions_after_update.sql`）
- `data_insert_` - 插入初始数据
- `data_update_` - 数据迁移

**使用建议**：
- 每个版本的数据库变更独立存放在对应版本目录下
- 文件名应清晰描述操作内容，便于追溯和回滚
- 相同版本的多个 SQL 文件按需执行顺序命名（可添加数字前缀）

## 开发规范

**详细规范参考 [.cursor/rules](./.cursor/rules/)**：
- `general.mdc` - 通用规范
- `golang.mdc` - Go 开发规范
- `git.mdc` - Git 提交规范
- `document.mdc` - 文档规范

**核心设计原则**：
- 🔑 **DRY（Don't Repeat Yourself）**：避免重复代码，提取共用逻辑到独立函数
- 🔑 **KISS（Keep It Simple, Stupid）**：保持代码简单直接，避免过度设计
- 单一职责：每个函数/类只做一件事，保持可测试性
- 优先使用成熟的库和工具，避免不必要的自定义实现

**关键要求**：
- ⚠️ **不要自动提交代码，除非有明确提示**
- 始终使用中文描述 Git commit message
- Git 格式：`[type]: [description]`（type: feat/fix/docs/refactor/test 等）
- 提交前必须通过：`go fmt ./...`、`go vet ./...`、`go test ./...`
- 代码完整性：不留 todos、占位符或缺失部分
- 安全优先：避免 SQL 注入、XSS 等漏洞

## 参考文档

- **命令和部署**：[README.md](./README.md)
- **开发规范**：[.cursor/rules](./.cursor/rules/)
