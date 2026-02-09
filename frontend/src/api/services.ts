import { apiClient } from './client';
import {
  User,
  LoginRequest,
  RegisterRequest,
  Product,
  Category,
  CartItem,
  Order,
  CreateOrderRequest,
  UserAddress,
  PaginatedResponse,
  ProductQueryParams,
} from '../types';

// 认证相关API
export const authApi = {
  login: (data: LoginRequest) => apiClient.post<{ user: User; token: string }>('/auth/login', data),
  register: (data: RegisterRequest) => apiClient.post<{ user: User; token: string }>('/auth/register', data),
  logout: () => apiClient.post('/auth/logout'),
  refreshToken: () => apiClient.post<{ token: string }>('/auth/refresh'),
  getProfile: () => apiClient.get<User>('/auth/profile'),
  updateProfile: (data: Partial<User>) => apiClient.put<User>('/auth/profile', data),
};

// 用户相关API
export const userApi = {
  getAddresses: () => apiClient.get<UserAddress[]>('/users/addresses'),
  addAddress: (data: Omit<UserAddress, 'id' | 'createdAt' | 'updatedAt'>) => 
    apiClient.post<UserAddress>('/users/addresses', data),
  updateAddress: (id: string, data: Partial<UserAddress>) => 
    apiClient.put<UserAddress>(`/users/addresses/${id}`, data),
  deleteAddress: (id: string) => apiClient.delete(`/users/addresses/${id}`),
  setDefaultAddress: (id: string) => apiClient.put(`/users/addresses/${id}/default`),
};

// 商品相关API
export const productApi = {
  getCategories: () => apiClient.get<Category[]>('/categories'),
  getCategoryProducts: (categoryId: string, params?: ProductQueryParams) => 
    apiClient.get<PaginatedResponse<Product>>(`/categories/${categoryId}/products`, { params }),
  getProducts: (params?: ProductQueryParams) => 
    apiClient.get<PaginatedResponse<Product>>('/products', { params }),
  getProduct: (id: string) => apiClient.get<Product>(`/products/${id}`),
  getProductSkus: (id: string) => apiClient.get(`/products/${id}/skus`),
  getProductReviews: (id: string) => apiClient.get(`/products/${id}/reviews`),
};

// 购物车相关API
export const cartApi = {
  getCart: () => apiClient.get<{ items: CartItem[]; total: number }>('/cart'),
  addToCart: (data: { productId: string; skuId?: string; quantity: number }) => 
    apiClient.post<CartItem>('/cart/items', data),
  updateCartItem: (id: string, data: { quantity: number }) => 
    apiClient.put<CartItem>(`/cart/items/${id}`, data),
  removeCartItem: (id: string) => apiClient.delete(`/cart/items/${id}`),
  clearCart: () => apiClient.delete('/cart'),
};

// 订单相关API
export const orderApi = {
  getOrders: (params?: { page?: number; pageSize?: number }) => 
    apiClient.get<PaginatedResponse<Order>>('/orders', { params }),
  createOrder: (data: CreateOrderRequest) => apiClient.post<Order>('/orders', data),
  getOrder: (id: string) => apiClient.get<Order>(`/orders/${id}`),
  cancelOrder: (id: string) => apiClient.put(`/orders/${id}/cancel`),
  payOrder: (id: string) => apiClient.post(`/orders/${id}/pay`),
};