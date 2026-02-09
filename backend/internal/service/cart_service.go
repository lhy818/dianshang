package service

import (
	"errors"

	"gorm.io/gorm"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/repository"
)

// CartService 购物车服务
type CartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

// NewCartService 创建购物车服务
func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// GetCartByUserID 根据用户ID获取购物车
func (s *CartService) GetCartByUserID(userID string) (*domain.ShoppingCart, error) {
	if userID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	
	return s.cartRepo.FindByUserID(userID)
}

// AddCartItem 添加商品到购物车
func (s *CartService) AddCartItem(userID, productID, skuID string, quantity int) (*domain.CartItem, error) {
	if userID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	
	if productID == "" {
		return nil, errors.New("商品ID不能为空")
	}
	
	if quantity <= 0 {
		return nil, errors.New("商品数量必须大于0")
	}
	
	// 获取或创建购物车
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建购物车
			cart = &domain.ShoppingCart{
				UserID: userID,
			}
			if err := s.cartRepo.Create(cart); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	
	// 获取商品信息
	product, err := s.productRepo.FindByID(productID)
	if err != nil || product == nil {
		return nil, errors.New("商品不存在")
	}
	
	if !product.IsActive {
		return nil, errors.New("商品已下架")
	}
	
	// 检查库存
	var sku *domain.ProductSKU
	var unitPrice float64
	var stockQuantity int
	
	if skuID != "" {
		// 使用SKU
		sku, err = s.productRepo.FindSKUByID(skuID)
		if err != nil || sku == nil {
			return nil, errors.New("商品SKU不存在")
		}
		
		if !sku.IsActive {
			return nil, errors.New("商品SKU已下架")
		}
		
		if sku.ProductID != productID {
			return nil, errors.New("商品SKU与商品不匹配")
		}
		
		unitPrice = sku.Price
		stockQuantity = sku.StockQuantity
	} else {
		// 使用商品
		unitPrice = product.Price
		stockQuantity = product.StockQuantity
	}
	
	// 检查库存是否充足
	if stockQuantity < quantity {
		return nil, errors.New("商品库存不足")
	}
	
	// 检查购物车中是否已存在该商品
	existingItem, err := s.cartRepo.FindItemByProduct(cart.ID, productID, skuID)
	if err == nil && existingItem != nil {
		// 更新数量
		newQuantity := existingItem.Quantity + quantity
		if stockQuantity < newQuantity {
			return nil, errors.New("商品库存不足")
		}
		
		existingItem.Quantity = newQuantity
		if err := s.cartRepo.UpdateItem(existingItem); err != nil {
			return nil, err
		}
		
		return existingItem, nil
	}
	
	// 创建新的购物车商品项
	cartItem := &domain.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		SKUID:     skuID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
	}
	
	if err := s.cartRepo.CreateItem(cartItem); err != nil {
		return nil, err
	}
	
	return cartItem, nil
}

// UpdateCartItem 更新购物车商品
func (s *CartService) UpdateCartItem(userID, cartItemID string, quantity int) (*domain.CartItem, error) {
	if userID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	
	if cartItemID == "" {
		return nil, errors.New("购物车商品ID不能为空")
	}
	
	if quantity <= 0 {
		return nil, errors.New("商品数量必须大于0")
	}
	
	// 获取购物车商品项
	cartItem, err := s.cartRepo.FindItemByID(cartItemID)
	if err != nil || cartItem == nil {
		return nil, errors.New("购物车商品不存在")
	}
	
	// 验证购物车属于该用户
	cart, err := s.cartRepo.FindByID(cartItem.CartID)
	if err != nil || cart == nil || cart.UserID != userID {
		return nil, errors.New("购物车不属于该用户")
	}
	
	// 检查库存
	var stockQuantity int
	
	if cartItem.SKUID != "" {
		// 检查SKU库存
		sku, err := s.productRepo.FindSKUByID(cartItem.SKUID)
		if err != nil || sku == nil || !sku.IsActive {
			return nil, errors.New("商品SKU不存在或已下架")
		}
		stockQuantity = sku.StockQuantity
	} else {
		// 检查商品库存
		product, err := s.productRepo.FindByID(cartItem.ProductID)
		if err != nil || product == nil || !product.IsActive {
			return nil, errors.New("商品不存在或已下架")
		}
		stockQuantity = product.StockQuantity
	}
	
	// 检查库存是否充足
	if stockQuantity < quantity {
		return nil, errors.New("商品库存不足")
	}
	
	// 更新数量
	cartItem.Quantity = quantity
	if err := s.cartRepo.UpdateItem(cartItem); err != nil {
		return nil, err
	}
	
	return cartItem, nil
}

// DeleteCartItem 删除购物车商品
func (s *CartService) DeleteCartItem(userID, cartItemID string) error {
	if userID == "" {
		return errors.New("用户ID不能为空")
	}
	
	if cartItemID == "" {
		return errors.New("购物车商品ID不能为空")
	}
	
	// 获取购物车商品项
	cartItem, err := s.cartRepo.FindItemByID(cartItemID)
	if err != nil || cartItem == nil {
		return errors.New("购物车商品不存在")
	}
	
	// 验证购物车属于该用户
	cart, err := s.cartRepo.FindByID(cartItem.CartID)
	if err != nil || cart == nil || cart.UserID != userID {
		return errors.New("购物车不属于该用户")
	}
	
	return s.cartRepo.DeleteItem(cartItemID)
}

// ClearCart 清空购物车
func (s *CartService) ClearCart(userID string) error {
	if userID == "" {
		return errors.New("用户ID不能为空")
	}
	
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // 购物车不存在，无需清空
		}
		return err
	}
	
	return s.cartRepo.ClearItems(cart.ID)
}

// GetCartTotal 获取购物车总金额
func (s *CartService) GetCartTotal(userID string) (float64, error) {
	if userID == "" {
		return 0, errors.New("用户ID不能为空")
	}
	
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil // 购物车不存在，总金额为0
		}
		return 0, err
	}
	
	items, err := s.cartRepo.FindItemsByCartID(cart.ID)
	if err != nil {
		return 0, err
	}
	
	var total float64
	for _, item := range items {
		total += item.UnitPrice * float64(item.Quantity)
	}
	
	return total, nil
}