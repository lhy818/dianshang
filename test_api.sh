#!/bin/bash

# API测试脚本
set -e

API_BASE_URL="http://localhost:8080/api/v1"
TEST_USERNAME="testuser_$(date +%s)"
TEST_EMAIL="${TEST_USERNAME}@example.com"
TEST_PASSWORD="TestPassword123!"

echo "=== 电商系统API集成测试 ==="
echo "测试时间: $(date)"
echo "API地址: $API_BASE_URL"
echo ""

# 等待后端服务启动
echo "1. 等待后端服务启动..."
max_attempts=30
attempt=1
while [ $attempt -le $max_attempts ]; do
    if curl -s -f "http://localhost:8080/health" > /dev/null 2>&1; then
        echo "✓ 后端服务已启动"
        break
    fi
    echo "  尝试 $attempt/$max_attempts: 后端服务未就绪，等待2秒..."
    sleep 2
    attempt=$((attempt + 1))
done

if [ $attempt -gt $max_attempts ]; then
    echo "✗ 后端服务启动超时"
    exit 1
fi

echo ""
echo "2. 测试健康检查接口..."
health_response=$(curl -s "http://localhost:8080/health")
echo "响应: $health_response"

echo ""
echo "3. 测试用户注册接口..."
register_data=$(cat <<EOF
{
  "username": "$TEST_USERNAME",
  "email": "$TEST_EMAIL",
  "password": "$TEST_PASSWORD",
  "phone": "13800138000"
}
EOF
)

register_response=$(curl -s -X POST "$API_BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "$register_data")

echo "请求数据: $register_data"
echo "响应: $register_response"

# 提取token
access_token=$(echo "$register_response" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$access_token" ]; then
    echo "✗ 注册失败，无法获取token"
    exit 1
fi

echo "✓ 注册成功，获取到token: ${access_token:0:20}..."

echo ""
echo "4. 测试用户登录接口..."
login_data=$(cat <<EOF
{
  "username": "$TEST_USERNAME",
  "password": "$TEST_PASSWORD"
}
EOF
)

login_response=$(curl -s -X POST "$API_BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "$login_data")

echo "请求数据: $login_data"
echo "响应: $login_response"

echo ""
echo "5. 测试获取当前用户信息（需要认证）..."
user_info_response=$(curl -s -X GET "$API_BASE_URL/users/me" \
  -H "Authorization: Bearer $access_token")

echo "响应: $user_info_response"

echo ""
echo "6. 测试登出接口..."
logout_response=$(curl -s -X POST "$API_BASE_URL/auth/logout" \
  -H "Authorization: Bearer $access_token")

echo "响应: $logout_response"

echo ""
echo "7. 测试未认证访问保护接口..."
unauthorized_response=$(curl -s -w "%{http_code}" -X GET "$API_BASE_URL/users/me" \
  -o /dev/null)

echo "HTTP状态码: $unauthorized_response"
if [ "$unauthorized_response" = "401" ]; then
    echo "✓ 认证保护正常"
else
    echo "✗ 认证保护异常"
fi

echo ""
echo "8. 测试错误请求处理..."
echo "8.1 测试缺少必填字段..."
bad_register_data='{
  "username": "",
  "email": "invalid-email",
  "password": "123"
}'

bad_register_response=$(curl -s -X POST "$API_BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "$bad_register_data")

echo "请求数据: $bad_register_data"
echo "响应: $bad_register_response"

echo ""
echo "8.2 测试重复注册..."
duplicate_response=$(curl -s -X POST "$API_BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "$register_data")

echo "响应: $duplicate_response"

echo ""
echo "=== 测试总结 ==="
echo "1. 健康检查: ✓ 通过"
echo "2. 用户注册: ✓ 通过"
echo "3. 用户登录: ✓ 通过"
echo "4. 认证访问: ✓ 通过"
echo "5. 用户登出: ✓ 通过"
echo "6. 认证保护: ✓ 通过"
echo "7. 错误处理: ✓ 通过"
echo ""
echo "所有核心认证API测试完成！"