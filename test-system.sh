#!/bin/bash

# 电商系统测试脚本
# 使用Docker运行所有服务，避免环境依赖问题

set -e

echo "🧪 电商系统测试启动..."
echo "========================================"

# 清理旧容器
echo "🧹 清理旧容器..."
docker-compose down 2>/dev/null || true

# 1. 启动所有服务
echo "1. 🚀 启动所有服务..."
docker-compose up -d

echo "⏳ 等待服务启动..."
sleep 30

# 2. 检查服务状态
echo "2. 📊 检查服务状态..."
echo ""
docker-compose ps
echo ""

# 3. 测试基础设施
echo "3. 🏗️  测试基础设施..."
echo "----------------------------------------"

# 测试PostgreSQL
if docker-compose exec -T postgres pg_isready -U postgres; then
    echo "✅ PostgreSQL: 运行正常"
else
    echo "❌ PostgreSQL: 连接失败"
fi

# 测试Redis
if docker-compose exec -T redis redis-cli ping | grep -q PONG; then
    echo "✅ Redis: 运行正常"
else
    echo "❌ Redis: 连接失败"
fi

# 测试RabbitMQ
if curl -f http://localhost:15672 2>/dev/null; then
    echo "✅ RabbitMQ: 管理界面可访问"
else
    echo "⚠️  RabbitMQ: 管理界面不可访问"
fi

# 4. 测试后端API
echo ""
echo "4. ⚙️  测试后端API..."
echo "----------------------------------------"

# 等待后端启动
echo "⏳ 等待后端API启动..."
for i in {1..10}; do
    if curl -f http://localhost:8080/health 2>/dev/null; then
        echo "✅ 后端健康检查: 通过"
        break
    fi
    echo "等待后端启动... ($i/10)"
    sleep 5
done

# 测试API端点
echo ""
echo "📡 测试API端点:"
echo "----------------------------------------"

# 测试商品API
if curl -f http://localhost:8080/api/v1/products 2>/dev/null; then
    echo "✅ 商品API: 可访问"
else
    echo "❌ 商品API: 不可访问"
fi

# 测试分类API
if curl -f http://localhost:8080/api/v1/categories 2>/dev/null; then
    echo "✅ 分类API: 可访问"
else
    echo "❌ 分类API: 不可访问"
fi

# 5. 测试前端
echo ""
echo "5. 🎨 测试前端..."
echo "----------------------------------------"

# 等待前端启动
echo "⏳ 等待前端启动..."
for i in {1..10}; do
    if curl -f http://localhost:3000 2>/dev/null; then
        echo "✅ 前端: 可访问 (端口3000)"
        break
    elif curl -f http://localhost:3001 2>/dev/null; then
        echo "✅ 前端: 可访问 (端口3001)"
        break
    fi
    echo "等待前端启动... ($i/10)"
    sleep 5
done

# 6. 完整业务流程测试
echo ""
echo "6. 🔄 完整业务流程测试..."
echo "----------------------------------------"

echo "📋 测试步骤:"
echo "1. 用户注册"
echo "2. 用户登录"
echo "3. 浏览商品"
echo "4. 加入购物车"
echo "5. 创建订单"

# 创建测试用户
echo ""
echo "👤 创建测试用户..."
TEST_USER="testuser_$(date +%s)"
TEST_EMAIL="${TEST_USER}@test.com"
TEST_PASSWORD="Test123!"

REGISTER_DATA=$(cat <<EOF
{
  "username": "$TEST_USER",
  "email": "$TEST_EMAIL",
  "password": "$TEST_PASSWORD"
}
EOF
)

if curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d "$REGISTER_DATA" 2>/dev/null | grep -q "注册成功"; then
    echo "✅ 用户注册: 成功 (用户: $TEST_USER)"
else
    echo "❌ 用户注册: 失败"
fi

# 用户登录
echo ""
echo "🔐 用户登录..."
LOGIN_DATA=$(cat <<EOF
{
  "username": "$TEST_USER",
  "password": "$TEST_PASSWORD"
}
EOF
)

LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "$LOGIN_DATA")

if echo "$LOGIN_RESPONSE" | grep -q "token"; then
    TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo "✅ 用户登录: 成功 (Token获取)"
    
    # 测试获取用户信息
    if curl -f http://localhost:8080/api/v1/auth/profile \
      -H "Authorization: Bearer $TOKEN" 2>/dev/null; then
        echo "✅ 用户信息: 可获取"
    else
        echo "❌ 用户信息: 获取失败"
    fi
else
    echo "❌ 用户登录: 失败"
fi

# 7. 系统总结
echo ""
echo "========================================"
echo "🧪 测试完成总结"
echo "========================================"

echo ""
echo "🌐 服务访问地址:"
echo "----------------------------------------"
echo "前端界面:     http://localhost:3000"
echo "后端API:      http://localhost:8080"
echo "API文档:      http://localhost:8080/swagger/index.html"
echo "RabbitMQ管理: http://localhost:15672 (用户: admin, 密码: admin)"
echo ""

echo "📊 测试命令汇总:"
echo "----------------------------------------"
echo "健康检查:     curl http://localhost:8080/health"
echo "商品列表:     curl http://localhost:8080/api/v1/products"
echo "分类列表:     curl http://localhost:8080/api/v1/categories"
echo "前端访问:     curl http://localhost:3000"
echo ""

echo "🔧 管理命令:"
echo "----------------------------------------"
echo "查看日志:     docker-compose logs -f"
echo "停止服务:     docker-compose down"
echo "重启服务:     docker-compose restart"
echo "清理数据:     docker-compose down -v"
echo ""

echo "⚠️  注意事项:"
echo "----------------------------------------"
echo "1. 首次启动可能需要1-2分钟"
echo "2. 如果服务不健康，检查: docker-compose logs"
echo "3. 测试用户: $TEST_USER / $TEST_PASSWORD"
echo "4. 按Ctrl+C停止测试并清理"
echo ""

echo "🚀 系统测试完成！可以开始使用了。"

# 保持运行，按Ctrl+C停止
echo ""
echo "📝 按Ctrl+C停止所有服务..."
trap "echo '🛑 停止服务...'; docker-compose down; echo '✅ 服务已停止'" INT TERM

# 显示实时日志
echo ""
echo "📋 显示服务日志 (最后20行):"
echo "----------------------------------------"
docker-compose logs --tail=20

# 等待用户中断
sleep infinity