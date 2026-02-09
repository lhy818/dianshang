#!/bin/bash
cd backend

# 检查是否在容器内
if [ -f /.dockerenv ]; then
    echo "在容器内，直接启动后端..."
    go run cmd/main.go
else
    echo "在主机上，使用Docker启动后端..."
    docker run --rm -d \
        --name ecommerce-backend \
        --network ecommerce-project_default \
        -p 8080:8080 \
        -e APP_ENV=development \
        -e DB_HOST=postgres \
        -e DB_PORT=5432 \
        -e REDIS_HOST=redis \
        -e RABBITMQ_HOST=rabbitmq \
        -v $(pwd):/app \
        -w /app \
        golang:1.21-alpine \
        sh -c "apk add --no-cache git && go mod download && go run cmd/main.go"
fi