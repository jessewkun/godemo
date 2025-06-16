# godemo

基于 gocommon 和 gin 的 Web 服务框架封装

## 项目简介

本项目是一个基于 [gocommon](https://github.com/jessewkun/gocommon) 和 [gin](https://github.com/gin-gonic/gin) 的 Web 服务框架封装，提供了以下特性：

-   基于 gin 的 HTTP 服务
-   使用 wire 进行依赖注入
-   集成 Swagger 文档
-   支持 Redis 缓存
-   支持 MySQL 数据库（使用 GORM）
-   配置管理（使用 Viper）
-   日志管理
-   优雅启动和关闭

## 项目结构

```
.
├── cmd/            # 主程序入口
├── config/         # 配置文件
│   ├── config.toml   # 项目启动实际使用的配置文件，请勿直接修改，使用make命令自动生成对应环境的该文件
│   ├── debug.toml    # 开发环境配置文件
│   ├── release.toml  # 生产环境配置文件
│   ├── test.toml     # 测试环境配置文件
├── internal/       # 业务逻辑实现
│   ├── handler/    # HTTP 处理器
│   ├── middleware/ # 业务中间件，通用的中间件在 gocommon 中
│   ├── model/      # 数据模型
│   ├── repository/ # 数据访问层
│   ├── service/    # 业务逻辑层
│   └── wire/       # 依赖注入配置
├── bin/           # 编译输出目录
└── Makefile       # 构建脚本
```

## 快速开始

### 环境要求

-   Go 1.23.10 或更高版本
-   MySQL
-   Redis

### 安装依赖

```bash
make mod
```

### 构建项目

项目支持三种构建模式：

1. 开发环境构建：

```bash
# 自动进行历史产物的清理，配置文件的生成，wire的注入
make debug
```

2. 测试环境构建：

```bash
# 自动进行历史产物的清理，配置文件的生成，wire的注入
make test
```

3. 生产环境构建：

```bash
# 自动进行历史产物的清理，配置文件的生成，wire的注入
make release
```

### 运行服务

```bash
make run
```

### 停止服务

```bash
make stop
```

### 生成 Swagger 文档

```bash
make swag
```

### 生成依赖注入代码

```bash
# 手动测试注入代码的生成
make wire
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
    - 在 `internal/service` 中实现业务逻辑
    - 在 `internal/repository` 中实现数据访问
    - 在 `internal/model` 中实现数据模型
    - 在 `internal/wire` 中注册依赖
    - 在 `internal/middleware` 中创建业务中间件

## 测试

运行测试并生成覆盖率报告：

```bash
make cover
```

## 注意事项

1. 确保在运行服务前已正确配置数据库和 Redis 连接信息
2. 开发新功能时注意遵循项目的代码结构和规范
3. 提交代码前请运行测试确保功能正常
