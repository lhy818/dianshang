import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider, Layout, message } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { QueryProvider } from './providers/QueryProvider';
import { useAuthStore } from './store/authStore';

// 导入组件
import { LoginForm } from './components/auth/LoginForm';
import { RegisterForm } from './components/auth/RegisterForm';
import { UserProfile } from './components/auth/UserProfile';
import { ProductList } from './components/products/ProductList';
import { ProductDetail } from './components/products/ProductDetail';
import { CartPage } from './components/cart/CartPage';
import { OrderList } from './components/orders/OrderList';
import { OrderDetail } from './components/orders/OrderDetail';

const { Header, Content, Footer } = Layout;

// 受保护的路由组件
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuthStore();
  
  if (!isAuthenticated) {
    message.warning('请先登录');
    return <Navigate to="/login" />;
  }
  
  return <>{children}</>;
};

// 公共布局组件
const MainLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ 
        background: '#fff', 
        padding: '0 24px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between'
      }}>
        <div style={{ fontSize: 20, fontWeight: 'bold', color: '#1890ff' }}>
          电商系统
        </div>
        <div>
          {/* 这里可以添加导航菜单 */}
        </div>
      </Header>
      
      <Content style={{ padding: '24px', background: '#f0f2f5' }}>
        {children}
      </Content>
      
      <Footer style={{ textAlign: 'center', background: '#fff' }}>
        电商系统 ©2026 版权所有
      </Footer>
    </Layout>
  );
};

function App() {
  return (
    <ConfigProvider
      locale={zhCN}
      theme={{
        token: {
          colorPrimary: '#1890ff',
          borderRadius: 6,
        },
      }}
    >
      <QueryProvider>
        <Router>
          <Routes>
            {/* 公开路由 */}
            <Route path="/login" element={<LoginForm />} />
            <Route path="/register" element={<RegisterForm />} />
            
            {/* 受保护的路由 */}
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
            
            {/* 默认重定向 */}
            <Route path="*" element={<Navigate to="/" />} />
          </Routes>
        </Router>
      </QueryProvider>
    </ConfigProvider>
  );
}

export default App;