import React, { useEffect, useState } from 'react';
import { Card, Avatar, Typography, Space, Button, Row, Col, Descriptions, Tag, message } from 'antd';
import { UserOutlined, MailOutlined, PhoneOutlined, EditOutlined, LogoutOutlined } from '@ant-design/icons';
import { useAuthStore } from '../../store/authStore';
import { authApi } from '../../api/services';
import { User } from '../../types';

const { Title, Text } = Typography;

export const UserProfile: React.FC = () => {
  const { user, logout } = useAuthStore();
  const [userData, setUserData] = useState<User | null>(user);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (!userData) {
      fetchUserProfile();
    }
  }, []);

  const fetchUserProfile = async () => {
    try {
      setIsLoading(true);
      const profile = await authApi.getProfile();
      setUserData(profile);
    } catch (error) {
      message.error('获取用户信息失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
    authApi.logout().catch(() => {
      // 忽略登出API错误，因为本地状态已清除
    });
    message.success('已退出登录');
  };

  if (!userData) {
    return (
      <Card loading={isLoading}>
        <div style={{ textAlign: 'center', padding: 50 }}>
          <Text type="secondary">加载中...</Text>
        </div>
      </Card>
    );
  }

  const getStatusTag = (status: User['status']) => {
    const statusMap = {
      active: { color: 'success', text: '活跃' },
      inactive: { color: 'default', text: '未激活' },
      banned: { color: 'error', text: '已禁用' },
    };
    const { color, text } = statusMap[status];
    return <Tag color={color}>{text}</Tag>;
  };

  return (
    <Card>
      <Row gutter={[24, 24]}>
        <Col xs={24} md={8}>
          <div style={{ textAlign: 'center' }}>
            <Avatar
              size={120}
              src={userData.avatarUrl}
              icon={!userData.avatarUrl && <UserOutlined />}
              style={{ marginBottom: 16 }}
            />
            <Title level={3}>{userData.username}</Title>
            <Space direction="vertical" size="small">
              {getStatusTag(userData.status)}
              {userData.isEmailVerified && <Tag color="blue">邮箱已验证</Tag>}
              {userData.isPhoneVerified && <Tag color="green">手机已验证</Tag>}
            </Space>
          </div>
        </Col>
        
        <Col xs={24} md={16}>
          <Descriptions title="个人信息" column={1} bordered>
            <Descriptions.Item label="用户名">
              <Text strong>{userData.username}</Text>
            </Descriptions.Item>
            
            <Descriptions.Item label="邮箱">
              <Space>
                <MailOutlined />
                <Text>{userData.email}</Text>
                {userData.isEmailVerified && (
                  <Tag color="success" size="small">已验证</Tag>
                )}
              </Space>
            </Descriptions.Item>
            
            {userData.phone && (
              <Descriptions.Item label="手机">
                <Space>
                  <PhoneOutlined />
                  <Text>{userData.phone}</Text>
                  {userData.isPhoneVerified && (
                    <Tag color="success" size="small">已验证</Tag>
                  )}
                </Space>
              </Descriptions.Item>
            )}
            
            <Descriptions.Item label="账户状态">
              {getStatusTag(userData.status)}
            </Descriptions.Item>
            
            <Descriptions.Item label="最后登录时间">
              <Text>
                {userData.lastLoginAt 
                  ? new Date(userData.lastLoginAt).toLocaleString()
                  : '从未登录'
                }
              </Text>
            </Descriptions.Item>
            
            <Descriptions.Item label="注册时间">
              <Text>{new Date(userData.createdAt).toLocaleString()}</Text>
            </Descriptions.Item>
          </Descriptions>
          
          <Space style={{ marginTop: 24 }}>
            <Button type="primary" icon={<EditOutlined />}>
              编辑资料
            </Button>
            <Button 
              danger 
              icon={<LogoutOutlined />}
              onClick={handleLogout}
            >
              退出登录
            </Button>
          </Space>
        </Col>
      </Row>
    </Card>
  );
};