// 基础类型
export interface BaseEntity {
  id: string;
  createdAt: string;
  updatedAt: string;
}

// 用户相关
export interface User extends BaseEntity {
  username: string;
  email: string;
  phone?: string;
  avatarUrl?: string;
  status: 'active' | 'inactive' | 'banned';
  isEmailVerified: boolean;
  isPhoneVerified: boolean;
  lastLoginAt?: string;
}

export interface UserProfile extends BaseEntity {
  userId: string;
  realName?: string;
  gender?: 'male' | 'female' | 'other';
  birthDate?: string;
  idCardNumber?: string;
}

export interface UserAddress extends BaseEntity {
  userId: string;
  recipientName: string;
  phone: string;
  province: string;
  city: string;
  district: string;
  streetAddress: string;
  postalCode?: string;
  isDefault: boolean;
}

// 商品相关
export interface Category extends BaseEntity {
  parentId?: string;
  name: string;
  slug: string;
  description?: string;
  imageUrl?: string;
  sortOrder: number;
  isActive: boolean;
}

export interface Product extends BaseEntity {
  categoryId?: string;
  name: string;
  slug: string;
  description?: string;
  shortDescription?: string;
  sku: string;
  price: number;
  originalPrice?: number;
  costPrice?: number;
  stockQuantity: number;
  soldQuantity: number;
  weight?: number;
  volume?: number;
  mainImageUrl?: string;
  imageUrls: string[];
  attributes: Record<string, any>;
  isActive: boolean;
  isRecommended: boolean;
  isHot: boolean;
  isNew: boolean;
  viewCount: number;
  rating: number;
  reviewCount: number;
}

export interface ProductSku extends BaseEntity {
  productId: string;
  skuCode: string;
  attributes: Record<string, any>;
  price: number;
  originalPrice?: number;
  stockQuantity: number;
  soldQuantity: number;
  imageUrl?: string;
  isActive: boolean;
}

// 购物车相关
export interface ShoppingCart extends BaseEntity {
  userId: string;
  sessionId?: string;
}

export interface CartItem extends BaseEntity {
  cartId: string;
  productId: string;
  skuId?: string;
  quantity: number;
  unitPrice: number;
  product?: Product;
}

// 订单相关
export interface Order extends BaseEntity {
  orderNo: string;
  userId: string;
  totalAmount: number;
  discountAmount: number;
  shippingFee: number;
  finalAmount: number;
  status: 'pending' | 'paid' | 'shipped' | 'delivered' | 'completed' | 'cancelled' | 'refunded';
  paymentStatus: 'unpaid' | 'paid' | 'refunded' | 'failed';
  paymentMethod?: string;
  paymentTime?: string;
  shippingAddress: Record<string, any>;
  shippingMethod?: string;
  shippingNo?: string;
  shippingTime?: string;
  deliveredTime?: string;
  buyerRemark?: string;
  sellerRemark?: string;
  cancelledReason?: string;
  cancelledTime?: string;
}

export interface OrderItem extends BaseEntity {
  orderId: string;
  productId: string;
  skuId?: string;
  productName: string;
  skuAttributes: Record<string, any>;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
}

// 支付相关
export interface Payment extends BaseEntity {
  orderId: string;
  paymentNo: string;
  paymentMethod: string;
  paymentChannel?: string;
  amount: number;
  status: 'pending' | 'success' | 'failed' | 'refunded';
  transactionId?: string;
  payerInfo: Record<string, any>;
  rawResponse: Record<string, any>;
  paidAt?: string;
}

// 请求/响应类型
export interface PaginatedResponse<T> {
  items: T[];
  pagination: {
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
  };
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  phone?: string;
}

export interface CreateOrderRequest {
  cartItemIds: string[];
  shippingAddressId: string;
  paymentMethod: string;
  buyerRemark?: string;
}

// API响应类型
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
  timestamp: number;
}

export interface ApiError {
  code: number;
  message: string;
  errors?: Array<{
    field: string;
    message: string;
  }>;
  timestamp: number;
}

// 分页查询参数
export interface PaginationParams {
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
}

// 商品查询参数
export interface ProductQueryParams extends PaginationParams {
  categoryId?: string;
  search?: string;
  minPrice?: number;
  maxPrice?: number;
  isRecommended?: boolean;
  isHot?: boolean;
  isNew?: boolean;
}