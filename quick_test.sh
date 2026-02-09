#!/bin/bash

echo "=== 快速API连通性测试 ==="
echo "测试时间: $(date)"
echo ""

# 测试健康检查
echo "1. 测试健康检查接口..."
for i in {1..30}; do
    if curl -s -f "http://localhost:8080/health" > /dev/null 2>&1; then
        echo "✓ 后端服务已启动"
        health_response=$(curl -s "http://localhost:8080/health")
        echo "响应: $health_response"
        break
    fi
    echo "  尝试 $i/30: 等待后端服务启动..."
    sleep 2
done

echo ""
echo "2. 测试API基础功能..."
echo "2.1 测试404处理..."
curl -s "http://localhost:8080/api/v1/nonexistent"

echo ""
echo "2.2 测试认证API结构..."
echo "尝试访问需要认证的接口（应该返回401）..."
curl -s -w "HTTP状态码: %{http_code}\n" -o /dev/null "http://localhost:8080/api/v1/users/me"

echo ""
echo "=== 测试完成 ==="