# godemo

基于 gocommon 和 gin 的 Web 服务框架封装

## 项目简介

本项目是一个基于 [gocommon](https://github.com/jessewkun/gocommon) 和 [gin](https://github.com/gin-gonic/gin) 的 Web 服务框架封装，提供了以下特性：

-   基于 gin 的 HTTP 服务
-   使用 wire 进行依赖注入
-   配置管理（使用 Viper）
-   日志管理
-   优雅启动和关闭

## 快速开始

### 环境要求

-   Go 1.24.11 或更高版本
-   MySQL
-   Redis

### 安装依赖

```bash
make mod
```

### 构建项目

```bash
# 克隆项目
git clone <repository-url>
cd godemo

make help              # 查看所有可用命令

# 主应用相关
make build             # 清理并构建主应用
make run               # 运行主应用（默认 debug 环境）
make stop              # 停止主应用
make status            # 查看主应用状态

# Cron定时任务相关
make build-cron        # 清理并构建cron应用
make run-cron          # 运行cron调度器（默认 debug 环境）
make run-cron-task TASK=<task_name>  # 手动执行指定任务
make stop-cron         # 停止cron应用
make status-cron       # 查看cron应用状态

# 开发工具
make clean             # 清理构建文件
make test              # 运行测试
make wire              # 生成依赖注入代码
make mod               # 整理 Go 模块

# 指定环境运行
make run ENV=debug     # 主应用开发环境，使用 config/debug.toml
make run ENV=test      # 主应用测试环境，使用 config/test.toml
make run ENV=release   # 主应用生产环境，使用 config/release.toml
make run-cron ENV=debug    # cron开发环境，使用 config/debug.toml
make run-cron ENV=test     # cron测试环境，使用 config/test.toml
make run-cron ENV=release  # cron生产环境，使用 config/release.toml
```

## 配置说明

配置文件位于 `config` 目录下，支持不同环境的配置：

-   `debug.toml`: 开发环境配置
-   `test.toml`: 测试环境配置
-   `release.toml`: 生产环境配置

## 开发指南

1. 下载该项目，全局搜索 godemo 替换为实际项目名称
2. 修改对应环境的配置文件
3. 添加新的 API 接口：

    - 在 `internal/router` 中定义路由
    - 在 `internal/dto` 中定义接口协议
    - 在 `internal/handler` 中创建处理器
    - 在 `internal/handler/provider.gp` 中添加新增的 handler
    - 在 `internal/service` 中实现业务逻辑
    - 在 `internal/service/provider.go` 中添加新增的 vervice
    - 在 `internal/repository` 中实现数据访问
    - 在 `internal/repository/provider.go` 中添加新增的 repository
    - 在 `internal/model` 中实现数据模型
    - 在 `internal/wire` 中注册依赖
    - 在 `internal/middleware` 中创建业务中间件

## 注意事项

1. 确保在运行服务前已正确配置数据库和 Redis 连接信息
2. 开发新功能时注意遵循项目的代码结构和规范
3. 提交代码前请运行测试确保功能正常

## Docker 使用说明

### 构建镜像

```bash
docker build -t godemo .
```

### 运行容器（指定不同环境配置文件）

默认使用开发环境配置（debug.toml）：

```bash
docker run --rm godemo
```

指定生产环境配置（release.toml）：

```bash
docker run --rm godemo -c ./config/release.toml
```

指定测试环境配置（test.toml）：

```bash
docker run --rm godemo -c ./config/test.toml
```
