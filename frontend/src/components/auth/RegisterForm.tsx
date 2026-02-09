import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Button, Form, Input, Card, Typography, message, Space } from 'antd';
import { UserOutlined, MailOutlined, LockOutlined, PhoneOutlined } from '@ant-design/icons';
import { useAuthStore } from '../../store/authStore';
import { authApi } from '../../api/services';
import { RegisterRequest } from '../../types';

const { Title, Text } = Typography;

export const RegisterForm: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const { login, setLoading, setError } = useAuthStore();
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (values: RegisterRequest) => {
    try {
      setIsLoading(true);
      setLoading(true);
      
      const response = await authApi.register(values);
      login(response.user, response.token);
      
      message.success('注册成功！');
      navigate('/');
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || '注册失败，请稍后重试';
      setError(errorMessage);
      message.error(errorMessage);
    } finally {
      setIsLoading(false);
      setLoading(false);
    }
  };

  const validatePassword = (_: any, value: string) => {
    if (!value) {
      return Promise.reject('请输入密码');
    }
    if (value.length < 6) {
      return Promise.reject('密码至少6个字符');
    }
    if (!/(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/.test(value)) {
      return Promise.reject('密码必须包含大小写字母和数字');
    }
    return Promise.resolve();
  };

  const validateConfirmPassword = ({ getFieldValue }: any) => ({
    validator(_: any, value: string) {
      if (!value || getFieldValue('password') === value) {
        return Promise.resolve();
      }
      return Promise.reject('两次输入的密码不一致');
    },
  });

  return (
    <Card style={{ maxWidth: 400, margin: '0 auto', marginTop: 50 }}>
      <div style={{ textAlign: 'center', marginBottom: 24 }}>
        <Title level={2}>用户注册</Title>
        <Text type="secondary">创建新账户，开始购物之旅</Text>
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
            { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线' },
          ]}
        >
          <Input
            prefix={<UserOutlined />}
            placeholder="用户名"
            size="large"
          />
        </Form.Item>
        
        <Form.Item
          name="email"
          rules={[
            { required: true, message: '请输入邮箱' },
            { type: 'email', message: '请输入有效的邮箱地址' },
          ]}
        >
          <Input
            prefix={<MailOutlined />}
            placeholder="邮箱"
            size="large"
          />
        </Form.Item>
        
        <Form.Item
          name="phone"
          rules={[
            { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码' },
          ]}
        >
          <Input
            prefix={<PhoneOutlined />}
            placeholder="手机号码（可选）"
            size="large"
          />
        </Form.Item>
        
        <Form.Item
          name="password"
          rules={[
            { required: true, validator: validatePassword },
          ]}
        >
          <Input.Password
            prefix={<LockOutlined />}
            placeholder="密码"
            size="large"
          />
        </Form.Item>
        
        <Form.Item
          name="confirmPassword"
          dependencies={['password']}
          rules={[
            { required: true, message: '请确认密码' },
            validateConfirmPassword,
          ]}
        >
          <Input.Password
            prefix={<LockOutlined />}
            placeholder="确认密码"
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
            注册
          </Button>
        </Form.Item>
        
        <Space direction="vertical" style={{ width: '100%', textAlign: 'center' }}>
          <Link to="/login">已有账户？立即登录</Link>
        </Space>
      </Form>
    </Card>
  );
};