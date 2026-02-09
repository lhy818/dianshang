import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { 
  Row, 
  Col, 
  Card, 
  Typography, 
  Button, 
  Table, 
  InputNumber, 
  Space, 
  Divider, 
  Empty, 
  Checkbox,
  message,
  Popconfirm
} from 'antd';
import { 
  DeleteOutlined, 
  ShoppingCartOutlined, 
  ArrowLeftOutlined,
  ShoppingOutlined
} from '@ant-design/icons';
import { useCartStore } from '../../store/cartStore';
import { CartItem } from '../../types';

const { Title, Text } = Typography;

export const CartPage: React.FC = () => {
  const navigate = useNavigate();
  const { 
    items, 
    totalItems, 
    totalAmount, 
    updateQuantity, 
    removeItem, 
    clearCart 
  } = useCartStore();
  
  const [selectedItems, setSelectedItems] = useState<string[]>([]);
  const [isCheckingOut, setIsCheckingOut] = useState(false);

  const columns = [
    {
      title: '商品',
      dataIndex: 'product',
      key: 'product',
      width: '40%',
      render: (_: any, record: CartItem) => (
        <Space>
          <img
            src={record.product?.mainImageUrl || 'https://via.placeholder.com/80'}
            alt={record.product?.name}
            style={{ width: 80, height: 80, objectFit: 'cover' }}
          />
          <Space direction="vertical" size="small">
            <Text strong>{record.product?.name}</Text>
            <Text type="secondary">SKU: {record.product?.sku}</Text>
            {record.skuId && (
              <Text type="secondary">规格: {record.skuId}</Text>
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
        <Text strong style={{ color: '#ff4d4f' }}>
          ¥{price.toFixed(2)}
        </Text>
      ),
    },
    {
      title: '数量',
      dataIndex: 'quantity',
      key: 'quantity',
      width: '20%',
      render: (quantity: number, record: CartItem) => (
        <InputNumber
          min={1}
          max={record.product?.stockQuantity || 99}
          value={quantity}
          onChange={(value) => updateQuantity(record.id, value || 1)}
          size="middle"
        />
      ),
    },
    {
      title: '小计',
      key: 'subtotal',
      width: '15%',
      render: (_: any, record: CartItem) => (
        <Text strong style={{ color: '#ff4d4f' }}>
          ¥{(record.unitPrice * record.quantity).toFixed(2)}
        </Text>
      ),
    },
    {
      title: '操作',
      key: 'actions',
      width: '10%',
      render: (_: any, record: CartItem) => (
        <Button
          type="text"
          danger
          icon={<DeleteOutlined />}
          onClick={() => removeItem(record.id)}
        />
      ),
    },
  ];

  const handleSelectAll = (checked: boolean) => {
    if (checked) {
      setSelectedItems(items.map(item => item.id));
    } else {
      setSelectedItems([]);
    }
  };

  const handleSelectItem = (itemId: string, checked: boolean) => {
    if (checked) {
      setSelectedItems([...selectedItems, itemId]);
    } else {
      setSelectedItems(selectedItems.filter(id => id !== itemId));
    }
  };

  const handleCheckout = () => {
    if (selectedItems.length === 0) {
      message.warning('请选择要结算的商品');
      return;
    }
    
    setIsCheckingOut(true);
    // 这里可以跳转到结算页面
    navigate('/checkout', { state: { selectedItems } });
  };

  const selectedTotal = items
    .filter(item => selectedItems.includes(item.id))
    .reduce((sum, item) => sum + item.unitPrice * item.quantity, 0);

  if (items.length === 0) {
    return (
      <Card>
        <Empty
          image={<ShoppingCartOutlined style={{ fontSize: 64, color: '#d9d9d9' }} />}
          description="购物车是空的"
        >
          <Button type="primary" onClick={() => navigate('/products')}>
            去逛逛
          </Button>
        </Empty>
      </Card>
    );
  }

  return (
    <div>
      <Title level={2}>
        <ShoppingCartOutlined /> 我的购物车
      </Title>
      
      <Card>
        {/* 购物车头部 */}
        <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
          <Col>
            <Space>
              <Checkbox
                checked={selectedItems.length === items.length && items.length > 0}
                indeterminate={selectedItems.length > 0 && selectedItems.length < items.length}
                onChange={(e) => handleSelectAll(e.target.checked)}
              >
                全选
              </Checkbox>
              <Text>已选择 {selectedItems.length} 件商品</Text>
            </Space>
          </Col>
          <Col>
            <Popconfirm
              title="确定要清空购物车吗？"
              onConfirm={clearCart}
              okText="确定"
              cancelText="取消"
            >
              <Button danger icon={<DeleteOutlined />}>
                清空购物车
              </Button>
            </Popconfirm>
          </Col>
        </Row>

        <Divider />

        {/* 购物车商品列表 */}
        <Table
          dataSource={items}
          columns={columns}
          rowKey="id"
          pagination={false}
          rowSelection={{
            selectedRowKeys: selectedItems,
            onChange: (selectedRowKeys) => {
              setSelectedItems(selectedRowKeys as string[]);
            },
            getCheckboxProps: (record: CartItem) => ({
              disabled: (record.product?.stockQuantity || 0) <= 0,
            }),
          }}
          scroll={{ x: 800 }}
        />

        <Divider />

        {/* 购物车底部 */}
        <Row justify="space-between" align="middle">
          <Col>
            <Space>
              <Button 
                icon={<ArrowLeftOutlined />}
                onClick={() => navigate('/products')}
              >
                继续购物
              </Button>
              <Popconfirm
                title="确定要清空购物车吗？"
                onConfirm={clearCart}
                okText="确定"
                cancelText="取消"
              >
                <Button danger icon={<DeleteOutlined />}>
                  清空购物车
                </Button>
              </Popconfirm>
            </Space>
          </Col>
          
          <Col>
            <Space direction="vertical" align="end" size="small">
              <Space>
                <Text>已选择 {selectedItems.length} 件商品</Text>
                <Text strong style={{ fontSize: 20, color: '#ff4d4f' }}>
                  合计: ¥{selectedTotal.toFixed(2)}
                </Text>
              </Space>
              <Text type="secondary">
                购物车共 {totalItems} 件商品，总计 ¥{totalAmount.toFixed(2)}
              </Text>
              <Button
                type="primary"
                size="large"
                icon={<ShoppingOutlined />}
                loading={isCheckingOut}
                onClick={handleCheckout}
                disabled={selectedItems.length === 0}
                style={{ width: 200 }}
              >
                去结算
              </Button>
            </Space>
          </Col>
        </Row>
      </Card>

      {/* 购物车提示 */}
      <Card style={{ marginTop: 24 }}>
        <Space direction="vertical" size="small">
          <Text strong>购物须知：</Text>
          <Text type="secondary">• 商品价格可能随时变动，请以结算时价格为准</Text>
          <Text type="secondary">• 库存有限，请及时下单</Text>
          <Text type="secondary">• 支持7天无理由退货</Text>
          <Text type="secondary">• 满99元包邮</Text>
        </Space>
      </Card>
    </div>
  );
};