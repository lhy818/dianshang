import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import {
  Card,
  Typography,
  Button,
  Space,
  Tag,
  Descriptions,
  Table,
  Timeline,
  Row,
  Col,
  Divider,
  Spin,
  message,
  Popconfirm
} from 'antd';
import {
  ArrowLeftOutlined,
  PrinterOutlined,
  DownloadOutlined,
  CloseCircleOutlined,
  CheckCircleOutlined,
  TruckOutlined,
  HomeOutlined,
  CreditCardOutlined
} from '@ant-design/icons';
import { orderApi } from '../../api/services';
import { Order, OrderItem } from '../../types';
import dayjs from 'dayjs';

const { Title, Text } = Typography;

export const OrderDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  // 获取订单详情
  const { data: order, isLoading, isError } = useQuery({
    queryKey: ['order', id],
    queryFn: () => orderApi.getOrder(id!),
    enabled: !!id,
  });

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

  const handleCancelOrder = async () => {
    try {
      await orderApi.cancelOrder(id!);
      message.success('订单已取消');
      navigate('/orders');
    } catch (error) {
      message.error('取消订单失败');
    }
  };

  const handlePayOrder = async () => {
    try {
      await orderApi.payOrder(id!);
      message.success('支付请求已发送');
      // 这里可以跳转到支付页面
    } catch (error) {
      message.error('支付失败');
    }
  };

  const orderItemsColumns = [
    {
      title: '商品',
      dataIndex: 'productName',
      key: 'productName',
      width: '40%',
      render: (productName: string, record: OrderItem) => (
        <Space>
          <div style={{ width: 60, height: 60, backgroundColor: '#f5f5f5' }} />
          <Space direction="vertical" size="small">
            <Text strong>{productName}</Text>
            {record.skuAttributes && Object.keys(record.skuAttributes).length > 0 && (
              <Text type="secondary">
                规格: {Object.values(record.skuAttributes).join(', ')}
              </Text>
            )}
          </Space>
        </Space>
      ),
    },
    {
      title: '单价',
      dataIndex: 'unitPrice',
      key: 'unitPrice',
      width: '15%',
      render: (price: number) => (
        <Text>¥{price.toFixed(2)}</Text>
      ),
    },
    {
      title: '数量',
      dataIndex: 'quantity',
      key: 'quantity',
      width: '15%',
    },
    {
      title: '小计',
      key: 'subtotal',
      width: '15%',
      render: (_: any, record: OrderItem) => (
        <Text strong>¥{(record.unitPrice * record.quantity).toFixed(2)}</Text>
      ),
    },
  ];

  const getOrderTimeline = (order: Order) => {
    const items = [
      {
        color: 'green',
        dot: <CheckCircleOutlined />,
        children: (
          <Space direction="vertical" size="small">
            <Text strong>订单创建</Text>
            <Text type="secondary">{dayjs(order.createdAt).format('YYYY-MM-DD HH:mm:ss')}</Text>
          </Space>
        ),
      },
    ];

    if (order.paymentTime) {
      items.push({
        color: 'blue',
        dot: <CreditCardOutlined />,
        children: (
          <Space direction="vertical" size="small">
            <Text strong>支付成功</Text>
            <Text type="secondary">{dayjs(order.paymentTime).format('YYYY-MM-DD HH:mm:ss')}</Text>
            <Text>支付方式: {order.paymentMethod}</Text>
          </Space>
        ),
      });
    }

    if (order.shippingTime) {
      items.push({
        color: 'cyan',
        dot: <TruckOutlined />,
        children: (
          <Space direction="vertical" size="small">
            <Text strong>已发货</Text>
            <Text type="secondary">{dayjs(order.shippingTime).format('YYYY-MM-DD HH:mm:ss')}</Text>
            {order.shippingNo && <Text>物流单号: {order.shippingNo}</Text>}
          </Space>
        ),
      });
    }

    if (order.deliveredTime) {
      items.push({
        color: 'green',
        dot: <HomeOutlined />,
        children: (
          <Space direction="vertical" size="small">
            <Text strong>已送达</Text>
            <Text type="secondary">{dayjs(order.deliveredTime).format('YYYY-MM-DD HH:mm:ss')}</Text>
          </Space>
        ),
      });
    }

    if (order.cancelledTime) {
      items.push({
        color: 'red',
        dot: <CloseCircleOutlined />,
        children: (
          <Space direction="vertical" size="small">
            <Text strong>订单取消</Text>
            <Text type="secondary">{dayjs(order.cancelledTime).format('YYYY-MM-DD HH:mm:ss')}</Text>
            {order.cancelledReason && <Text>原因: {order.cancelledReason}</Text>}
          </Space>
        ),
      });
    }

    return items;
  };

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: 50 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (isError || !order) {
    return (
      <Card>
        <div style={{ textAlign: 'center', padding: 50 }}>
          <Title level={3}>订单不存在或加载失败</Title>
          <Button type="primary" onClick={() => navigate('/orders')}>
            返回订单列表
          </Button>
        </div>
      </Card>
    );
  }

  return (
    <div>
      {/* 头部操作栏 */}
      <Row justify="space-between" align="middle" style={{ marginBottom: 24 }}>
        <Col>
          <Button
            type="link"
            icon={<ArrowLeftOutlined />}
            onClick={() => navigate('/orders')}
          >
            返回订单列表
          </Button>
        </Col>
        <Col>
          <Space>
            <Button icon={<PrinterOutlined />}>打印订单</Button>
            <Button icon={<DownloadOutlined />}>下载发票</Button>
            
            {order.status === 'pending' && order.paymentStatus === 'unpaid' && (
              <>
                <Popconfirm
                  title="确定要取消订单吗？"
                  onConfirm={handleCancelOrder}
                  okText="确定"
                  cancelText="取消"
                >
                  <Button danger>取消订单</Button>
                </Popconfirm>
                <Button type="primary" onClick={handlePayOrder}>
                  去支付
                </Button>
              </>
            )}
          </Space>
        </Col>
      </Row>

      <Row gutter={[24, 24]}>
        {/* 订单信息 */}
        <Col xs={24} lg={16}>
          <Card title="订单信息">
            <Descriptions column={2} bordered>
              <Descriptions.Item label="订单号">
                <Text strong copyable>{order.orderNo}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="订单状态">
                {getStatusTag(order.status)}
              </Descriptions.Item>
              <Descriptions.Item label="支付状态">
                {getPaymentStatusTag(order.paymentStatus)}
              </Descriptions.Item>
              <Descriptions.Item label="支付方式">
                {order.paymentMethod || '未选择'}
              </Descriptions.Item>
              <Descriptions.Item label="下单时间">
                {dayjs(order.createdAt).format('YYYY-MM-DD HH:mm:ss')}
              </Descriptions.Item>
              <Descriptions.Item label="支付时间">
                {order.paymentTime 
                  ? dayjs(order.paymentTime).format('YYYY-MM-DD HH:mm:ss')
                  : '-'
                }
              </Descriptions.Item>
              {order.shippingNo && (
                <Descriptions.Item label="物流单号">
                  <Text copyable>{order.shippingNo}</Text>
                </Descriptions.Item>
              )}
              {order.shippingMethod && (
                <Descriptions.Item label="配送方式">
                  {order.shippingMethod}
                </Descriptions.Item>
              )}
              {order.buyerRemark && (
                <Descriptions.Item label="买家留言" span={2}>
                  {order.buyerRemark}
                </Descriptions.Item>
              )}
            </Descriptions>
          </Card>

          <Divider />

          {/* 订单商品 */}
          <Card title="订单商品">
            <Table
              dataSource={[]} // 这里需要从API获取订单商品列表
              columns={orderItemsColumns}
              rowKey="id"
              pagination={false}
              summary={() => (
                <Table.Summary>
                  <Table.Summary.Row>
                    <Table.Summary.Cell index={0} colSpan={3} align="right">
                      <Text strong>商品合计:</Text>
                    </Table.Summary.Cell>
                    <Table.Summary.Cell index={1}>
                      <Text strong>¥{order.totalAmount.toFixed(2)}</Text>
                    </Table.Summary.Cell>
                  </Table.Summary.Row>
                  <Table.Summary.Row>
                    <Table.Summary.Cell index={0} colSpan={3} align="right">
                      <Text>运费:</Text>
                    </Table.Summary.Cell>
                    <Table.Summary.Cell index={1}>
                      <Text>¥{order.shippingFee.toFixed(2)}</Text>
                    </Table.Summary.Cell>
                  </Table.Summary.Row>
                  <Table.Summary.Row>
                    <Table.Summary.Cell index={0} colSpan={3} align="right">
                      <Text>优惠:</Text>
                    </Table.Summary.Cell>
                    <Table.Summary.Cell index={1}>
                      <Text type="success">-¥{order.discountAmount.toFixed(2)}</Text>
                    </Table.Summary.Cell>
                  </Table.Summary.Row>
                  <Table.Summary.Row>
                    <Table.Summary.Cell index={0} colSpan={3} align="right">
                      <Text strong>实付金额:</Text>
                    </Table.Summary.Cell>
                    <Table.Summary.Cell index={1}>
                      <Text strong style={{ color: '#ff4d4f', fontSize: 16 }}>
                        ¥{order.finalAmount.toFixed(2)}
                      </Text>
                    </Table.Summary.Cell>
                  </Table.Summary.Row>
                </Table.Summary>
              )}
            />
          </Card>
        </Col>

        {/* 订单时间线 */}
        <Col xs={24} lg={8}>
          <Card title="订单进度">
            <Timeline items={getOrderTimeline(order)} />
          </Card>

          {/* 收货地址 */}
          <Card title="收货地址" style={{ marginTop: 24 }}>
            {order.shippingAddress ? (
              <Space direction="vertical">
                <Text strong>{order.shippingAddress.recipientName}</Text>
                <Text>{order.shippingAddress.phone}</Text>
                <Text>
                  {order.shippingAddress.province}
                  {order.shippingAddress.city}
                  {order.shippingAddress.district}
                  {order.shippingAddress.streetAddress}
                </Text>
                {order.shippingAddress.postalCode && (
                  <Text>邮编: {order.shippingAddress.postalCode}</Text>
                )}
              </Space>
            ) : (
              <Text type="secondary">暂无收货地址信息</Text>
            )}
          </Card>
        </Col>
      </Row>
    </div>
  );
};