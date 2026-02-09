import React, { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Row, Col, Card, Typography, Button, Input, Select, Pagination, Spin, Empty, Tag, Space } from 'antd';
import { SearchOutlined, FilterOutlined, ShoppingCartOutlined, EyeOutlined } from '@ant-design/icons';
import { productApi } from '../../api/services';
import { Product, ProductQueryParams } from '../../types';
import { useCartStore } from '../../store/cartStore';

const { Title, Text } = Typography;
const { Search } = Input;
const { Option } = Select;

interface ProductListProps {
  categoryId?: string;
  showFilters?: boolean;
}

export const ProductList: React.FC<ProductListProps> = ({ 
  categoryId, 
  showFilters = true 
}) => {
  const [searchParams, setSearchParams] = useState<ProductQueryParams>({
    page: 1,
    pageSize: 12,
    sortBy: 'createdAt',
    sortOrder: 'desc',
  });
  
  const [searchText, setSearchText] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string | undefined>(categoryId);
  const [priceRange, setPriceRange] = useState<[number?, number?]>([]);
  
  const { addItem } = useCartStore();

  // 获取商品列表
  const { data, isLoading, isError } = useQuery({
    queryKey: ['products', searchParams, selectedCategory],
    queryFn: () => {
      const params = { ...searchParams };
      if (selectedCategory) {
        params.categoryId = selectedCategory;
      }
      return productApi.getProducts(params);
    },
  });

  // 获取分类列表
  const { data: categories } = useQuery({
    queryKey: ['categories'],
    queryFn: () => productApi.getCategories(),
  });

  const handleSearch = (value: string) => {
    setSearchText(value);
    setSearchParams(prev => ({
      ...prev,
      search: value,
      page: 1,
    }));
  };

  const handleSortChange = (value: string) => {
    const [sortBy, sortOrder] = value.split('_');
    setSearchParams(prev => ({
      ...prev,
      sortBy,
      sortOrder: sortOrder as 'asc' | 'desc',
      page: 1,
    }));
  };

  const handleFilterChange = (filterType: string, value: any) => {
    setSearchParams(prev => ({
      ...prev,
      [filterType]: value,
      page: 1,
    }));
  };

  const handlePageChange = (page: number, pageSize: number) => {
    setSearchParams(prev => ({
      ...prev,
      page,
      pageSize,
    }));
  };

  const handleAddToCart = (product: Product) => {
    addItem(product, 1);
    // 这里可以添加成功提示
  };

  if (isError) {
    return (
      <Card>
        <Empty description="加载商品失败，请稍后重试" />
      </Card>
    );
  }

  return (
    <div>
      {/* 搜索和筛选区域 */}
      {showFilters && (
        <Card style={{ marginBottom: 24 }}>
          <Row gutter={[16, 16]} align="middle">
            <Col xs={24} md={8}>
              <Search
                placeholder="搜索商品..."
                enterButton={<SearchOutlined />}
                size="large"
                value={searchText}
                onChange={(e) => setSearchText(e.target.value)}
                onSearch={handleSearch}
              />
            </Col>
            
            <Col xs={24} md={8}>
              <Select
                placeholder="选择分类"
                style={{ width: '100%' }}
                size="large"
                value={selectedCategory}
                onChange={setSelectedCategory}
                allowClear
              >
                {categories?.map(category => (
                  <Option key={category.id} value={category.id}>
                    {category.name}
                  </Option>
                ))}
              </Select>
            </Col>
            
            <Col xs={24} md={8}>
              <Select
                placeholder="排序方式"
                style={{ width: '100%' }}
                size="large"
                defaultValue="createdAt_desc"
                onChange={handleSortChange}
              >
                <Option value="createdAt_desc">最新上架</Option>
                <Option value="price_asc">价格从低到高</Option>
                <Option value="price_desc">价格从高到低</Option>
                <Option value="soldQuantity_desc">销量最高</Option>
                <Option value="rating_desc">评分最高</Option>
              </Select>
            </Col>
          </Row>
          
          <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
            <Col>
              <Button 
                icon={<FilterOutlined />}
                onClick={() => handleFilterChange('isRecommended', true)}
              >
                推荐商品
              </Button>
            </Col>
            <Col>
              <Button 
                icon={<FilterOutlined />}
                onClick={() => handleFilterChange('isHot', true)}
              >
                热销商品
              </Button>
            </Col>
            <Col>
              <Button 
                icon={<FilterOutlined />}
                onClick={() => handleFilterChange('isNew', true)}
              >
                新品上市
              </Button>
            </Col>
          </Row>
        </Card>
      )}

      {/* 商品列表 */}
      {isLoading ? (
        <div style={{ textAlign: 'center', padding: 50 }}>
          <Spin size="large" />
        </div>
      ) : (
        <>
          <Row gutter={[16, 24]}>
            {data?.items.map(product => (
              <Col key={product.id} xs={24} sm={12} md={8} lg={6}>
                <Card
                  hoverable
                  cover={
                    <div style={{ height: 200, overflow: 'hidden' }}>
                      <img
                        alt={product.name}
                        src={product.mainImageUrl || 'https://via.placeholder.com/300x200'}
                        style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                      />
                    </div>
                  }
                  actions={[
                    <EyeOutlined key="view" />,
                    <ShoppingCartOutlined 
                      key="cart" 
                      onClick={() => handleAddToCart(product)}
                    />,
                  ]}
                >
                  <Card.Meta
                    title={
                      <Space direction="vertical" size="small" style={{ width: '100%' }}>
                        <Text strong ellipsis={{ tooltip: product.name }}>
                          {product.name}
                        </Text>
                        <Space size="small">
                          {product.isNew && <Tag color="red">新品</Tag>}
                          {product.isHot && <Tag color="orange">热销</Tag>}
                          {product.isRecommended && <Tag color="blue">推荐</Tag>}
                        </Space>
                      </Space>
                    }
                    description={
                      <Space direction="vertical" size="small" style={{ width: '100%' }}>
                        <Text type="secondary" ellipsis>
                          {product.shortDescription || product.description?.substring(0, 50)}
                        </Text>
                        <Space>
                          <Text strong style={{ color: '#ff4d4f', fontSize: 18 }}>
                            ¥{product.price.toFixed(2)}
                          </Text>
                          {product.originalPrice && (
                            <Text delete type="secondary">
                              ¥{product.originalPrice.toFixed(2)}
                            </Text>
                          )}
                        </Space>
                        <Space>
                          <Text type="secondary">库存: {product.stockQuantity}</Text>
                          <Text type="secondary">销量: {product.soldQuantity}</Text>
                        </Space>
                        <Space>
                          <Text type="secondary">评分: {product.rating.toFixed(1)}</Text>
                          <Text type="secondary">浏览: {product.viewCount}</Text>
                        </Space>
                      </Space>
                    }
                  />
                </Card>
              </Col>
            ))}
          </Row>

          {/* 分页 */}
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
    </div>
  );
};