#!/bin/bash

# 部署脚本
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
APP_NAME="godemo"
APP_DIR="/root/godemo"
BACKUP_DIR="/root/godemo_backup"
LOG_FILE="$APP_DIR/logs/app.log"
PID_FILE="$APP_DIR/godemo.pid"
HEALTH_CHECK_URL="http://localhost:8001/healthcheck/ping"
MAX_WAIT_TIME=30

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 健康检查函数
health_check() {
    local retries=0
    local max_retries=10

    log_info "开始健康检查..."

    while [ $retries -lt $max_retries ]; do
        if curl -f -s "$HEALTH_CHECK_URL" > /dev/null 2>&1; then
            log_info "健康检查通过！"
            return 0
        fi

        retries=$((retries + 1))
        log_warn "健康检查失败，重试 $retries/$max_retries..."
        sleep 3
    done

    log_error "健康检查失败，服务可能未正常启动"
    return 1
}

# 停止服务
stop_service() {
    log_info "停止现有服务..."
    cd "$APP_DIR"

    # 先尝试 make stop
    if make stop; then
        log_info "make stop 执行成功"
    else
        log_warn "make stop 执行失败，检查进程状态..."
        # 检查进程是否真的还在运行
        if pgrep -f "bin/$APP_NAME" > /dev/null; then
            log_warn "进程仍在运行，尝试直接停止..."
            pkill -f "bin/$APP_NAME" || true
            sleep 1
            if pgrep -f "bin/$APP_NAME" > /dev/null; then
                log_warn "进程仍在运行，强制停止..."
                pkill -9 -f "bin/$APP_NAME" || true
                sleep 1
            fi
            log_info "进程已停止"
        else
            log_info "进程已经停止，make stop 的退出码可能是误报"
        fi
    fi
}

# 备份当前版本
backup_current() {
    if [ -d "$APP_DIR" ] && [ -f "$APP_DIR/bin/$APP_NAME" ]; then
        log_info "备份当前版本..."
        rm -rf "$BACKUP_DIR"
        cp -r "$APP_DIR" "$BACKUP_DIR"
        log_info "备份完成: $BACKUP_DIR"
    else
        log_info "没有可备份的版本，跳过备份"
    fi
}

# 启动服务
start_service() {
    local env=${1:-test}  # 默认使用 test 环境
    log_info "启动服务 [环境: $env]..."
    cd "$APP_DIR"
    make run ENV=$env
}

# 回滚函数
rollback() {
    log_error "部署失败，开始回滚..."

    stop_service

    if [ -d "$BACKUP_DIR" ]; then
        log_info "恢复备份版本..."
        rm -rf "$APP_DIR"
        mv "$BACKUP_DIR" "$APP_DIR"
        start_service "$env"

        if health_check; then
            log_info "回滚成功！"
            return 0
        else
            log_error "回滚失败！"
            return 1
        fi
    else
        log_error "没有可用的备份版本"
        return 1
    fi
}

# 主部署流程
main() {
    local env=${1:-test}  # 默认使用 test 环境
    log_info "开始部署 $APP_NAME [环境: $env]..."

    # 检查应用目录是否存在
    if [ ! -d "$APP_DIR" ]; then
        log_error "应用目录不存在: $APP_DIR"
        exit 1
    fi

    # 检查可执行文件是否存在
    if [ ! -f "$APP_DIR/bin/$APP_NAME" ]; then
        log_error "可执行文件不存在: $APP_DIR/bin/$APP_NAME"
        exit 1
    fi

    # 备份当前版本（在停止服务之前）
    backup_current

    # 停止现有服务
    stop_service

    # 启动新服务
    start_service "$env"

    # 等待服务启动
    log_info "等待服务启动..."
    sleep 5

    # 健康检查
    if health_check; then
        log_info "✅ 部署成功！"

        # 显示服务状态
        log_info "服务状态："
        cd "$APP_DIR"
        make status

        # 显示最新日志
        log_info "最新日志："
        tail -n 10 "$LOG_FILE" || true

        # 保留备份一段时间（用于回滚）
        log_info "部署成功，备份保留在: $BACKUP_DIR"

        exit 0
    else
        log_error "❌ 部署失败！"

        # 显示错误日志
        log_error "错误日志："
        tail -n 20 "$LOG_FILE" || true

        # 回滚
        if rollback; then
            exit 0
        else
            exit 1
        fi
    fi
}

# 清理旧备份
cleanup_backup() {
    if [ -d "$BACKUP_DIR" ]; then
        log_info "清理旧备份: $BACKUP_DIR"
        rm -rf "$BACKUP_DIR"
        log_info "备份清理完成"
    else
        log_info "没有找到旧备份"
    fi
}

# 执行主函数
if [ "$1" = "cleanup" ]; then
    cleanup_backup
else
    main "$@"
fi
