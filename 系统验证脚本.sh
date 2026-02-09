#!/bin/bash

# 电商系统集成验证脚本
# 版本: 1.0
# 日期: 2026-02-08

set -e

echo "========================================="
echo "电商系统集成验证脚本"
echo "========================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

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

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        log_error "命令 $1 未安装"
        exit 1
    fi
}

# 检查Docker服务状态
check_docker() {
    log_info "检查Docker服务状态..."
    if ! docker info &> /dev/null; then
        log_error "Docker服务未运行"
        exit 1
    fi
    log_info "Docker服务正常"
}

# 检查基础设施服务
check_infrastructure() {
    log_info "检查基础设施服务..."
    
    services=("postgres" "redis" "rabbitmq")
    for service in "${services[@]}"; do
        container_name="ecommerce-$service"
        if docker ps --format '{{.Names}}' | grep -q "^${container_name}$"; then
            log_info "✓ $service 服务运行中"
        else
            log_warn "⚠ $service 服务未运行"
        fi
    done
}

# 检查后端服务
check_backend() {
    log_info "检查后端服务..."
    
    if docker ps --format '{{.Names}}' | grep -q "^ecommerce-backend$"; then
        log_info "✓ 后端服务运行中"
        
        # 测试健康检查端点
        log_info "测试后端健康检查..."
        if curl -s http://localhost:8080/health &> /dev/null; then
            log_info "✓ 后端健康检查通过"
        else
            log_warn "⚠ 后端健康检查失败"
        fi
    else
        log_warn "⚠ 后端服务未运行"
    fi
}

# 检查前端服务
check_frontend() {
    log_info "检查前端服务..."
    
    if docker ps --format '{{.Names}}' | grep -q "^ecommerce-frontend$"; then
        log_info "✓ 前端服务运行中"
    else
        log_warn "⚠ 前端服务未运行"
    fi
}

# 检查网络连通性
check_network() {
    log_info "检查服务网络连通性..."
    
    # 检查Docker网络
    if docker network ls | grep -q "ecommerce-project_default"; then
        log_info "✓ Docker网络存在"
    else
        log_warn "⚠ Docker网络不存在"
    fi
}

# 检查数据库连接
check_database() {
    log_info "检查数据库连接..."
    
    if docker ps --format '{{.Names}}' | grep -q "^ecommerce-postgres$"; then
        # 尝试连接数据库
        if docker exec ecommerce-postgres pg_isready -U postgres &> /dev/null; then
            log_info "✓ 数据库连接正常"
            
            # 检查表数量
            table_count=$(docker exec ecommerce-postgres psql -U postgres -d ecommerce -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null || echo "0")
            log_info "✓ 数据库表数量: $table_count"
        else
            log_warn "⚠ 数据库连接失败"
        fi
    else
        log_warn "⚠ 数据库服务未运行"
    fi
}

# 检查Redis连接
check_redis() {
    log_info "检查Redis连接..."
    
    if docker ps --format '{{.Names}}' | grep -q "^ecommerce-redis$"; then
        if docker exec ecommerce-redis redis-cli ping | grep -q "PONG"; then
            log_info "✓ Redis连接正常"
        else
            log_warn "⚠ Redis连接失败"
        fi
    else
        log_warn "⚠ Redis服务未运行"
    fi
}

# 检查RabbitMQ连接
check_rabbitmq() {
    log_info "检查RabbitMQ连接..."
    
    if docker ps --format '{{.Names}}' | grep -q "^ecommerce-rabbitmq$"; then
        if docker exec ecommerce-rabbitmq rabbitmq-diagnostics ping &> /dev/null; then
            log_info "✓ RabbitMQ连接正常"
        else
            log_warn "⚠ RabbitMQ连接失败"
        fi
    else
        log_warn "⚠ RabbitMQ服务未运行"
    fi
}

# 性能测试（基础）
performance_test() {
    log_info "执行基础性能测试..."
    
    # 测试API响应时间
    if curl -s http://localhost:8080/health &> /dev/null; then
        start_time=$(date +%s%N)
        curl -s http://localhost:8080/health &> /dev/null
        end_time=$(date +%s%N)
        response_time=$((($end_time - $start_time)/1000000))
        log_info "✓ API响应时间: ${response_time}ms"
        
        if [ $response_time -lt 100 ]; then
            log_info "✓ 响应时间优秀 (<100ms)"
        elif [ $response_time -lt 500 ]; then
            log_info "✓ 响应时间良好 (<500ms)"
        else
            log_warn "⚠ 响应时间较慢 (>500ms)"
        fi
    else
        log_warn "⚠ 无法测试API性能"
    fi
}

# 生成验证报告
generate_report() {
    log_info "生成验证报告..."
    
    report_file="integration_validation_report_$(date +%Y%m%d_%H%M%S).txt"
    
    cat > $report_file << EOF
电商系统集成验证报告
=====================

验证时间: $(date)
验证环境: $(uname -a)

1. 基础设施验证
----------------
$(docker-compose ps 2>/dev/null || echo "Docker Compose未运行")

2. 服务状态
-----------
$(docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep ecommerce || echo "无电商服务运行")

3. 网络检查
-----------
$(docker network ls | grep ecommerce || echo "无电商网络")

4. 性能指标
-----------
API响应时间: ${response_time}ms

5. 问题汇总
-----------
$(if [ $response_time -gt 500 ]; then echo "- API响应时间较慢"; fi)
$(if ! docker ps --format '{{.Names}}' | grep -q "ecommerce-backend"; then echo "- 后端服务未运行"; fi)
$(if ! docker ps --format '{{.Names}}' | grep -q "ecommerce-frontend"; then echo "- 前端服务未运行"; fi)

6. 验证结论
-----------
$(if [ $response_time -lt 100 ] && docker ps --format '{{.Names}}' | grep -q "ecommerce-backend"; then
    echo "✅ 系统集成验证通过"
else
    echo "⚠ 系统集成验证部分通过，需要优化"
fi)

EOF
    
    log_info "验证报告已生成: $report_file"
}

# 主函数
main() {
    echo ""
    log_info "开始系统集成验证..."
    echo ""
    
    # 检查依赖
    check_command docker
    check_command curl
    
    # 检查Docker
    check_docker
    
    # 检查各项服务
    check_infrastructure
    echo ""
    
    check_backend
    echo ""
    
    check_frontend
    echo ""
    
    check_network
    echo ""
    
    check_database
    echo ""
    
    check_redis
    echo ""
    
    check_rabbitmq
    echo ""
    
    # 性能测试
    performance_test
    echo ""
    
    # 生成报告
    generate_report
    
    echo ""
    log_info "验证完成！"
    echo "========================================="
}

# 执行主函数
main "$@"