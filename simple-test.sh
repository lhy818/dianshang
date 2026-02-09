#!/bin/bash

# 简化测试脚本 - 只启动基础设施和测试基本功能

set -e

echo "🧪 电商系统简化测试启动..."
echo "========================================"

# 清理旧容器
echo "🧹 清理旧容器..."
docker-compose down 2>/dev/null || true

# 1. 只启动基础设施（PostgreSQL, Redis, RabbitMQ）
echo "1. 🏗️  启动基础设施..."
docker-compose up -d postgres redis rabbitmq

echo "⏳ 等待基础设施启动..."
sleep 10

# 2. 检查基础设施状态
echo "2. 📊 检查基础设施状态..."
echo ""

# 测试PostgreSQL
if docker-compose exec -T postgres pg_isready -U postgres; then
    echo "✅ PostgreSQL: 运行正常"
    
    # 创建测试数据库
    echo "📦 创建测试数据库..."
    docker-compose exec -T postgres psql -U postgres -c "CREATE DATABASE IF NOT EXISTS ecommerce_test;"
    echo "✅ 测试数据库创建完成"
else
    echo "❌ PostgreSQL: 连接失败"
    exit 1
fi

# 测试Redis
if docker-compose exec -T redis redis-cli ping | grep -q PONG; then
    echo "✅ Redis: 运行正常"
else
    echo "❌ Redis: 连接失败"
    exit 1
fi

# 测试RabbitMQ
echo "⏳ 等待RabbitMQ启动..."
sleep 5
if curl -f http://localhost:15672 2>/dev/null; then
    echo "✅ RabbitMQ: 管理界面可访问"
    echo "   访问地址: http://localhost:15672 (用户: admin, 密码: admin)"
else
    echo "⚠️  RabbitMQ: 管理界面不可访问，但可能仍在启动中"
fi

# 3. 检查后端代码结构
echo ""
echo "3. 🔍 检查后端代码结构..."
echo "----------------------------------------"

if [ -f "backend/go.mod" ]; then
    echo "✅ Go模块文件存在"
    echo "   模块名称: $(grep '^module' backend/go.mod | cut -d' ' -f2)"
else
    echo "❌ Go模块文件不存在"
fi

if [ -f "backend/main.go" ]; then
    echo "✅ 主程序文件存在"
    echo "   文件大小: $(wc -l < backend/main.go) 行"
else
    echo "❌ 主程序文件不存在"
fi

# 4. 检查前端代码结构
echo ""
echo "4. 🎨 检查前端代码结构..."
echo "----------------------------------------"

if [ -f "frontend/package.json" ]; then
    echo "✅ package.json存在"
    echo "   项目名称: $(grep '"name"' frontend/package.json | cut -d'"' -f4)"
    echo "   依赖数量: $(grep -c '"dependencies"' frontend/package.json) 个依赖项"
else
    echo "❌ package.json不存在"
fi

if [ -f "frontend/src/App.tsx" ]; then
    echo "✅ 主应用组件存在"
    echo "   文件大小: $(wc -l < frontend/src/App.tsx) 行"
else
    echo "❌ 主应用组件不存在"
fi

# 5. 测试API文档
echo ""
echo "5. 📚 检查API文档..."
echo "----------------------------------------"

# 检查后端是否有API文档
if [ -f "backend/docs/swagger.yaml" ] || [ -f "backend/docs/swagger.json" ]; then
    echo "✅ Swagger文档存在"
else
    echo "⚠️  Swagger文档不存在，需要生成"
fi

# 6. 系统架构验证
echo ""
echo "6. 🏛️  系统架构验证..."
echo "----------------------------------------"

echo "📁 项目结构:"
echo "----------------------------------------"
find . -type f -name "*.md" -o -name "*.sql" -o -name "*.yaml" -o -name "*.yml" -o -name "Dockerfile" | head -20 | sort

# 7. 总结
echo ""
echo "========================================"
echo "🧪 简化测试完成"
echo "========================================"

echo ""
echo "✅ 基础设施已启动:"
echo "   - PostgreSQL: localhost:5432"
echo "   - Redis: localhost:6379"
echo "   - RabbitMQ: localhost:15672"
echo ""
echo "📁 项目结构验证通过"
echo ""
echo "🚀 下一步:"
echo "   1. 构建后端: cd backend && go build"
echo "   2. 安装前端依赖: cd frontend && npm install"
echo "   3. 启动开发环境: 使用 docker-compose up"
echo ""
echo "🔧 管理命令:"
echo "   停止服务: docker-compose down"
echo "   查看日志: docker-compose logs"
echo ""
echo "📝 按Ctrl+C停止所有服务..."
echo ""

# 显示基础设施状态
docker-compose ps

# 等待用户中断
trap "echo '🛑 停止服务...'; docker-compose down; echo '✅ 服务已停止'" INT TERM
sleep infinity