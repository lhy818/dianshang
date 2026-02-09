import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Button, Form, Input, Card, Typography, message, Space } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useAuthStore } from '../../store/authStore';
import { authApi } from '../../api/services';
import { LoginRequest } from '../../types';

const { Title, Text } = Typography;

export const LoginForm: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const { login, setLoading, setError } = useAuthStore();
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (values: LoginRequest) => {
    try {
      setIsLoading(true);
      setLoading(true);
      
      const response = await authApi.login(values);
      login(response.user, response.token);
      
      message.success('登录成功！');
      navigate('/');
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || '登录失败，请检查用户名和密码';
      setError(errorMessage);
      message.error(errorMessage);
    } finally {
      setIsLoading(false);
      setLoading(false);
    }
  };

  return (
    <Card style={{ maxWidth: 400, margin: '0 auto', marginTop: 50 }}>
      <div style={{ textAlign: 'center', marginBottom: 24 }}>
        <Title level={2}>用户登录</Title>
        <Text type="secondary">欢迎回来，请登录您的账户</Text>
      </div>
      
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        autoComplete="off"
      >
        <Form.Item
          name="username"
          rules={[
            { required: true, message: '请输入用户名' },
            { min: 3, message: '用户名至少3个字符' },
            { max: 20, message: '用户名最多20个字符' },
          ]}
        >
          <Input
            prefix={<UserOutlined />}
            placeholder="用户名"
            size="large"
          />
        </Form.Item>
        
        <Form.Item
          name="password"
          rules={[
            { required: true, message: '请输入密码' },
            { min: 6, message: '密码至少6个字符' },
          ]}
        >
          <Input.Password
            prefix={<LockOutlined />}
            placeholder="密码"
            size="large"
          />
        </Form.Item>
        
        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            size="large"
            loading={isLoading}
            block
          >
            登录
          </Button>
        </Form.Item>
        
        <Space direction="vertical" style={{ width: '100%', textAlign: 'center' }}>
          <Link to="/register">还没有账户？立即注册</Link>
          <Link to="/forgot-password">忘记密码？</Link>
        </Space>
      </Form>
    </Card>
  );
};