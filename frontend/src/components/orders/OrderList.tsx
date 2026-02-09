import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import {
  Card,
  Table,
  Tag,
  Button,
  Space,
  Typography,
  Input,
  Select,
  DatePicker,
  Row,
  Col,
  Pagination,
  Spin,
  Empty,
  Dropdown,
  MenuProps,
  message
} from 'antd';
import {
  SearchOutlined,
  EyeOutlined,
  MoreOutlined,
  ReloadOutlined,
  FileTextOutlined
} from '@ant-design/icons';
import { orderApi } from '../../api/services';
import { Order } from '../../types';
import dayjs from 'dayjs';

const { Title, Text } = Typography;
const { Search } = Input;
const { Option } = Select;
const { RangePicker } = DatePicker;

export const OrderList: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useState({
    page: 1,
    pageSize: 10,
    status: undefined as string | undefined,
    search: '',
    startDate: '',
    endDate: '',
  });

  // 获取订单列表
  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['orders', searchParams],
    queryFn: () => orderApi.getOrders(searchParams),
  });

  const handleSearch = (value: string) => {
    setSearchParams(prev => ({ ...prev, search: value, page: 1 }));
  };

  const handleStatusChange = (value: string) => {
    setSearchParams(prev => ({ ...prev, status: value, page: 1 }));
  };

  const handleDateChange = (dates: any) => {
    if (dates) {
      setSearchParams(prev => ({
        ...prev,
        startDate: dates[0].format('YYYY-MM-DD'),
        endDate: dates[1].format('YYYY-MM-DD'),
        page: 1,
      }));
    } else {
      setSearchParams(prev => ({
        ...prev,
        startDate: '',
        endDate: '',
        page: 1,
      }));
    }
  };

  const handlePageChange = (page: number, pageSize: number) => {
    setSearchParams(prev => ({ ...prev, page, pageSize }));
  };

  const getStatusTag = (status: Order['status']) => {
    const statusConfig = {
      pending: { color: 'orange', text: '待付款' },
      paid: { color: 'blue', text: '已付款' },
      shipped: { color: 'cyan', text: '已发货' },
      delivered: { color: 'green', text: '已送达' },
      completed: { color: 'success', text: '已完成' },
      cancelled: { color: 'red', text: '已取消' },
      refunded: { color: 'purple', text: '已退款' },
    };
    const config = statusConfig[status];
    return <Tag color={config.color}>{config.text}</Tag>;
  };

  const getPaymentStatusTag = (status: Order['paymentStatus']) => {
    const statusConfig = {
      unpaid: { color: 'orange', text: '未支付' },
      paid: { color: 'green', text: '已支付' },
      refunded: { color: 'purple', text: '已退款' },
      failed: { color: 'red', text: '支付失败' },
    };
    const config = statusConfig[status];
    return <Tag color={config.color}>{config.text}</Tag>;
  };

  const handleCancelOrder = async (orderId: string) => {
    try {
      await orderApi.cancelOrder(orderId);
      message.success('订单已取消');
      refetch();
    } catch (error) {
      message.error('取消订单失败');
    }
  };

  const columns = [
    {
      title: '订单号',
      dataIndex: 'orderNo',
      key: 'orderNo',
      width: 180,
      render: (orderNo: string) => (
        <Text strong copyable>{orderNo}</Text>
      ),
    },
    {
      title: '订单状态',
      dataIndex: 'status',
      key: 'status',
      width: 120,
      render: (status: Order['status']) => getStatusTag(status),
    },
    {
      title: '支付状态',
      dataIndex: 'paymentStatus',
      key: 'paymentStatus',
      width: 120,
      render: (status: Order['paymentStatus']) => getPaymentStatusTag(status),
    },
    {
      title: '订单金额',
      dataIndex: 'finalAmount',
      key: 'finalAmount',
      width: 120,
      render: (amount: number) => (
        <Text strong style={{ color: '#ff4d4f' }}>
          ¥{amount.toFixed(2)}
        </Text>
      ),
    },
    {
      title: '下单时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 180,
      render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      key: 'actions',
      width: 150,
      render: (_: any, record: Order) => {
        const items: MenuProps['items'] = [
          {
            key: 'view',
            label: '查看详情',
            icon: <EyeOutlined />,
            onClick: () => navigate(`/orders/${record.id}`),
          },
          {
            key: 'invoice',
            label: '下载发票',
            icon: <FileTextOutlined />,
            disabled: record.paymentStatus !== 'paid',
          },
        ];

        if (record.status === 'pending') {
          items.push({
            key: 'cancel',
            label: '取消订单',
            danger: true,
            onClick: () => handleCancelOrder(record.id),
          });
        }

        if (record.status === 'pending' && record.paymentStatus === 'unpaid') {
          items.push({
            key: 'pay',
            label: '去支付',
            onClick: () => navigate(`/payment/${record.id}`),
          });
        }

        return (
          <Space>
            <Button
              type="link"
              icon={<EyeOutlined />}
              onClick={() => navigate(`/orders/${record.id}`)}
            >
              查看
            </Button>
            <Dropdown menu={{ items }} placement="bottomRight">
              <Button type="text" icon={<MoreOutlined />} />
            </Dropdown>
          </Space>
        );
      },
    },
  ];

  if (isError) {
    return (
      <Card>
        <Empty description="加载订单失败，请稍后重试" />
      </Card>
    );
  }

  return (
    <div>
      <Title level={2}>我的订单</Title>
      
      {/* 搜索和筛选 */}
      <Card style={{ marginBottom: 24 }}>
        <Row gutter={[16, 16]}>
          <Col xs={24} md={8}>
            <Search
              placeholder="搜索订单号"
              enterButton={<SearchOutlined />}
              size="large"
              onSearch={handleSearch}
              allowClear
            />
          </Col>
          
          <Col xs={24} md={8}>
            <Select
              placeholder="订单状态"
              style={{ width: '100%' }}
              size="large"
              allowClear
              onChange={handleStatusChange}
            >
              <Option value="pending">待付款</Option>
              <Option value="paid">已付款</Option>
              <Option value="shipped">已发货</Option>
              <Option value="delivered">已送达</Option>
              <Option value="completed">已完成</Option>
              <Option value="cancelled">已取消</Option>
            </Select>
          </Col>
          
          <Col xs={24} md={8}>
            <RangePicker
              style={{ width: '100%' }}
              size="large"
              onChange={handleDateChange}
              format="YYYY-MM-DD"
            />
          </Col>
        </Row>
        
        <Row justify="end" style={{ marginTop: 16 }}>
          <Button
            icon={<ReloadOutlined />}
            onClick={() => refetch()}
          >
            刷新
          </Button>
        </Row>
      </Card>

      {/* 订单列表 */}
      <Card>
        {isLoading ? (
          <div style={{ textAlign: 'center', padding: 50 }}>
            <Spin size="large" />
          </div>
        ) : (
          <>
            <Table
              dataSource={data?.items}
              columns={columns}
              rowKey="id"
              pagination={false}
              scroll={{ x: 800 }}
            />
            
            {data && data.pagination.totalPages > 1 && (
              <div style={{ textAlign: 'center', marginTop: 24 }}>
                <Pagination
                  current={data.pagination.page}
                  pageSize={data.pagination.pageSize}
                  total={data.pagination.total}
                  onChange={handlePageChange}
                  showSizeChanger
                  showQuickJumper
                  showTotal={(total, range) => 
                    `第 ${range[0]}-${range[1]} 条，共 ${total} 条`
                  }
                />
              </div>
            )}
          </>
        )}
      </Card>
    </div>
  );
};