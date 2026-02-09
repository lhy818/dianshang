import React from 'react';
import { Card, Typography, Space, Divider } from 'antd';
import { LoginForm } from '../components/auth/LoginForm';
import { RegisterForm } from '../components/auth/RegisterForm';
import { UserProfile } from '../components/auth/UserProfile';
import { ProductList } from '../components/products/ProductList';
import { ProductDetail } from '../components/products/ProductDetail';
import { CartPage } from '../components/cart/CartPage';
import { OrderList } from '../components/orders/OrderList';
import { OrderDetail } from '../components/orders/OrderDetail';

const { Title, Text } = Typography;

export const ComponentExamples: React.FC = () => {
  return (
    <div style={{ padding: 24 }}>
      <Title level={2}>电商系统组件使用示例</Title>
      <Text type="secondary">以下展示了所有核心组件的使用方式</Text>
      
      <Divider />
      
      <Space direction="vertical" size="large" style={{ width: '100%' }}>
        
        {/* 认证模块 */}
        <Card title="1. 用户认证模块">
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <div>
              <Title level={4}>登录组件 (LoginForm)</Title>
              <Text type="secondary">用于用户登录，包含表单验证和API调用</Text>
              <div style={{ marginTop: 16 }}>
                <LoginForm />
              </div>
            </div>
            
            <Divider />
            
            <div>
              <Title level={4}>注册组件 (RegisterForm)</Title>
              <Text type="secondary">用于新用户注册，包含密码强度验证</Text>
              <div style={{ marginTop: 16 }}>
                <RegisterForm />
              </div>
            </div>
            
            <Divider />
            
            <div>
              <Title level={4}>用户信息组件 (UserProfile)</Title>
              <Text type="secondary">显示用户个人信息和状态</Text>
              <div style={{ marginTop: 16 }}>
                <UserProfile />
              </div>
            </div>
          </Space>
        </Card>
        
        {/* 商品模块 */}
        <Card title="2. 商品展示模块">
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <div>
              <Title level={4}>商品列表组件 (ProductList)</Title>
              <Text type="secondary">支持分页、筛选、排序的商品列表</Text>
              <div style={{ marginTop: 16 }}>
                <ProductList showFilters={false} />
              </div>
            </div>
            
            <Divider />
            
            <div>
              <Title level={4}>商品详情组件 (ProductDetail)</Title>
              <Text type="secondary">商品详细信息页面，包含图片轮播、规格参数等</Text>
              <Text type="secondary" style={{ display: 'block', marginTop: 8 }}>
                注：需要商品ID参数，这里展示的是组件结构
              </Text>
            </div>
          </Space>
        </Card>
        
        {/* 购物车模块 */}
        <Card title="3. 购物车模块">
          <div>
            <Title level={4}>购物车页面组件 (CartPage)</Title>
            <Text type="secondary">购物车管理页面，支持批量操作和结算</Text>
            <div style={{ marginTop: 16 }}>
              <CartPage />
            </div>
          </div>
        </Card>
        
        {/* 订单模块 */}
        <Card title="4. 订单模块">
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <div>
              <Title level={4}>订单列表组件 (OrderList)</Title>
              <Text type="secondary">用户订单列表，支持搜索和筛选</Text>
              <div style={{ marginTop: 16 }}>
                <OrderList />
              </div>
            </div>
            
            <Divider />
            
            <div>
              <Title level={4}>订单详情组件 (OrderDetail)</Title>
              <Text type="secondary">订单详细信息页面，包含时间线和商品列表</Text>
              <Text type="secondary" style={{ display: 'block', marginTop: 8 }}>
                注：需要订单ID参数，这里展示的是组件结构
              </Text>
            </div>
          </Space>
        </Card>
        
      </Space>
      
      <Divider />
      
      <Card title="组件使用说明">
        <Space direction="vertical" size="small">
          <Text strong>安装依赖：</Text>
          <Text code>npm install @tanstack/react-query zustand antd @ant-design/icons axios react-router-dom</Text>
          
          <Text strong>状态管理：</Text>
          <Text>使用Zustand进行全局状态管理，已创建authStore和cartStore</Text>
          
          <Text strong>数据获取：</Text>
          <Text>使用React Query进行数据获取和缓存，已配置QueryProvider</Text>
          
          <Text strong>API调用：</Text>
          <Text>已封装axios客户端，包含请求拦截器和错误处理</Text>
          
          <Text strong>类型安全：</Text>
          <Text>使用TypeScript严格模式，所有类型定义在src/types/index.ts</Text>
          
          <Text strong>响应式设计：</Text>
          <Text>使用Ant Design的Grid系统和响应式组件</Text>
        </Space>
      </Card>
    </div>
  );
};