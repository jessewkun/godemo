.PHONY: help build run stop clean test wire mod

# 默认环境
ENV ?= debug

# 应用配置
BINARY_NAME = godemo
CMD_FILE = cmd/main.go
CONFIG_DIR = config

# 颜色定义
SUCCESS = \033[32m
ERROR = \033[31m
WARNING = \033[33m
RESET = \033[0m

# 帮助信息
help:
	@echo "可用的命令："
	@echo "  make build                           - 清理并构建应用"
	@echo "  make run [ENV=debug|test|release]    - 构建并运行应用"
	@echo "  make stop                            - 停止应用"
	@echo "  make status                          - 查看应用状态"
	@echo "  make clean                           - 清理构建文件"
	@echo "  make test                            - 运行测试"
	@echo "  make wire                            - 生成依赖注入代码"
	@echo "  make mod                             - 整理 Go 模块"
	@echo ""
	@echo "示例："
	@echo "  make run                             - 以 debug 环境运行"
	@echo "  make run ENV=debug                   - 以 debug 环境运行"
	@echo "  make run ENV=test                    - 以 test 环境运行"
	@echo "  make run ENV=release                 - 以 release 环境运行"
	@echo "  make build                           - 清理并构建应用"

# 构建应用
build: clean wire
	@echo "$(SUCCESS)===> 构建 $(BINARY_NAME)$(RESET)"
	@go build -o bin/$(BINARY_NAME) $(CMD_FILE)
	@chmod +x bin/$(BINARY_NAME)
	@echo "$(SUCCESS)===> 构建完成$(RESET)"

# 运行应用
run:
	@echo "$(SUCCESS)===> 启动 $(BINARY_NAME) [$(ENV) 环境]$(RESET)"
	@cp $(CONFIG_DIR)/$(ENV).toml $(CONFIG_DIR)/config.toml
	@mkdir -p logs
	@nohup bin/$(BINARY_NAME) -c $(CONFIG_DIR)/config.toml > logs/app.log 2>&1 &
	@sleep 2
	@make status

# 停止应用
stop:
	@echo "$(WARNING)===> 停止 $(BINARY_NAME)$(RESET)"
	@echo "检查进程: pgrep -f 'bin/$(BINARY_NAME)'"
	@if pgrep -f "bin/$(BINARY_NAME)" > /dev/null; then \
		echo "发现进程，正在停止..."; \
		pkill -f "bin/$(BINARY_NAME)" || true; \
		sleep 1; \
		# 再次检查进程是否真的停止了 \
		if pgrep -f "bin/$(BINARY_NAME)" > /dev/null; then \
			echo "进程仍在运行，强制停止..."; \
			pkill -9 -f "bin/$(BINARY_NAME)" || true; \
			sleep 1; \
		fi; \
		echo "$(SUCCESS)===> 应用已停止$(RESET)"; \
	else \
		echo "未发现进程"; \
		echo "$(SUCCESS)===> 应用未运行$(RESET)"; \
	fi

# 查看应用状态
status:
	@if pgrep -f "bin/$(BINARY_NAME)" > /dev/null; then \
		echo "$(SUCCESS)===> $(BINARY_NAME) 正在运行$(RESET)"; \
		ps aux | grep -v grep | grep "bin/$(BINARY_NAME)"; \
	else \
		echo "$(ERROR)===> $(BINARY_NAME) 未运行$(RESET)"; \
	fi

# 清理构建文件
clean:
	@echo "$(WARNING)===> 清理构建文件$(RESET)"
	@rm -rf bin/*
	@rm -f $(CONFIG_DIR)/config.toml
	@rm -rf logs/*
	@rm -f nohup.out
	@echo "$(SUCCESS)===> 清理完成$(RESET)"

# 运行测试
test:
	@echo "$(SUCCESS)===> 运行测试$(RESET)"
	@go test -v ./...

# 生成依赖注入代码
wire:
	@echo "$(SUCCESS)===> 生成依赖注入代码$(RESET)"
	@cd internal/wire && wire
	@echo "$(SUCCESS)===> 依赖注入代码生成完成$(RESET)"

# 整理 Go 模块
mod:
	@echo "$(SUCCESS)===> 整理 Go 模块$(RESET)"
	@go mod tidy
	@go mod download
	@echo "$(SUCCESS)===> 模块整理完成$(RESET)"



# 默认目标
default: help
