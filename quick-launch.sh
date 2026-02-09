#!/bin/bash

# 快速启动脚本 - 分阶段启动电商系统

set -e

echo "🚀 电商系统快速启动"
echo "========================================"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 函数：打印带颜色的消息
print_status() {
    echo -e "${BLUE}[$(date '+%H:%M:%S')]${NC} $1"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 阶段1: 检查基础设施
print_status "阶段1: 检查基础设施..."
if docker-compose ps | grep -q "Up"; then
    print_success "基础设施已在运行"
else
    print_status "启动基础设施..."
    docker-compose up -d postgres redis rabbitmq
    sleep 10
fi

# 检查基础设施健康状态
print_status "检查基础设施健康状态..."
if docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then
    print_success "PostgreSQL: 健康"
else
    print_error "PostgreSQL: 不健康"
    exit 1
fi

if docker-compose exec -T redis redis-cli ping 2>/dev/null | grep -q PONG; then
    print_success "Redis: 健康"
else
    print_error "Redis: 不健康"
    exit 1
fi

# 阶段2: 构建后端
print_status "阶段2: 构建后端应用程序..."
cd backend

# 检查Go模块
if [ ! -f "go.mod" ]; then
    print_error "go.mod文件不存在"
    exit 1
fi

print_status "下载Go依赖..."
if go mod download; then
    print_success "Go依赖下载完成"
else
    print_warning "Go依赖下载失败，尝试继续构建..."
fi

print_status "构建后端二进制文件..."
if go build -o ../ecommerce-backend ./...; then
    print_success "后端构建成功"
    print_status "二进制文件: $(pwd)/../ecommerce-backend"
else
    print_error "后端构建失败"
    exit 1
fi

cd ..

# 阶段3: 启动后端
print_status "阶段3: 启动后端服务..."
export DATABASE_URL="host=localhost user=postgres password=postgres dbname=ecommerce port=5432 sslmode=disable"
export PORT="8080"

# 在后台启动后端
./ecommerce-backend &
BACKEND_PID=$!

print_status "等待后端启动..."
sleep 5

# 检查后端健康
print_status "检查后端健康状态..."
for i in {1..10}; do
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        print_success "后端健康检查通过"
        break
    fi
    if [ $i -eq 10 ]; then
        print_error "后端启动失败"
        kill $BACKEND_PID 2>/dev/null || true
        exit 1
    fi
    print_status "等待后端启动... ($i/10)"
    sleep 2
done

# 阶段4: 测试API
print_status "阶段4: 测试API端点..."
echo ""
echo "📡 API端点测试:"
echo "----------------------------------------"

# 测试根端点
if curl -s http://localhost:8080/ | grep -q "电商系统API"; then
    print_success "根端点: 正常"
else
    print_warning "根端点: 响应异常"
fi

# 测试健康检查
HEALTH_RESPONSE=$(curl -s http://localhost:8080/health)
if echo "$HEALTH_RESPONSE" | grep -q "healthy"; then
    print_success "健康检查: 正常"
    echo "   响应: $(echo "$HEALTH_RESPONSE" | jq -r '.status' 2>/dev/null || echo "$HEALTH_RESPONSE")"
else
    print_warning "健康检查: 响应异常"
fi

# 测试商品API
if curl -s http://localhost:8080/api/v1/products | grep -q "商品列表"; then
    print_success "商品API: 正常"
else
    print_warning "商品API: 开发中"
fi

# 测试认证API
if curl -s http://localhost:8080/api/v1/auth/register | grep -q "注册功能"; then
    print_success "注册API: 正常"
else
    print_warning "注册API: 开发中"
fi

# 阶段5: 系统总结
echo ""
echo "========================================"
echo "🎉 系统启动完成！"
echo "========================================"
echo ""
echo "🌐 服务访问地址:"
echo "----------------------------------------"
echo "后端API:      http://localhost:8080"
echo "健康检查:     http://localhost:8080/health"
echo "API文档:      http://localhost:8080/"
echo "PostgreSQL:   localhost:5432"
echo "Redis:        localhost:6379"
echo "RabbitMQ管理: http://localhost:15672"
echo "              (用户: admin, 密码: admin)"
echo ""
echo "📊 系统状态:"
echo "----------------------------------------"
echo "基础设施:     ✅ 运行中"
echo "后端服务:     ✅ 运行中 (PID: $BACKEND_PID)"
echo "API端点:      🚧 开发中"
echo "前端界面:     🚧 待启动"
echo ""
echo "🔧 管理命令:"
echo "----------------------------------------"
echo "停止后端:     kill $BACKEND_PID"
echo "停止基础设施: docker-compose down"
echo "查看日志:     docker-compose logs -f"
echo ""
echo "🧪 测试命令:"
echo "----------------------------------------"
echo "健康检查:     curl http://localhost:8080/health"
echo "商品列表:     curl http://localhost:8080/api/v1/products"
echo "用户注册:     curl -X POST http://localhost:8080/api/v1/auth/register"
echo ""
echo "⚠️  注意事项:"
echo "----------------------------------------"
echo "1. 这是一个开发版本，功能正在实现中"
echo "2. 前端界面尚未集成"
echo "3. 数据库迁移需要手动运行"
echo "4. 按Ctrl+C停止所有服务"
echo ""
echo "🚀 下一步:"
echo "   1. 运行数据库迁移: docker-compose exec postgres psql -U postgres -d ecommerce -f /migrations/001_initial_schema.sql"
echo "   2. 启动前端开发服务器"
echo "   3. 集成CRMEB功能模块"
echo ""
print_status "系统正在运行中..."
echo "按Ctrl+C停止所有服务"

# 捕获Ctrl+C信号
trap 'echo ""; print_status "正在停止服务..."; kill $BACKEND_PID 2>/dev/null || true; docker-compose down; print_success "服务已停止"; exit 0' INT TERM

# 保持脚本运行
wait $BACKEND_PID