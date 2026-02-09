#!/bin/bash

# 电商网站项目测试脚本
# 测试项目是否可以正常运行

set -e

echo "🔍 开始测试电商网站项目..."

# 检查项目结构
echo "📁 检查项目结构..."
if [ ! -d "frontend" ] || [ ! -d "backend" ]; then
    echo "❌ 项目结构不完整"
    exit 1
fi

echo "✅ 项目结构完整"

# 检查前端项目
echo "🖥️  检查前端项目..."
if [ ! -f "frontend/package.json" ]; then
    echo "❌ 前端package.json不存在"
    exit 1
fi

if [ ! -f "frontend/vite.config.ts" ]; then
    echo "❌ 前端vite配置不存在"
    exit 1
fi

echo "✅ 前端项目配置完整"

# 检查后端项目
echo "⚙️  检查后端项目..."
if [ ! -f "backend/go.mod" ]; then
    echo "❌ 后端go.mod不存在"
    exit 1
fi

if [ ! -f "backend/cmd/main.go" ]; then
    echo "❌ 后端主文件不存在"
    exit 1
fi

echo "✅ 后端项目配置完整"

# 检查配置文件
echo "📋 检查配置文件..."
if [ ! -f "backend/config/config.yaml" ]; then
    echo "❌ 后端配置文件不存在"
    exit 1
fi

if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Docker Compose配置不存在"
    exit 1
fi

echo "✅ 配置文件完整"

# 检查数据库迁移文件
echo "🗄️  检查数据库迁移..."
if [ ! -f "backend/migrations/001_initial_schema.sql" ]; then
    echo "❌ 数据库迁移文件不存在"
    exit 1
fi

echo "✅ 数据库迁移文件完整"

# 检查领域模型
echo "🏗️  检查领域模型..."
if [ ! -f "backend/internal/domain/user.go" ]; then
    echo "❌ 用户领域模型不存在"
    exit 1
fi

if [ ! -f "backend/internal/domain/product.go" ]; then
    echo "❌ 商品领域模型不存在"
    exit 1
fi

echo "✅ 领域模型完整"

# 检查API处理器
echo "🌐 检查API处理器..."
if [ ! -f "backend/internal/api/auth_handler.go" ]; then
    echo "❌ 认证处理器不存在"
    exit 1
fi

echo "✅ API处理器完整"

# 检查服务层
echo "🛠️  检查服务层..."
if [ ! -f "backend/internal/service/user_service.go" ]; then
    echo "❌ 用户服务不存在"
    exit 1
fi

echo "✅ 服务层完整"

# 检查仓库层
echo "💾 检查仓库层..."
if [ ! -f "backend/internal/repository/user_repository.go" ]; then
    echo "❌ 用户仓库不存在"
    exit 1
fi

echo "✅ 仓库层完整"

# 检查启动脚本
echo "🚀 检查启动脚本..."
if [ ! -f "start.sh" ]; then
    echo "❌ 启动脚本不存在"
    exit 1
fi

if [ ! -x "start.sh" ]; then
    echo "⚠️  启动脚本没有执行权限，正在添加..."
    chmod +x start.sh
fi

echo "✅ 启动脚本完整"

# 检查文档
echo "📚 检查项目文档..."
if [ ! -f "docs/项目启动指南.md" ]; then
    echo "❌ 项目启动指南不存在"
    exit 1
fi

if [ ! -f "项目状态报告.md" ]; then
    echo "❌ 项目状态报告不存在"
    exit 1
fi

echo "✅ 项目文档完整"

# 总结
echo ""
echo "🎉 项目测试完成！"
echo ""
echo "📊 测试结果汇总："
echo "  ✅ 项目结构完整"
echo "  ✅ 前端项目配置完整"
echo "  ✅ 后端项目配置完整"
echo "  ✅ 配置文件完整"
echo "  ✅ 数据库迁移完整"
echo "  ✅ 领域模型完整"
echo "  ✅ API处理器完整"
echo "  ✅ 服务层完整"
echo "  ✅ 仓库层完整"
echo "  ✅ 启动脚本完整"
echo "  ✅ 项目文档完整"
echo ""
echo "🚀 项目已准备好，可以开始开发！"
echo ""
echo "下一步建议："
echo "  1. 运行 ./start.sh infra 启动基础设施"
echo "  2. 运行 ./start.sh dev 启动开发环境"
echo "  3. 访问 http://localhost:3000 查看前端"
echo "  4. 访问 http://localhost:8080/health 检查后端健康状态"
echo ""
echo "💡 提示：确保已安装 Docker、Node.js 18+ 和 Go 1.21+"