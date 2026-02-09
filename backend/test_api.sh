#!/bin/bash

# API测试脚本
BASE_URL="http://localhost:8080/api/v1"

echo "=== 电商系统API测试 ==="
echo ""

# 1. 测试健康检查
echo "1. 测试健康检查..."
curl -s -X GET "$BASE_URL/../health" | jq .
echo ""

# 2. 测试用户注册
echo "2. 测试用户注册..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser_'$(date +%s)'",
    "email": "test_'$(date +%s)'@example.com",
    "password": "password123"
  }')

echo "$REGISTER_RESPONSE" | jq .
ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.access_token')
echo ""

# 3. 测试用户登录
echo "3. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }')

echo "$LOGIN_RESPONSE" | jq .
if [ "$ACCESS_TOKEN" = "null" ]; then
  ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')
fi
echo ""

# 4. 测试获取分类列表
echo "4. 测试获取分类列表..."
curl -s -X GET "$BASE_URL/categories" | jq .
echo ""

# 5. 测试获取商品列表
echo "5. 测试获取商品列表..."
curl -s -X GET "$BASE_URL/products?page=1&page_size=10" | jq .
echo ""

# 6. 测试获取用户资料（需要认证）
echo "6. 测试获取用户资料..."
if [ "$ACCESS_TOKEN" != "null" ]; then
  curl -s -X GET "$BASE_URL/auth/profile" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
else
  echo "未获取到访问令牌，跳过认证测试"
fi
echo ""

# 7. 测试获取购物车（需要认证）
echo "7. 测试获取购物车..."
if [ "$ACCESS_TOKEN" != "null" ]; then
  curl -s -X GET "$BASE_URL/cart" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
fi
echo ""

# 8. 测试获取订单列表（需要认证）
echo "8. 测试获取订单列表..."
if [ "$ACCESS_TOKEN" != "null" ]; then
  curl -s -X GET "$BASE_URL/orders" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
fi
echo ""

echo "=== API测试完成 ==="