import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { CartItem, Product } from '../types';

interface CartState {
  items: CartItem[];
  totalItems: number;
  totalAmount: number;
  isLoading: boolean;
  
  // Actions
  addItem: (product: Product, quantity: number, skuId?: string) => void;
  updateQuantity: (itemId: string, quantity: number) => void;
  removeItem: (itemId: string) => void;
  clearCart: () => void;
  setLoading: (loading: boolean) => void;
  calculateTotals: () => void;
}

export const useCartStore = create<CartState>()(
  persist(
    (set, get) => ({
      items: [],
      totalItems: 0,
      totalAmount: 0,
      isLoading: false,
      
      addItem: (product, quantity, skuId) => {
        const { items } = get();
        const existingItem = items.find(item => 
          item.productId === product.id && item.skuId === skuId
        );
        
        if (existingItem) {
          // 更新现有商品数量
          const updatedItems = items.map(item =>
            item.id === existingItem.id
              ? { ...item, quantity: item.quantity + quantity }
              : item
          );
          set({ items: updatedItems });
        } else {
          // 添加新商品
          const newItem: CartItem = {
            id: `temp_${Date.now()}`,
            cartId: 'temp',
            productId: product.id,
            skuId,
            quantity,
            unitPrice: product.price,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            product,
          };
          set({ items: [...items, newItem] });
        }
        
        get().calculateTotals();
      },
      
      updateQuantity: (itemId, quantity) => {
        const { items } = get();
        if (quantity <= 0) {
          get().removeItem(itemId);
          return;
        }
        
        const updatedItems = items.map(item =>
          item.id === itemId ? { ...item, quantity } : item
        );
        set({ items: updatedItems });
        get().calculateTotals();
      },
      
      removeItem: (itemId) => {
        const { items } = get();
        const updatedItems = items.filter(item => item.id !== itemId);
        set({ items: updatedItems });
        get().calculateTotals();
      },
      
      clearCart: () => {
        set({ items: [], totalItems: 0, totalAmount: 0 });
      },
      
      setLoading: (isLoading) => set({ isLoading }),
      
      calculateTotals: () => {
        const { items } = get();
        const totalItems = items.reduce((sum, item) => sum + item.quantity, 0);
        const totalAmount = items.reduce(
          (sum, item) => sum + item.unitPrice * item.quantity,
          0
        );
        set({ totalItems, totalAmount });
      },
    }),
    {
      name: 'cart-storage',
      partialize: (state) => ({ items: state.items }),
    }
  )
);