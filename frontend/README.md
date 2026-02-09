# 电商系统前端

基于React + TypeScript + Vite + Ant Design的现代化电商系统前端。

## 技术栈

- **React 18** - 前端框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Ant Design 5** - UI组件库
- **React Query** - 数据获取和缓存
- **Zustand** - 状态管理
- **React Router 6** - 路由管理
- **Axios** - HTTP客户端

## 项目结构

```
src/
├── api/                    # API相关
│   ├── client.ts          # Axios客户端配置
│   └── services.ts        # API服务封装
├── components/            # 组件目录
│   ├── auth/              # 认证模块组件
│   │   ├── LoginForm.tsx  # 登录组件
│   │   ├── RegisterForm.tsx # 注册组件
│   │   └── UserProfile.tsx  # 用户信息组件
│   ├── products/          # 商品模块组件
│   │   ├── ProductList.tsx  # 商品列表组件
│   │   └── ProductDetail.tsx # 商品详情组件
│   ├── cart/              # 购物车模块组件
│   │   └── CartPage.tsx   # 购物车页面组件
│   └── orders/            # 订单模块组件
│       ├── OrderList.tsx  # 订单列表组件
│       └── OrderDetail.tsx # 订单详情组件
├── store/                 # 状态管理
│   ├── authStore.ts       # 认证状态
│   └── cartStore.ts       # 购物车状态
├── types/                 # TypeScript类型定义
│   └── index.ts           # 所有类型定义
├── providers/             # 上下文提供者
│   └── QueryProvider.tsx  # React Query提供者
├── examples/              # 组件示例
│   └── ComponentExamples.tsx # 组件使用示例
├── styles/                # 样式文件
│   └── index.css          # 全局样式
├── App.tsx               # 主应用组件
└── main.tsx             # 应用入口
```

## 核心功能模块

### 1. 用户认证模块
- **登录组件** (`LoginForm`): 用户登录，包含表单验证
- **注册组件** (`RegisterForm`): 新用户注册，密码强度验证
- **用户信息组件** (`UserProfile`): 显示用户个人信息和状态

### 2. 商品展示模块
- **商品列表组件** (`ProductList`): 支持分页、筛选、排序的商品列表
- **商品详情组件** (`ProductDetail`): 商品详细信息，包含图片轮播、规格参数

### 3. 购物车模块
- **购物车页面组件** (`CartPage`): 购物车管理，支持批量操作和结算

### 4. 订单模块
- **订单列表组件** (`OrderList`): 用户订单列表，支持搜索和筛选
- **订单详情组件** (`OrderDetail`): 订单详细信息，包含时间线和商品列表

## 状态管理

### 认证状态 (`authStore`)
```typescript
interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  
  // Actions
  login: (user: User, token: string) => void;
  logout: () => void;
  // ...
}
```

### 购物车状态 (`cartStore`)
```typescript
interface CartState {
  items: CartItem[];
  totalItems: number;
  totalAmount: number;
  
  // Actions
  addItem: (product: Product, quantity: number, skuId?: string) => void;
  updateQuantity: (itemId: string, quantity: number) => void;
  removeItem: (itemId: string) => void;
  clearCart: () => void;
  // ...
}
```

## API服务

### 认证API
```typescript
const authApi = {
  login: (data: LoginRequest) => Promise<{ user: User; token: string }>;
  register: (data: RegisterRequest) => Promise<{ user: User; token: string }>;
  logout: () => Promise<void>;
  getProfile: () => Promise<User>;
  // ...
};
```

### 商品API
```typescript
const productApi = {
  getProducts: (params?: ProductQueryParams) => Promise<PaginatedResponse<Product>>;
  getProduct: (id: string) => Promise<Product>;
  getCategories: () => Promise<Category[]>;
  // ...
};
```

### 购物车API
```typescript
const cartApi = {
  getCart: () => Promise<{ items: CartItem[]; total: number }>;
  addToCart: (data: { productId: string; skuId?: string; quantity: number }) => Promise<CartItem>;
  updateCartItem: (id: string, data: { quantity: number }) => Promise<CartItem>;
  // ...
};
```

### 订单API
```typescript
const orderApi = {
  getOrders: (params?: { page?: number; pageSize?: number }) => Promise<PaginatedResponse<Order>>;
  createOrder: (data: CreateOrderRequest) => Promise<Order>;
  getOrder: (id: string) => Promise<Order>;
  // ...
};
```

## 开发指南

### 环境配置
1. 复制环境变量文件：
   ```bash
   cp .env.example .env
   ```

2. 配置环境变量：
   ```env
   VITE_API_BASE_URL=http://localhost:8080/api/v1
   ```

### 安装依赖
```bash
npm install
```

### 启动开发服务器
```bash
npm run dev
```

### 构建生产版本
```bash
npm run build
```

### 代码检查
```bash
npm run lint
```

### 代码格式化
```bash
npm run format
```

### 运行测试
```bash
npm run test
```

## 组件使用示例

### 商品列表组件
```tsx
import { ProductList } from './components/products/ProductList';

function App() {
  return (
    <ProductList 
      categoryId="123"  // 可选：分类ID
      showFilters={true} // 可选：是否显示筛选器
    />
  );
}
```

### 购物车组件
```tsx
import { CartPage } from './components/cart/CartPage';
import { useCartStore } from './store/cartStore';

function App() {
  const { addItem } = useCartStore();
  
  // 添加商品到购物车
  const handleAddToCart = (product: Product) => {
    addItem(product, 1);
  };
  
  return <CartPage />;
}
```

### 认证组件
```tsx
import { LoginForm } from './components/auth/LoginForm';
import { RegisterForm } from './components/auth/RegisterForm';
import { useAuthStore } from './store/authStore';

function App() {
  const { isAuthenticated, user } = useAuthStore();
  
  if (!isAuthenticated) {
    return <LoginForm />;
  }
  
  return <div>欢迎回来，{user?.username}</div>;
}
```

## 路由配置

```tsx
<Routes>
  <Route path="/login" element={<LoginForm />} />
  <Route path="/register" element={<RegisterForm />} />
  <Route path="/" element={
    <ProtectedRoute>
      <MainLayout>
        <ProductList />
      </MainLayout>
    </ProtectedRoute>
  } />
  <Route path="/products" element={
    <ProtectedRoute>
      <MainLayout>
        <ProductList />
      </MainLayout>
    </ProtectedRoute>
  } />
  <Route path="/products/:id" element={
    <ProtectedRoute>
      <MainLayout>
        <ProductDetail />
      </MainLayout>
    </ProtectedRoute>
  } />
  <Route path="/cart" element={
    <ProtectedRoute>
      <MainLayout>
        <CartPage />
      </MainLayout>
    </ProtectedRoute>
  } />
  <Route path="/orders" element={
    <ProtectedRoute>
      <MainLayout>
        <OrderList />
      </MainLayout>
    </ProtectedRoute>
  } />
  <Route path="/orders/:id" element={
    <ProtectedRoute>
      <MainLayout>
        <OrderDetail />
      </MainLayout>
    </ProtectedRoute>
  } />
  <Route path="/profile" element={
    <ProtectedRoute>
      <MainLayout>
        <UserProfile />
      </MainLayout>
    </ProtectedRoute>
  } />
</Routes>
```

## 技术特性

### 1. 类型安全
- 使用TypeScript严格模式
- 完整的类型定义
- 接口与后端API保持一致

### 2. 性能优化
- React Query数据缓存
- 组件懒加载
- 图片懒加载
- 代码分割

### 3. 响应式设计
- 移动端优先
- 自适应布局
- 触摸友好的交互

### 4. 用户体验
- 加载状态管理
- 错误边界处理
- 表单验证
- 操作反馈

### 5. 代码质量
- ESLint代码检查
- Prettier代码格式化
- 组件化架构
- 可复用组件

## 与后端集成

### API规范
- RESTful API设计
- 统一响应格式
- JWT认证
- 分页参数标准化

### 数据模型
前端类型定义与后端数据模型完全一致，确保类型安全。

### 错误处理
- 统一的错误处理中间件
- 用户友好的错误提示
- 网络错误重试机制

## 部署

### 构建
```bash
npm run build
```

### 部署到静态服务器
将`dist`目录部署到任何静态文件服务器。

### Docker部署
```dockerfile
FROM nginx:alpine
COPY dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License