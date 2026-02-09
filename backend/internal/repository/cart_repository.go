package repository

import (
	"errors"

	"gorm.io/gorm"

	"ecommerce-backend/internal/domain"
)

// CartRepository 购物车仓库接口
type CartRepository interface {
	// 购物车相关
	Create(cart *domain.ShoppingCart) error
	FindByID(id string) (*domain.ShoppingCart, error)
	FindByUserID(userID string) (*domain.ShoppingCart, error)
	Update(cart *domain.ShoppingCart) error
	Delete(id string) error

	// 购物车商品项相关
	CreateItem(item *domain.CartItem) error
	FindItemByID(id string) (*domain.CartItem, error)
	FindItemByProduct(cartID, productID, skuID string) (*domain.CartItem, error)
	FindItemsByCartID(cartID string) ([]*domain.CartItem, error)
	UpdateItem(item *domain.CartItem) error
	DeleteItem(id string) error
	ClearItems(cartID string) error
	ClearItemsWithTx(tx *gorm.DB, cartID string) error
}

// GormCartRepository GORM实现的购物车仓库
type GormCartRepository struct {
	db *gorm.DB
}

// NewGormCartRepository 创建GORM购物车仓库
func NewGormCartRepository(db *gorm.DB) CartRepository {
	return &GormCartRepository{db: db}
}

// Create 创建购物车
func (r *GormCartRepository) Create(cart *domain.ShoppingCart) error {
	return r.db.Create(cart).Error
}

// FindByID 根据ID查找购物车
func (r *GormCartRepository) FindByID(id string) (*domain.ShoppingCart, error) {
	var cart domain.ShoppingCart
	err := r.db.Preload("Items").Preload("Items.Product").Preload("Items.SKU").
		Where("id = ?", id).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &cart, err
}

// FindByUserID 根据用户ID查找购物车
func (r *GormCartRepository) FindByUserID(userID string) (*domain.ShoppingCart, error) {
	var cart domain.ShoppingCart
	err := r.db.Preload("Items").Preload("Items.Product").Preload("Items.SKU").
		Where("user_id = ?", userID).First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &cart, err
}

// Update 更新购物车
func (r *GormCartRepository) Update(cart *domain.ShoppingCart) error {
	return r.db.Save(cart).Error
}

// Delete 删除购物车
func (r *GormCartRepository) Delete(id string) error {
	return r.db.Delete(&domain.ShoppingCart{}, "id = ?", id).Error
}

// CreateItem 创建购物车商品项
func (r *GormCartRepository) CreateItem(item *domain.CartItem) error {
	return r.db.Create(item).Error
}

// FindItemByID 根据ID查找购物车商品项
func (r *GormCartRepository) FindItemByID(id string) (*domain.CartItem, error) {
	var item domain.CartItem
	err := r.db.Preload("Product").Preload("SKU").
		Where("id = ?", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

// FindItemByProduct 根据商品查找购物车商品项
func (r *GormCartRepository) FindItemByProduct(cartID, productID, skuID string) (*domain.CartItem, error) {
	var item domain.CartItem
	
	query := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID)
	if skuID != "" {
		query = query.Where("sku_id = ?", skuID)
	} else {
		query = query.Where("sku_id IS NULL OR sku_id = ''")
	}
	
	err := query.First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

// FindItemsByCartID 根据购物车ID查找所有商品项
func (r *GormCartRepository) FindItemsByCartID(cartID string) ([]*domain.CartItem, error) {
	var items []*domain.CartItem
	err := r.db.Preload("Product").Preload("SKU").
		Where("cart_id = ?", cartID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

// UpdateItem 更新购物车商品项
func (r *GormCartRepository) UpdateItem(item *domain.CartItem) error {
	return r.db.Save(item).Error
}

// DeleteItem 删除购物车商品项
func (r *GormCartRepository) DeleteItem(id string) error {
	return r.db.Delete(&domain.CartItem{}, "id = ?", id).Error
}

// ClearItems 清空购物车商品项
func (r *GormCartRepository) ClearItems(cartID string) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&domain.CartItem{}).Error
}

// ClearItemsWithTx 在事务中清空购物车商品项
func (r *GormCartRepository) ClearItemsWithTx(tx *gorm.DB, cartID string) error {
	return tx.Where("cart_id = ?", cartID).Delete(&domain.CartItem{}).Error
}