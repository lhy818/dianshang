#!/bin/bash

# 电商系统快速启动脚本
# 避免Docker构建问题，分别启动服务

set -e

echo "🚀 电商系统快速启动..."
echo "========================================"

# 1. 启动基础设施
echo "1. 🏗️  启动基础设施..."
docker-compose up -d postgres redis rabbitmq

echo "⏳ 等待基础设施就绪..."
sleep 15

# 2. 启动后端（直接运行Go）
echo "2. ⚙️  启动后端API..."
cd backend
echo "📦 下载Go依赖..."
go mod download
echo "🚀 启动Go服务..."
go run cmd/main.go &
BACKEND_PID=$!
cd ..

echo "⏳ 等待后端启动..."
sleep 10

# 检查后端
if curl -f http://localhost:8080/health 2>/dev/null; then
    echo "✅ 后端API运行正常 (PID: $BACKEND_PID)"
else
    echo "❌ 后端启动失败"
    exit 1
fi

# 3. 启动前端（直接运行Vite）
echo "3. 🎨 启动前端..."
cd frontend
echo "📦 安装前端依赖..."
npm install --no-audit --no-fund
echo "🚀 启动Vite开发服务器..."
npm run dev &
FRONTEND_PID=$!
cd ..

echo "⏳ 等待前端启动..."
sleep 10

# 检查前端
if curl -f http://localhost:3000 2>/dev/null; then
    echo "✅ 前端运行正常 (PID: $FRONTEND_PID)"
else
    echo "⚠️  前端可能在3001端口..."
    if curl -f http://localhost:3001 2>/dev/null; then
        echo "✅ 前端在3001端口运行正常"
    else
        echo "❌ 前端启动失败"
        exit 1
    fi
fi

echo ""
echo "========================================"
echo "🎉 电商系统启动成功！"
echo ""
echo "🌐 访问地址："
echo "----------------------------------------"
echo "前端： http://localhost:3000"
echo "后端： http://localhost:8080"
echo ""
echo "📊 测试命令："
echo "----------------------------------------"
echo "curl http://localhost:8080/health"
echo "curl http://localhost:8080/api/v1/products"
echo ""
echo "🛑 按Ctrl+C停止所有服务"

# 保存进程ID
echo $BACKEND_PID > .backend.pid
echo $FRONTEND_PID > .frontend.pid

# 清理函数
cleanup() {
    echo ""
    echo "🛑 停止服务..."
    kill $BACKEND_PID 2>/dev/null || true
    kill $FRONTEND_PID 2>/dev/null || true
    docker-compose down
    rm -f .backend.pid .frontend.pid
    echo "✅ 服务已停止"
    exit 0
}

trap cleanup INT TERM

# 等待
wait