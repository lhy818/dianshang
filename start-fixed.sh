#!/bin/bash

# 电商系统修复版启动脚本
# 解决Docker构建问题，优化启动流程

set -e  # 遇到错误立即退出

echo "🚀 启动电商系统（修复版）..."
echo "========================================"

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker未安装，请先安装Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose未安装，请先安装Docker Compose"
    exit 1
fi

echo "✅ Docker和Docker Compose已安装"

# 创建必要的目录
echo "📁 创建必要的目录..."
mkdir -p ./backend/logs
mkdir -p ./nginx/ssl
mkdir -p ./frontend/dist

# 停止并删除旧容器
echo "🧹 清理旧容器..."
docker-compose -f docker-compose.yml down 2>/dev/null || true
docker-compose -f docker-compose.prod.yml down 2>/dev/null || true

# 删除旧的镜像
echo "🗑️  删除旧的Docker镜像..."
docker rmi ecommerce-backend ecommerce-frontend 2>/dev/null || true

# 先单独构建前端（解决npm ci问题）
echo "🔨 构建前端应用..."
cd frontend
echo "📦 安装前端依赖..."
npm install --no-audit --no-fund
echo "🏗️  构建前端..."
npm run build
cd ..

# 创建.env文件（如果不存在）
if [ ! -f .env ]; then
    echo "📝 创建.env配置文件..."
    cat > .env << EOF
# 数据库配置
DB_PASSWORD=ecommerce_password_123

# Redis配置
REDIS_PASSWORD=redis_password_123

# RabbitMQ配置
RABBITMQ_USER=admin
RABBITMQ_PASSWORD=admin_password_123

# JWT密钥（生产环境请修改）
JWT_SECRET=your_super_secret_jwt_key_change_this_in_production_123
EOF
    echo "✅ .env文件已创建"
fi

# 启动基础设施服务
echo "🐘 启动PostgreSQL数据库..."
docker-compose up -d postgres

echo "🔴 启动Redis缓存..."
docker-compose up -d redis

echo "🐇 启动RabbitMQ消息队列..."
docker-compose up -d rabbitmq

# 等待基础设施服务就绪
echo "⏳ 等待基础设施服务就绪..."
sleep 10

# 检查服务健康状态
echo "🏥 检查服务健康状态..."
if docker-compose ps | grep -q "unhealthy"; then
    echo "⚠️  有些服务不健康，检查日志..."
    docker-compose logs --tail=20
else
    echo "✅ 所有基础设施服务健康"
fi

# 构建并启动后端
echo "⚙️  构建后端应用..."
cd backend
echo "📦 下载Go依赖..."
go mod download
echo "🏗️  构建Go应用..."
CGO_ENABLED=0 GOOS=linux go build -o main ./cmd
cd ..

# 启动后端服务
echo "🚀 启动后端API服务..."
docker-compose up -d backend

# 等待后端启动
echo "⏳ 等待后端服务启动..."
sleep 15

# 检查后端健康
echo "🏥 检查后端健康状态..."
if curl -f http://localhost:8080/health 2>/dev/null; then
    echo "✅ 后端API服务运行正常"
else
    echo "⚠️  后端API服务可能有问题，检查日志..."
    docker-compose logs backend --tail=20
fi

# 启动前端服务（使用开发模式，避免构建问题）
echo "🎨 启动前端开发服务器..."
cd frontend
echo "🚀 启动Vite开发服务器..."
npm run dev &
FRONTEND_PID=$!
cd ..

echo "⏳ 等待前端服务启动..."
sleep 10

# 检查前端是否运行
if curl -f http://localhost:3000 2>/dev/null; then
    echo "✅ 前端服务运行正常"
else
    echo "⚠️  前端服务可能有问题，尝试其他端口..."
    # 尝试3001端口
    if curl -f http://localhost:3001 2>/dev/null; then
        echo "✅ 前端服务在端口3001运行正常"
    else
        echo "❌ 前端服务启动失败，检查日志..."
        ps aux | grep -i vite || true
    fi
fi

echo ""
echo "========================================"
echo "🎉 电商系统启动完成！"
echo ""
echo "📊 服务状态："
echo "----------------------------------------"
docker-compose ps
echo ""
echo "🌐 访问地址："
echo "----------------------------------------"
echo "前端界面： http://localhost:3000"
echo "后端API：  http://localhost:8080"
echo "API文档：  http://localhost:8080/swagger/index.html"
echo "数据库：   localhost:5432 (用户: ecommerce_user)"
echo "Redis：    localhost:6379"
echo "RabbitMQ管理： http://localhost:15672 (用户: admin)"
echo ""
echo "🔧 管理命令："
echo "----------------------------------------"
echo "查看日志：    docker-compose logs -f"
echo "停止服务：    docker-compose down"
echo "重启服务：    docker-compose restart"
echo "清理数据：    docker-compose down -v"
echo ""
echo "📋 测试命令："
echo "----------------------------------------"
echo "健康检查：    curl http://localhost:8080/health"
echo "API测试：     curl http://localhost:8080/api/v1/products"
echo "前端测试：    curl http://localhost:3000"
echo ""
echo "⚠️  注意事项："
echo "----------------------------------------"
echo "1. 首次启动可能需要几分钟时间"
echo "2. 确保端口3000、8080、5432、6379、5672、15672未被占用"
echo "3. 生产环境请修改.env文件中的密码和密钥"
echo "4. 查看详细日志：docker-compose logs --tail=50"
echo ""
echo "🚀 开始使用电商系统吧！"

# 保存进程ID
echo $FRONTEND_PID > .frontend.pid

# 设置退出时清理
trap "echo '🛑 停止服务...'; kill $FRONTEND_PID 2>/dev/null || true; docker-compose down; echo '✅ 服务已停止'" EXIT

# 保持脚本运行
echo ""
echo "📝 按Ctrl+C停止所有服务..."
wait $FRONTEND_PID