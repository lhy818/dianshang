#!/bin/bash

# 电商网站项目启动脚本
# 用法: ./start.sh [dev|prod|infra]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Docker是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    print_success "Docker和Docker Compose已安装"
}

# 检查端口是否被占用
check_port() {
    local port=$1
    local service=$2
    
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        print_warning "端口 $port 已被占用 ($service)"
        return 1
    fi
    return 0
}

# 启动基础设施
start_infra() {
    print_info "启动基础设施服务..."
    
    # 检查端口
    check_port 5432 "PostgreSQL" || return 1
    check_port 6379 "Redis" || return 1
    check_port 5672 "RabbitMQ" || return 1
    check_port 15672 "RabbitMQ Management" || return 1
    
    # 启动服务
    docker-compose up -d postgres redis rabbitmq
    
    # 等待服务就绪
    print_info "等待服务启动..."
    sleep 10
    
    # 检查服务状态
    if docker-compose ps | grep -q "Up"; then
        print_success "基础设施服务启动成功"
        echo ""
        echo "服务状态:"
        echo "  PostgreSQL: localhost:5432"
        echo "  Redis: localhost:6379"
        echo "  RabbitMQ: localhost:5672"
        echo "  RabbitMQ管理界面: http://localhost:15672 (admin/admin)"
    else
        print_error "基础设施服务启动失败"
        docker-compose logs
        exit 1
    fi
}

# 启动开发环境
start_dev() {
    print_info "启动开发环境..."
    
    # 检查端口
    check_port 3000 "前端" || return 1
    check_port 8080 "后端" || return 1
    
    # 启动基础设施
    start_infra
    
    # 启动后端
    print_info "启动后端服务..."
    cd backend
    if [ ! -f go.mod ]; then
        print_error "后端项目目录不正确"
        exit 1
    fi
    
    # 安装依赖
    print_info "安装Go依赖..."
    go mod download
    
    # 启动服务（后台运行）
    print_info "启动Go服务..."
    go run cmd/main.go &
    BACKEND_PID=$!
    
    cd ..
    
    # 启动前端
    print_info "启动前端服务..."
    cd frontend
    if [ ! -f package.json ]; then
        print_error "前端项目目录不正确"
        exit 1
    fi
    
    # 安装依赖
    print_info "安装Node.js依赖..."
    npm install
    
    # 启动服务（后台运行）
    print_info "启动前端开发服务器..."
    npm run dev &
    FRONTEND_PID=$!
    
    cd ..
    
    # 保存PID到文件
    echo $BACKEND_PID > .backend.pid
    echo $FRONTEND_PID > .frontend.pid
    
    print_success "开发环境启动成功"
    echo ""
    echo "访问地址:"
    echo "  前端: http://localhost:3000"
    echo "  后端API: http://localhost:8080"
    echo "  后端健康检查: http://localhost:8080/api/v1/health"
    echo ""
    echo "按 Ctrl+C 停止所有服务"
    
    # 等待用户中断
    wait
}

# 停止服务
stop_services() {
    print_info "停止服务..."
    
    # 停止前端
    if [ -f .frontend.pid ]; then
        FRONTEND_PID=$(cat .frontend.pid)
        if kill -0 $FRONTEND_PID 2>/dev/null; then
            kill $FRONTEND_PID
            print_info "前端服务已停止"
        fi
        rm -f .frontend.pid
    fi
    
    # 停止后端
    if [ -f .backend.pid ]; then
        BACKEND_PID=$(cat .backend.pid)
        if kill -0 $BACKEND_PID 2>/dev/null; then
            kill $BACKEND_PID
            print_info "后端服务已停止"
        fi
        rm -f .backend.pid
    fi
    
    # 停止基础设施
    docker-compose down
    
    print_success "所有服务已停止"
}

# 清理环境
cleanup() {
    print_info "清理环境..."
    
    # 停止服务
    stop_services
    
    # 清理Docker资源
    docker-compose down -v
    
    # 清理临时文件
    rm -f .backend.pid .frontend.pid
    
    print_success "环境清理完成"
}

# 显示帮助
show_help() {
    echo "电商网站项目启动脚本"
    echo ""
    echo "用法: $0 [命令]"
    echo ""
    echo "命令:"
    echo "  dev     启动开发环境（前端+后端+基础设施）"
    echo "  infra   只启动基础设施（数据库、缓存、消息队列）"
    echo "  stop    停止所有服务"
    echo "  clean   清理所有资源"
    echo "  help    显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 dev      # 启动完整开发环境"
    echo "  $0 infra    # 只启动基础设施"
    echo "  $0 stop     # 停止所有服务"
}

# 主函数
main() {
    local command=${1:-"dev"}
    
    case $command in
        "dev")
            check_docker
            start_dev
            ;;
        "infra")
            check_docker
            start_infra
            ;;
        "stop")
            stop_services
            ;;
        "clean")
            cleanup
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "未知命令: $command"
            show_help
            exit 1
            ;;
    esac
}

# 捕获Ctrl+C
trap 'print_info "正在停止服务..."; stop_services; exit 0' INT

# 运行主函数
main "$@"