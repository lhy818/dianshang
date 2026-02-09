# 前端包优化配置指南

## 1. 依赖包优化策略

### 1.1 按需加载配置

#### Ant Design 按需加载
```json
// package.json 更新
{
  "dependencies": {
    "antd": "^5.12.0",
    "@ant-design/icons": "^5.2.6",
    "babel-plugin-import": "^1.13.8"
  }
}
```

```javascript
// vite.config.ts 添加
import { createStyleImportPlugin, AntdResolve } from 'vite-plugin-style-import'

export default defineConfig({
  plugins: [
    createStyleImportPlugin({
      resolves: [AntdResolve()],
    }),
  ],
})
```

#### React Icons 优化
```javascript
// 使用具体图标而不是完整包
import { FaUser } from 'react-icons/fa'
// 而不是
import { FaUser } from 'react-icons'
```

### 1.2 替换重型依赖

考虑替换方案：
1. **recharts** → **victory** 或 **chart.js** (更轻量)
2. **framer-motion** → **react-spring** (性能更好)
3. **antd** → **保持但按需加载**

## 2. 代码分割配置

### 2.1 路由级代码分割
```typescript
// 使用React.lazy进行路由懒加载
import React, { lazy, Suspense } from 'react'

const ProductList = lazy(() => import('./components/products/ProductList'))
const ProductDetail = lazy(() => import('./components/products/ProductDetail'))
const CartPage = lazy(() => import('./components/cart/CartPage'))

// 在路由中使用
<Route path="/products" element={
  <Suspense fallback={<LoadingSpinner />}>
    <ProductList />
  </Suspense>
} />
```

### 2.2 组件级代码分割
```typescript
// 对于大型组件使用动态导入
const HeavyComponent = React.lazy(() => 
  import('./components/HeavyComponent').then(module => ({
    default: module.HeavyComponent
  }))
)
```

## 3. 图片和资源优化

### 3.1 图片优化策略
1. **使用WebP格式**: 比JPEG小25-35%
2. **实现懒加载**: 使用Intersection Observer
3. **响应式图片**: 使用srcset属性
4. **CDN加速**: 配置图片CDN

### 3.2 字体优化
```css
/* 使用字体子集和woff2格式 */
@font-face {
  font-family: 'CustomFont';
  src: url('/fonts/custom.woff2') format('woff2'),
       url('/fonts/custom.woff') format('woff');
  font-display: swap; /* 避免FOIT */
}
```

## 4. 构建优化配置

### 4.1 安装优化插件
```bash
npm install --save-dev \
  vite-plugin-compression \
  vite-plugin-image-optimizer \
  rollup-plugin-visualizer \
  vite-plugin-style-import \
  @rollup/plugin-terser
```

### 4.2 环境特定配置
```typescript
// vite.config.prod.ts
export default defineConfig({
  build: {
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true,
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          // 生产环境更细粒度的代码分割
        },
      },
    },
  },
})
```

## 5. 性能监控配置

### 5.1 Web Vitals监控
```typescript
// src/utils/webVitals.ts
import { onCLS, onFID, onLCP } from 'web-vitals'

export function reportWebVitals() {
  onCLS(console.log)
  onFID(console.log)
  onLCP(console.log)
  
  // 发送到分析服务
  function sendToAnalytics(metric) {
    const body = JSON.stringify(metric)
    navigator.sendBeacon('/analytics', body)
  }
  
  onCLS(sendToAnalytics)
  onFID(sendToAnalytics)
  onLCP(sendToAnalytics)
}
```

### 5.2 错误监控
```typescript
// src/utils/errorTracking.ts
import * as Sentry from '@sentry/react'

Sentry.init({
  dsn: process.env.VITE_SENTRY_DSN,
  integrations: [
    new Sentry.BrowserTracing({
      tracingOrigins: ['localhost', 'your-domain.com'],
    }),
  ],
  tracesSampleRate: 1.0,
})
```

## 6. 缓存策略

### 6.1 静态资源缓存
```nginx
# Nginx配置
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
  expires 1y;
  add_header Cache-Control "public, immutable";
}

location ~* \.(woff|woff2|ttf|eot)$ {
  expires 1y;
  add_header Cache-Control "public, immutable";
}
```

### 6.2 Service Worker缓存
```javascript
// public/sw.js
const CACHE_NAME = 'v1'
const urlsToCache = [
  '/',
  '/index.html',
  '/assets/main.css',
  '/assets/main.js',
]

self.addEventListener('install', event => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(cache => cache.addAll(urlsToCache))
  )
})
```

## 7. 移动端优化

### 7.1 触摸优化
```css
/* 增加触摸目标大小 */
button, 
a, 
input[type="submit"] {
  min-height: 44px;
  min-width: 44px;
}

/* 禁用双击缩放 */
* {
  touch-action: manipulation;
}
```

### 7.2 性能预算
```json
// .performance-budget.json
{
  "budgets": [
    {
      "resourceType": "document",
      "budget": 50
    },
    {
      "resourceType": "script",
      "budget": 200
    },
    {
      "resourceType": "image",
      "budget": 300
    }
  ]
}
```

## 8. 实施步骤

### 阶段1: 立即实施 (1-2天)
1. 更新Vite配置使用优化版本
2. 配置Ant Design按需加载
3. 实现路由级代码分割
4. 添加图片优化插件

### 阶段2: 短期优化 (1周)
1. 集成性能监控
2. 实现Service Worker
3. 优化字体加载
4. 配置CDN

### 阶段3: 长期优化 (1个月)
1. 替换重型依赖
2. 实现高级缓存策略
3. 深度性能调优
4. A/B测试集成

## 9. 预期效果

| 指标 | 当前 | 优化后 | 改进 |
|------|------|--------|------|
| 包大小 | ~5MB | ~2MB | -60% |
| 首次加载 | 3-5s | 1-2s | -60% |
| LCP | 3-5s | <2.5s | -50% |
| CLS | 0.1-0.3 | <0.1 | -70% |
| 缓存命中率 | 0% | 80% | +80% |

---
*优化配置生成时间: 2026-02-08 13:28*