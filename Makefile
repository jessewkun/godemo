.PHONY: mod,run,cover,clean,debug,test,release,stop,swag,wire

# 自定义echo，避免颜色无法输出
# echo -e 参数有终端不识别
ECHO = printf

# 颜色定义
SUCCESS := \033[32m
ERROR := \033[31m
WARNING := \033[33m
RESET := \033[0m

CUR_PATH:=$(shell pwd)
APP_PATH:=$(CUR_PATH)
CONFIG_NAME:=$(CUR_PATH)/config.toml
CMD_FILE:=$(CUR_PATH)/cmd/main.go
BINARY_NAME = godemo

export GO111MODULE=on
export GOPROXY=https://goproxy.cn
# export GOSUMDB=off
export GO111MODULE=on

default: debug

clean:
	@$(ECHO) "================================\n"
	@$(ECHO) "   Cleaning godemo Service       \n"
	@$(ECHO) "================================\n"
	@rm -rf $(APP_PATH)/bin/*
	@rm -rf $(APP_PATH)/config/config.toml
	@rm -rf ./logs/*
	@rm -rf ./nohup.out
	@$(ECHO) "$(SUCCESS)$(BINARY_NAME) clean up completed$(RESET)\n"

mod:
	@go mod tidy -v
	@go mod download

debug: clean cmd/main.go go.sum go.mod wire
	@$(ECHO) "================================\n"
	@$(ECHO) "   Building godemo Service       \n"
	@$(ECHO) "================================\n"
	@cp $(CUR_PATH)/config/debug.toml $(CUR_PATH)/config/config.toml
	@go env
	@go build -o $(APP_PATH)/bin/$(BINARY_NAME) $(CMD_FILE)
	@$(ECHO) "$(SUCCESS)[debug] $(BINARY_NAME) build success$(RESET)\n"

test: clean cmd/main.go go.sum go.mod wire
	@$(ECHO) "================================\n"
	@$(ECHO) "   Building godemo Service       \n"
	@$(ECHO) "================================\n"
	@cp $(CUR_PATH)/config/test.toml $(CUR_PATH)/config/config.toml
	@go env
	@go build -o $(APP_PATH)/bin/$(BINARY_NAME) $(CMD_FILE)
	@$(ECHO) "$(SUCCESS)[test] $(BINARY_NAME) build success$(RESET)\n"

release: clean cmd/main.go go.sum go.mod wire
	@$(ECHO) "================================\n"
	@$(ECHO) "   Building godemo Service       \n"
	@$(ECHO) "================================\n"
	@cp $(CUR_PATH)/config/release.toml $(CUR_PATH)/config/config.toml
	@go env
	@go build -o $(APP_PATH)/bin/$(BINARY_NAME) $(CMD_FILE)
	@$(ECHO) "$(SUCCESS)[release] $(BINARY_NAME) build success$(RESET)\n"

run:
	@make -s stop
	@$(ECHO) "================================\n"
	@$(ECHO) "   Running godemo Service       \n"
	@$(ECHO) "================================\n"
	@nohup $(APP_PATH)/bin/$(BINARY_NAME) -c $(CUR_PATH)/config/config.toml > /dev/null 2>&1 &
	@make -s check-process

stop:
	@$(ECHO) "================================\n"
	@$(ECHO) "   Stoping godemo Service       \n"
	@$(ECHO) "================================\n"
	@ps -ef | grep bin/$(BINARY_NAME) | grep -v grep | awk '{print $$2}' | xargs -r kill -9
	@$(ECHO) "$(WARNING)$(BINARY_NAME) service is shutdown$(RESET)\n"

check-process:
	@if ps aux | grep -v grep | grep bin/$(BINARY_NAME); then \
		$(ECHO) "$(SUCCESS)$(BINARY_NAME) service is running$(RESET)\n"; \
	else \
		$(ECHO) "$(ERROR)$(BINARY_NAME) service is not running$(RESET)\n"; \
	fi

swag:
	@swag init

cover:
	@go vet $(APP_PATH)
	@go test -coverpkg="./..." -cover $(APP_PATH)/... -gcflags='all=-N -l'

wire:
	@$(ECHO) "================================\n"
	@$(ECHO) "   Wire generate godemo Service       \n"
	@$(ECHO) "================================\n"
	@cd internal/wire && rm -rf wire_gen.go && wire
	@$(ECHO) "$(SUCCESS)wire generate success$(RESET)\n"
