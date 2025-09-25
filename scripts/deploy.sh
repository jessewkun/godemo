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
CRON_APP_NAME="godemo-cron"
APP_DIR="/root/godemo"
CRON_APP_DIR="/root/godemo-cron"
BACKUP_DIR="/root/godemo_backup"
CRON_BACKUP_DIR="/root/godemo-cron_backup"
LOG_FILE="$APP_DIR/logs/app.log"
CRON_LOG_FILE="$CRON_APP_DIR/logs/cron.log"
HEALTH_CHECK_URL="http://localhost:8001/health/ping"
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
    local component=${1:-"app"}  # 默认为app，可选：app, cron, all
    log_info "停止现有服务 [组件: $component]..."
    cd "$APP_DIR"

    if [ "$component" = "app" ] || [ "$component" = "all" ]; then
        # 停止主应用
        log_info "停止主应用..."
        if make stop; then
            log_info "make stop 执行成功"
        else
            log_warn "make stop 执行失败，检查进程状态..."
            if pgrep -f "bin/$APP_NAME" > /dev/null; then
                log_warn "主应用进程仍在运行，尝试直接停止..."
                pkill -f "bin/$APP_NAME" || true
                sleep 1
                if pgrep -f "bin/$APP_NAME" > /dev/null; then
                    log_warn "进程仍在运行，强制停止..."
                    pkill -9 -f "bin/$APP_NAME" || true
                    sleep 1
                fi
                log_info "主应用进程已停止"
            else
                log_info "主应用进程已经停止，make stop 的退出码可能是误报"
            fi
        fi
    fi

    if [ "$component" = "cron" ] || [ "$component" = "all" ]; then
        # 停止cron应用
        log_info "停止cron应用..."
        if make stop-cron; then
            log_info "make stop-cron 执行成功"
        else
            log_warn "make stop-cron 执行失败，检查进程状态..."
            if pgrep -f "bin/$CRON_APP_NAME" > /dev/null; then
                log_warn "cron进程仍在运行，尝试直接停止..."
                pkill -f "bin/$CRON_APP_NAME" || true
                sleep 1
                if pgrep -f "bin/$CRON_APP_NAME" > /dev/null; then
                    log_warn "进程仍在运行，强制停止..."
                    pkill -9 -f "bin/$CRON_APP_NAME" || true
                    sleep 1
                fi
                log_info "cron进程已停止"
            else
                log_info "cron进程已经停止，make stop-cron 的退出码可能是误报"
            fi
        fi
    fi
}

# 备份当前版本
backup_current() {
    local component=${1:-"app"}  # 默认为app，可选：app, cron, all

    if [ "$component" = "app" ] || [ "$component" = "all" ]; then
        if [ -d "$APP_DIR" ] && [ -f "$APP_DIR/bin/$APP_NAME" ]; then
            log_info "备份主应用版本..."
            rm -rf "$BACKUP_DIR"
            cp -r "$APP_DIR" "$BACKUP_DIR"
            log_info "备份完成: $BACKUP_DIR"
        else
            log_info "主应用目录不存在或可执行文件不存在，跳过备份"
        fi
    fi

    if [ "$component" = "cron" ] || [ "$component" = "all" ]; then
        if [ -d "$CRON_APP_DIR" ] && [ -f "$CRON_APP_DIR/bin/$CRON_APP_NAME" ]; then
            log_info "备份cron应用版本..."
            rm -rf "$CRON_BACKUP_DIR"
            cp -r "$CRON_APP_DIR" "$CRON_BACKUP_DIR"
            log_info "备份完成: $CRON_BACKUP_DIR"
        else
            log_info "cron应用目录不存在或可执行文件不存在，跳过备份"
        fi
    fi
}

# 启动服务
start_service() {
    local env=${1:-test}  # 默认使用 test 环境
    local component=${2:-"app"}  # 默认为app，可选：app, cron, all
    log_info "启动服务 [环境: $env, 组件: $component]..."

    if [ "$component" = "app" ] || [ "$component" = "all" ]; then
        log_info "启动主应用..."
        cd "$APP_DIR"
        make run ENV=$env
    fi

    if [ "$component" = "cron" ] || [ "$component" = "all" ]; then
        log_info "启动cron应用..."
        cd "$CRON_APP_DIR"
        make run-cron ENV=$env
    fi
}

# 回滚函数
rollback() {
    local component=${1:-"app"}  # 默认为app，可选：app, cron, all
    log_error "部署失败，开始回滚 [组件: $component]..."

    stop_service "$component"

    if [ "$component" = "app" ] || [ "$component" = "all" ]; then
        if [ -d "$BACKUP_DIR" ]; then
            log_info "恢复主应用备份版本..."
            rm -rf "$APP_DIR"
            mv "$BACKUP_DIR" "$APP_DIR"
        else
            log_error "没有可用的主应用备份版本"
            return 1
        fi
    fi

    if [ "$component" = "cron" ] || [ "$component" = "all" ]; then
        if [ -d "$CRON_BACKUP_DIR" ]; then
            log_info "恢复cron应用备份版本..."
            rm -rf "$CRON_APP_DIR"
            mv "$CRON_BACKUP_DIR" "$CRON_APP_DIR"
        else
            log_error "没有可用的cron应用备份版本"
            return 1
        fi
    fi

    start_service "$env" "$component"

    # 对于cron组件，跳过健康检查
    if [ "$component" = "cron" ]; then
        log_info "cron回滚成功！"
        return 0
    elif health_check; then
        log_info "回滚成功！"
        return 0
    else
        log_error "回滚失败！"
        return 1
    fi
}

# 主部署流程
main() {
    # 检查参数
    if [ $# -lt 2 ]; then
        log_error "参数不足！"
        show_usage
        exit 1
    fi

    local env=$1
    local component=$2
    local app_name=""

    # 验证环境参数
    if [ "$env" != "test" ] && [ "$env" != "release" ]; then
        log_error "无效的环境参数: $env"
        log_error "支持的环境: test, release"
        exit 1
    fi

    # 验证组件参数
    if [ "$component" != "app" ] && [ "$component" != "cron" ] && [ "$component" != "all" ]; then
        log_error "无效的组件参数: $component"
        log_error "支持的组件: app, cron, all"
        exit 1
    fi

    # 根据组件确定应用名称
    if [ "$component" = "app" ]; then
        app_name="$APP_NAME"
    elif [ "$component" = "cron" ]; then
        app_name="$CRON_APP_NAME"
    elif [ "$component" = "all" ]; then
        app_name="$APP_NAME + $CRON_APP_NAME"
    fi

    log_info "开始部署 $app_name [环境: $env, 组件: $component]..."

    # 检查应用目录和可执行文件是否存在
    if [ "$component" = "app" ] || [ "$component" = "all" ]; then
        if [ ! -d "$APP_DIR" ]; then
            log_error "主应用目录不存在: $APP_DIR"
            exit 1
        fi
        if [ ! -f "$APP_DIR/bin/$APP_NAME" ]; then
            log_error "主应用可执行文件不存在: $APP_DIR/bin/$APP_NAME"
            exit 1
        fi
    fi

    if [ "$component" = "cron" ] || [ "$component" = "all" ]; then
        if [ ! -d "$CRON_APP_DIR" ]; then
            log_error "cron应用目录不存在: $CRON_APP_DIR"
            exit 1
        fi
        if [ ! -f "$CRON_APP_DIR/bin/$CRON_APP_NAME" ]; then
            log_error "cron应用可执行文件不存在: $CRON_APP_DIR/bin/$CRON_APP_NAME"
            exit 1
        fi
    fi

    # 备份当前版本（在停止服务之前）
    backup_current "$component"

    # 停止现有服务
    stop_service "$component"

    # 启动新服务
    start_service "$env" "$component"

    # 等待服务启动
    log_info "等待服务启动..."
    sleep 5

    # 健康检查（仅对主应用进行）
    if [ "$component" = "app" ] || [ "$component" = "all" ]; then
        if health_check; then
            log_info "✅ 主应用部署成功！"
        else
            log_error "❌ 主应用部署失败！"
            log_error "错误日志："
            tail -n 20 "$LOG_FILE" || true
            if rollback "$component"; then
                exit 0
            else
                exit 1
            fi
        fi
    fi

    if [ "$component" = "cron" ] || [ "$component" = "all" ]; then
        log_info "✅ cron应用部署成功！"
    fi

    # 显示服务状态
    log_info "服务状态："
    cd "$APP_DIR"
    if [ "$component" = "app" ] || [ "$component" = "all" ]; then
        make status
    fi
    if [ "$component" = "cron" ] || [ "$component" = "all" ]; then
        make status-cron
    fi

    # 显示最新日志
    log_info "最新日志："
    if [ "$component" = "app" ] || [ "$component" = "all" ]; then
        log_info "主应用日志："
        tail -n 10 "$LOG_FILE" || true
    fi
    if [ "$component" = "cron" ] || [ "$component" = "all" ]; then
        log_info "cron应用日志："
        tail -n 10 "$CRON_LOG_FILE" || true
    fi

    # 保留备份一段时间（用于回滚）
    log_info "部署成功，备份保留在: $BACKUP_DIR"
    exit 0
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

# 显示使用说明
show_usage() {
    echo "使用方法: $0 [环境] [组件]"
    echo ""
    echo "必需参数:"
    echo "  环境      - test | release"
    echo "  组件      - app | cron | all"
    echo ""
    echo "环境:"
    echo "  test      - 测试环境"
    echo "  release   - 生产环境"
    echo ""
    echo "组件:"
    echo "  app       - 主应用"
    echo "  cron      - Cron定时任务"
    echo "  all       - 主应用 + Cron定时任务"
    echo ""
    echo "特殊命令:"
    echo "  cleanup   - 清理备份"
    echo ""
    echo "示例:"
    echo "  $0 test app          # 部署主应用到测试环境"
    echo "  $0 test cron         # 部署cron到测试环境"
    echo "  $0 test all          # 部署主应用和cron到测试环境"
    echo "  $0 release app       # 部署主应用到生产环境"
    echo "  $0 release cron      # 部署cron到生产环境"
    echo "  $0 release all       # 部署主应用和cron到生产环境"
    echo "  $0 cleanup           # 清理备份"
}

# 执行主函数
if [ "$1" = "cleanup" ]; then
    cleanup_backup
elif [ "$1" = "help" ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_usage
else
    main "$@"
fi
