# 电商系统后端API

基于Go + Gin + PostgreSQL + Redis开发的电商系统后端API。

## 功能特性

- ✅ 用户认证（注册、登录、登出、JWT令牌）
- ✅ 商品管理（分类、商品列表、商品详情、SKU管理）
- ✅ 购物车管理（添加、更新、删除商品）
- ✅ 订单管理（创建订单、订单列表、订单详情、取消订单）
- ✅ 统一响应格式
- ✅ 输入参数验证
- ✅ 数据库事务处理
- ✅ 并发安全和性能优化

## 技术栈

- **后端框架**: Gin
- **数据库**: PostgreSQL
- **ORM**: GORM
- **缓存**: Redis
- **认证**: JWT
- **配置管理**: Viper
- **密码加密**: bcrypt

## 项目结构

```
ecommerce-backend/
├── cmd/
│   └── main.go              # 应用入口
├── config/
│   └── config.yaml          # 配置文件
├── internal/
│   ├── api/                 # API处理器
│   │   ├── auth_handler.go
│   │   ├── product_handler.go
│   │   ├── cart_handler.go
│   │   └── order_handler.go
│   ├── domain/              # 领域模型
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── cart.go
│   │   └── order.go
│   ├── repository/          # 数据仓库
│   │   ├── user_repository.go
│   │   ├── product_repository.go
│   │   ├── cart_repository.go
│   │   └── order_repository.go
│   └── service/             # 业务服务
│       ├── user_service.go
│       ├── product_service.go
│       ├── cart_service.go
│       └── order_service.go
├── migrations/              # 数据库迁移
│   └── 001_initial_schema.sql
├── go.mod                   # Go模块文件
├── go.sum                   # 依赖校验
└── README.md               # 项目说明
```

## API接口

### 认证相关
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户登出
- `GET /api/v1/auth/profile` - 获取用户资料

### 商品相关
- `GET /api/v1/categories` - 分类列表
- `GET /api/v1/products` - 商品列表（支持搜索、筛选、排序、分页）
- `GET /api/v1/products/{id}` - 商品详情
- `GET /api/v1/products/{id}/skus` - 商品SKU列表

### 购物车相关
- `GET /api/v1/cart` - 获取购物车
- `POST /api/v1/cart/items` - 添加商品到购物车
- `PUT /api/v1/cart/items/{id}` - 更新购物车商品
- `DELETE /api/v1/cart/items/{id}` - 删除购物车商品

### 订单相关
- `GET /api/v1/orders` - 订单列表
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders/{id}` - 订单详情
- `PUT /api/v1/orders/{id}/cancel` - 取消订单

## 快速开始

### 1. 环境要求

- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### 2. 安装依赖

```bash
cd ecommerce-project/backend
go mod download
```

### 3. 配置数据库

创建PostgreSQL数据库：
```sql
CREATE DATABASE ecommerce;
```

运行数据库迁移：
```bash
psql -U postgres -d ecommerce -f migrations/001_initial_schema.sql
```

### 4. 配置应用

编辑 `config/config.yaml` 文件，更新数据库连接信息：
```yaml
database:
  postgres:
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "your_password"
    dbname: "ecommerce"
```

### 5. 运行应用

```bash
go run cmd/main.go
```

应用将在 `http://localhost:8080` 启动。

### 6. 测试API

运行测试脚本：
```bash
chmod +x test_api.sh
./test_api.sh
```

## 开发指南

### 添加新的API

1. 在 `internal/domain/` 中添加领域模型
2. 在 `internal/repository/` 中添加数据仓库
3. 在 `internal/service/` 中添加业务服务
4. 在 `internal/api/` 中添加API处理器
5. 在 `cmd/main.go` 中注册路由

### 数据库迁移

开发环境使用GORM自动迁移：
```go
db.AutoMigrate(&YourModel{})
```

生产环境使用SQL迁移文件：
```bash
psql -U postgres -d ecommerce -f migrations/002_new_feature.sql
```

### 测试

运行单元测试：
```bash
go test ./...
```

## 部署

### Docker部署

创建Dockerfile：
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config/config.yaml ./config/
EXPOSE 8080
CMD ["./main"]
```

构建并运行：
```bash
docker build -t ecommerce-backend .
docker run -p 8080:8080 ecommerce-backend
```

### Kubernetes部署

创建deployment.yaml：
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ecommerce-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ecommerce-backend
  template:
    metadata:
      labels:
        app: ecommerce-backend
    spec:
      containers:
      - name: ecommerce-backend
        image: ecommerce-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: url
```

## 性能优化

1. **数据库索引**: 所有查询字段都建立了索引
2. **连接池**: 配置了数据库连接池
3. **缓存**: 热点数据使用Redis缓存
4. **异步处理**: 非核心业务异步化
5. **CDN**: 静态资源使用CDN加速

## 安全考虑

1. **认证授权**: JWT令牌认证，RBAC权限控制
2. **数据加密**: 敏感数据加密存储
3. **输入验证**: 所有输入参数都经过验证
4. **SQL注入防护**: 使用预编译语句
5. **XSS防护**: 输入过滤，输出编码
6. **CSRF防护**: Token验证
7. **限流**: 接口请求限流

## 监控和日志

1. **应用监控**: Prometheus + Grafana
2. **业务监控**: 关键业务指标监控
3. **日志系统**: 结构化日志，分布式追踪
4. **错误追踪**: Sentry或类似服务

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交代码变更
4. 编写测试用例
5. 提交Pull Request

## 许可证

MIT License