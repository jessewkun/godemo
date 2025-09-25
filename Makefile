.PHONY: help build run stop clean test wire mod fmt build-cron run-cron run-cron-task stop-cron status-cron
# 默认环境
ENV ?= debug
SHELL := /bin/bash

# 应用配置
CONFIG_DIR = config
BINARY_NAME = gomode
CMD_FILE = cmd/main.go

# Cron配置
CRON_BINARY_NAME = gomode-cron
CRON_CMD_FILE = cmd/cron/main.go

# 颜色定义
SUCCESS = \033[32m
ERROR = \033[31m
WARNING = \033[33m
RESET = \033[0m

# 帮助信息
help:
	@echo "可用的命令："
	@echo ""
	@echo "应用相关："
	@echo "  make build                           - 清理并构建应用"
	@echo "  make run [ENV=debug|test|release]    - 构建并运行应用"
	@echo "  make stop                            - 停止应用"
	@echo "  make status                          - 查看应用状态"
	@echo ""
	@echo "Cron相关："
	@echo "  make build-cron                      - 清理并构建cron应用"
	@echo "  make run-cron [ENV=debug|test|release] - 构建并运行cron调度器"
	@echo "  make run-cron-task TASK=<task_name> [ENV=debug|test|release] - 手动执行指定任务"
	@echo "  make stop-cron                       - 停止cron应用"
	@echo "  make status-cron                     - 查看cron应用状态"
	@echo ""
	@echo "开发工具："
	@echo "  make clean                           - 清理构建文件"
	@echo "  make test                            - 运行测试"
	@echo "  make wire                            - 生成依赖注入代码"
	@echo "  make mod                             - 整理 Go 模块"
	@echo "  make fmt                             - 格式化代码"
	@echo ""
	@echo "示例："
	@echo "  make run                             - 以 debug 环境运行应用"
	@echo "  make run-cron                        - 以 debug 环境运行cron调度器"
	@echo "  make run-cron-task TASK=statistics_student_english_words_sections - 执行指定任务"
	@echo "  make build-cron                      - 构建cron应用"

# 格式化代码
fmt:
	@echo -e "$(WARNING)===> 格式化代码$(RESET)"
	@go fmt ./...
	@go vet ./...
	@echo -e "$(SUCCESS)===> 格式化代码完成$(RESET)"

# 构建应用
build: clean wire fmt
	@echo -e "$(WARNING)===> 构建 $(BINARY_NAME)$(RESET)"
	@go build -o bin/$(BINARY_NAME) $(CMD_FILE)
	@chmod +x bin/$(BINARY_NAME)
	@echo -e "$(SUCCESS)===> 构建完成$(RESET)"

# 运行应用
run:
	@echo -e "$(SUCCESS)===> 启动 $(BINARY_NAME) [$(ENV) 环境]$(RESET)"
	@cp $(CONFIG_DIR)/$(ENV).toml $(CONFIG_DIR)/config.toml
	@mkdir -p logs
	@nohup bin/$(BINARY_NAME) -c $(CONFIG_DIR)/config.toml > logs/app.log 2>&1 &
	@sleep 2
	@make status

# 查看应用状态
status:
	@if pgrep -fx 'bin/$(BINARY_NAME) -c config/config.toml' > /dev/null; then \
		echo -e "$(SUCCESS)===> $(BINARY_NAME) 正在运行$(RESET)"; \
		echo -e "$(WARNING)===> 进程信息:$(RESET)"; \
		ps -efww | grep -v grep | grep 'bin/$(BINARY_NAME) -c config/config.toml' || true; \
	else \
		echo -e "$(ERROR)===> $(BINARY_NAME) 未运行$(RESET)"; \
	fi

# 停止应用
stop: status
	@if pgrep -fx 'bin/$(BINARY_NAME) -c config/config.toml' > /dev/null; then \
		echo -e "$(WARNING)===> 正在停止进程...$(RESET)"; \
		pkill -fx 'bin/$(BINARY_NAME) -c config/config.toml' || true; \
		sleep 1; \
		if pgrep -fx 'bin/$(BINARY_NAME) -c config/config.toml' > /dev/null; then \
			echo "进程仍在运行，强制停止..."; \
			pkill -9 -fx 'bin/$(BINARY_NAME) -c config/config.toml' || true; \
			sleep 1; \
		fi; \
		echo -e "$(SUCCESS)===> 应用已停止$(RESET)"; \
	else \
		echo -e "$(WARNING)===> 未发现进程，无需停止$(RESET)"; \
	fi

# 清理构建文件
clean:
	@echo -e "$(WARNING)===> 清理构建文件$(RESET)"
	@rm -rf bin/*
	@rm -f $(CONFIG_DIR)/config.toml
	@rm -rf logs/*
	@rm -f nohup.out
	@echo -e "$(SUCCESS)===> 清理完成$(RESET)"

# 运行测试
test:
	@echo -e "$(SUCCESS)===> 运行测试$(RESET)"
	@go test -v ./...

# 生成依赖注入代码
wire:
	@echo -e "$(WARNING)===> 生成依赖注入代码$(RESET)"
	@cd internal/wire && wire
	@echo -e "$(SUCCESS)===> 依赖注入代码生成完成$(RESET)"

# 整理 Go 模块
mod:
	@echo -e "$(WARNING)===> 整理 Go 模块$(RESET)"
	@go mod tidy
	@go mod download
	@echo -e "$(SUCCESS)===> 模块整理完成$(RESET)"

# 构建cron应用
build-cron: clean wire fmt
	@echo -e "$(WARNING)===> 构建 $(CRON_BINARY_NAME)$(RESET)"
	@go build -o bin/$(CRON_BINARY_NAME) $(CRON_CMD_FILE)
	@chmod +x bin/$(CRON_BINARY_NAME)
	@echo -e "$(SUCCESS)===> cron应用构建完成$(RESET)"

# 运行cron调度器
run-cron:
	@echo -e "$(SUCCESS)===> 启动 $(CRON_BINARY_NAME) [$(ENV) 环境]$(RESET)"
	@cp $(CONFIG_DIR)/$(ENV).toml $(CONFIG_DIR)/config.toml
	@mkdir -p logs
	@nohup bin/$(CRON_BINARY_NAME) -c $(CONFIG_DIR)/config.toml > logs/cron.log 2>&1 &
	@sleep 2
	@make status-cron

# 手动执行cron任务
run-cron-task:
ifndef TASK
	@echo -e "$(ERROR)===> 请指定任务名称，使用 TASK=<task_name>$(RESET)"
	@echo "示例: make run-cron-task TASK=statistics_student_english_words_sections"
	@exit 1
endif
	@echo -e "$(SUCCESS)===> 执行cron任务: $(TASK) [$(ENV) 环境]$(RESET)"
	@cp $(CONFIG_DIR)/$(ENV).toml $(CONFIG_DIR)/config.toml
	@mkdir -p logs
	@bin/$(CRON_BINARY_NAME) -c $(CONFIG_DIR)/config.toml -t $(TASK)

# 查看cron应用状态
status-cron:
	@if pgrep -fx 'bin/$(CRON_BINARY_NAME) -c config/config.toml' > /dev/null; then \
		echo -e "$(SUCCESS)===> $(CRON_BINARY_NAME) 正在运行$(RESET)"; \
		echo -e "$(WARNING)===> 进程信息:$(RESET)"; \
		ps -efww | grep -v grep | grep 'bin/$(CRON_BINARY_NAME) -c config/config.toml' || true; \
	else \
		echo -e "$(ERROR)===> $(CRON_BINARY_NAME) 未运行$(RESET)"; \
	fi

# 停止cron应用
stop-cron: status-cron
	@if pgrep -fx 'bin/$(CRON_BINARY_NAME) -c config/config.toml' > /dev/null; then \
		echo -e "$(WARNING)===> 正在停止进程...$(RESET)"; \
		pkill -fx 'bin/$(CRON_BINARY_NAME) -c config/config.toml' || true; \
		sleep 1; \
		if pgrep -fx 'bin/$(CRON_BINARY_NAME) -c config/config.toml' > /dev/null; then \
			echo "进程仍在运行，强制停止..."; \
			pkill -9 -fx 'bin/$(CRON_BINARY_NAME) -c config/config.toml' || true; \
			sleep 1; \
		fi; \
		echo -e "$(SUCCESS)===> cron应用已停止$(RESET)"; \
	else \
		echo -e "$(WARNING)===> 未发现进程，无需停止$(RESET)"; \
	fi

# 默认目标
default: help
