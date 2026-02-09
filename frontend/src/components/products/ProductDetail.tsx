import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { 
  Row, 
  Col, 
  Card, 
  Typography, 
  Button, 
  Image, 
  Space, 
  Tag, 
  InputNumber, 
  Tabs, 
  Rate, 
  Divider, 
  Spin, 
  message,
  Carousel
} from 'antd';
import { 
  ShoppingCartOutlined, 
  HeartOutlined, 
  ShareAltOutlined,
  LeftOutlined,
  CheckCircleOutlined
} from '@ant-design/icons';
import { productApi } from '../../api/services';
import { useCartStore } from '../../store/cartStore';
import { Product } from '../../types';

const { Title, Text } = Typography;
const { TabPane } = Tabs;

export const ProductDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { addItem } = useCartStore();
  
  const [quantity, setQuantity] = useState(1);
  const [selectedImageIndex, setSelectedImageIndex] = useState(0);
  const [isAddingToCart, setIsAddingToCart] = useState(false);

  // 获取商品详情
  const { data: product, isLoading, isError } = useQuery({
    queryKey: ['product', id],
    queryFn: () => productApi.getProduct(id!),
    enabled: !!id,
  });

  // 获取商品SKU
  const { data: skus } = useQuery({
    queryKey: ['product-skus', id],
    queryFn: () => productApi.getProductSkus(id!),
    enabled: !!id,
  });

  // 获取商品评价
  const { data: reviews } = useQuery({
    queryKey: ['product-reviews', id],
    queryFn: () => productApi.getProductReviews(id!),
    enabled: !!id,
  });

  const handleAddToCart = async () => {
    if (!product) return;
    
    try {
      setIsAddingToCart(true);
      addItem(product, quantity);
      message.success('已添加到购物车');
    } catch (error) {
      message.error('添加到购物车失败');
    } finally {
      setIsAddingToCart(false);
    }
  };

  const handleBuyNow = async () => {
    if (!product) return;
    
    try {
      setIsAddingToCart(true);
      addItem(product, quantity);
      message.success('已添加到购物车');
      navigate('/cart');
    } catch (error) {
      message.error('操作失败');
    } finally {
      setIsAddingToCart(false);
    }
  };

  if (isLoading) {
    return (
      <div style={{ textAlign: 'center', padding: 50 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (isError || !product) {
    return (
      <Card>
        <div style={{ textAlign: 'center', padding: 50 }}>
          <Title level={3}>商品不存在或加载失败</Title>
          <Button type="primary" onClick={() => navigate('/products')}>
            返回商品列表
          </Button>
        </div>
      </Card>
    );
  }

  const images = product.imageUrls.length > 0 
    ? product.imageUrls 
    : product.mainImageUrl 
      ? [product.mainImageUrl]
      : ['https://via.placeholder.com/600x400'];

  return (
    <div>
      <Button 
        type="link" 
        icon={<LeftOutlined />}
        onClick={() => navigate(-1)}
        style={{ marginBottom: 16 }}
      >
        返回
      </Button>

      <Card>
        <Row gutter={[32, 32]}>
          {/* 商品图片 */}
          <Col xs={24} md={12}>
            <div style={{ marginBottom: 16 }}>
              <Image
                src={images[selectedImageIndex]}
                alt={product.name}
                style={{ width: '100%', maxHeight: 500, objectFit: 'contain' }}
                preview={false}
              />
            </div>
            
            {images.length > 1 && (
              <Carousel dots={{ className: 'custom-dots' }} autoplay>
                {images.map((img, index) => (
                  <div key={index}>
                    <Image
                      src={img}
                      alt={`${product.name} - ${index + 1}`}
                      style={{ width: '100%', height: 100, objectFit: 'cover' }}
                      preview={false}
                      onClick={() => setSelectedImageIndex(index)}
                    />
                  </div>
                ))}
              </Carousel>
            )}
          </Col>

          {/* 商品信息 */}
          <Col xs={24} md={12}>
            <Space direction="vertical" size="large" style={{ width: '100%' }}>
              <div>
                <Title level={2}>{product.name}</Title>
                <Space size="middle">
                  <Text type="secondary">SKU: {product.sku}</Text>
                  <Space>
                    {product.isNew && <Tag color="red">新品</Tag>}
                    {product.isHot && <Tag color="orange">热销</Tag>}
                    {product.isRecommended && <Tag color="blue">推荐</Tag>}
                  </Space>
                </Space>
              </div>

              <div>
                <Space align="baseline" size="large">
                  <Title level={3} style={{ color: '#ff4d4f', margin: 0 }}>
                    ¥{product.price.toFixed(2)}
                  </Title>
                  {product.originalPrice && (
                    <Text delete type="secondary" style={{ fontSize: 16 }}>
                      ¥{product.originalPrice.toFixed(2)}
                    </Text>
                  )}
                </Space>
              </div>

              <div>
                <Space direction="vertical" size="small">
                  <Text strong>商品评分</Text>
                  <Space>
                    <Rate disabled defaultValue={product.rating} allowHalf />
                    <Text>({product.reviewCount} 条评价)</Text>
                  </Space>
                </Space>
              </div>

              <Divider />

              {/* 库存信息 */}
              <div>
                <Space direction="vertical" size="small">
                  <Text strong>库存状态</Text>
                  {product.stockQuantity > 0 ? (
                    <Space>
                      <CheckCircleOutlined style={{ color: '#52c41a' }} />
                      <Text type="success">有货</Text>
                      <Text type="secondary">剩余 {product.stockQuantity} 件</Text>
                    </Space>
                  ) : (
                    <Text type="danger">缺货</Text>
                  )}
                </Space>
              </div>

              {/* 数量选择 */}
              <div>
                <Space direction="vertical" size="small">
                  <Text strong>购买数量</Text>
                  <InputNumber
                    min={1}
                    max={product.stockQuantity}
                    value={quantity}
                    onChange={setQuantity}
                    size="large"
                    style={{ width: 120 }}
                  />
                </Space>
              </div>

              <Divider />

              {/* 操作按钮 */}
              <Space size="large">
                <Button
                  type="primary"
                  size="large"
                  icon={<ShoppingCartOutlined />}
                  loading={isAddingToCart}
                  onClick={handleAddToCart}
                  disabled={product.stockQuantity <= 0}
                >
                  加入购物车
                </Button>
                
                <Button
                  type="primary"
                  size="large"
                  danger
                  onClick={handleBuyNow}
                  disabled={product.stockQuantity <= 0}
                >
                  立即购买
                </Button>
                
                <Button size="large" icon={<HeartOutlined />}>
                  收藏
                </Button>
                
                <Button size="large" icon={<ShareAltOutlined />}>
                  分享
                </Button>
              </Space>

              {/* 商品属性 */}
              {Object.keys(product.attributes).length > 0 && (
                <div>
                  <Text strong>商品属性</Text>
                  <Row gutter={[16, 8]} style={{ marginTop: 8 }}>
                    {Object.entries(product.attributes).map(([key, value]) => (
                      <Col span={12} key={key}>
                        <Text type="secondary">{key}: </Text>
                        <Text strong>{value}</Text>
                      </Col>
                    ))}
                  </Row>
                </div>
              )}
            </Space>
          </Col>
        </Row>

        {/* 商品详情标签页 */}
        <Divider />
        
        <Tabs defaultActiveKey="description">
          <TabPane tab="商品详情" key="description">
            <div style={{ padding: 24 }}>
              {product.description ? (
                <div dangerouslySetInnerHTML={{ __html: product.description }} />
              ) : (
                <Text type="secondary">暂无商品详情描述</Text>
              )}
            </div>
          </TabPane>
          
          <TabPane tab="规格参数" key="specifications">
            <div style={{ padding: 24 }}>
              {Object.keys(product.attributes).length > 0 ? (
                <Row gutter={[16, 16]}>
                  {Object.entries(product.attributes).map(([key, value]) => (
                    <Col span={12} key={key}>
                      <Card size="small">
                        <Text strong>{key}: </Text>
                        <Text>{value}</Text>
                      </Card>
                    </Col>
                  ))}
                </Row>
              ) : (
                <Text type="secondary">暂无规格参数</Text>
              )}
            </div>
          </TabPane>
          
          <TabPane tab="商品评价" key="reviews">
            <div style={{ padding: 24 }}>
              {reviews && reviews.length > 0 ? (
                <Space direction="vertical" size="large" style={{ width: '100%' }}>
                  {reviews.map((review: any) => (
                    <Card key={review.id} size="small">
                      <Space direction="vertical" size="small">
                        <Space>
                          <Rate disabled defaultValue={review.rating} />
                          <Text strong>{review.userName}</Text>
                          <Text type="secondary">
                            {new Date(review.createdAt).toLocaleDateString()}
                          </Text>
                        </Space>
                        <Text>{review.content}</Text>
                      </Space>
                    </Card>
                  ))}
                </Space>
              ) : (
                <Text type="secondary">暂无评价</Text>
              )}
            </div>
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
};