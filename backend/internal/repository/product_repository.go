package repository

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"ecommerce-backend/internal/domain"
)

// ProductRepository 商品仓库接口
type ProductRepository interface {
	// 分类相关
	CreateCategory(category *domain.Category) error
	FindCategoryByID(id string) (*domain.Category, error)
	FindCategoryBySlug(slug string) (*domain.Category, error)
	FindCategories() ([]*domain.Category, error)
	UpdateCategory(category *domain.Category) error
	DeleteCategory(id string) error

	// 商品相关
	Create(product *domain.Product) error
	FindByID(id string) (*domain.Product, error)
	FindBySKU(sku string) (*domain.Product, error)
	FindBySlug(slug string) (*domain.Product, error)
	FindProducts(query string, conditions map[string]interface{}, priceRange map[string]interface{}, page, pageSize int, sortBy, sortOrder string) ([]*domain.Product, int64, error)
	Update(product *domain.Product) error
	Delete(id string) error
	UpdateStock(productID string, quantity int) error
	UpdateStockWithTx(tx *gorm.DB, productID string, quantity int) error
	IncrementViewCount(productID string) error
	UpdateRating(productID string, rating float64, reviewCount int) error

	// SKU相关
	CreateSKU(sku *domain.ProductSKU) error
	FindSKUByID(id string) (*domain.ProductSKU, error)
	FindSKUsByProductID(productID string) ([]*domain.ProductSKU, error)
	UpdateSKU(sku *domain.ProductSKU) error
	DeleteSKU(id string) error
	UpdateSKUStock(skuID string, quantity int) error
	UpdateSKUStockWithTx(tx *gorm.DB, skuID string, quantity int) error

	// 评价相关
	CreateReview(review *domain.ProductReview) error
	FindReviewByID(id string) (*domain.ProductReview, error)
	FindReviewByOrderItemID(orderItemID string) (*domain.ProductReview, error)
	FindReviewsByProductID(productID string, page, pageSize int) ([]*domain.ProductReview, int64, error)
	FindAllReviewsByProductID(productID string) ([]*domain.ProductReview, error)
	UpdateReview(review *domain.ProductReview) error
	DeleteReview(id string) error
}

// GormProductRepository GORM实现的商品仓库
type GormProductRepository struct {
	db *gorm.DB
}

// NewGormProductRepository 创建GORM商品仓库
func NewGormProductRepository(db *gorm.DB) ProductRepository {
	return &GormProductRepository{db: db}
}

// CreateCategory 创建分类
func (r *GormProductRepository) CreateCategory(category *domain.Category) error {
	return r.db.Create(category).Error
}

// FindCategoryByID 根据ID查找分类
func (r *GormProductRepository) FindCategoryByID(id string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

// FindCategoryBySlug 根据Slug查找分类
func (r *GormProductRepository) FindCategoryBySlug(slug string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

// FindCategories 查找所有分类
func (r *GormProductRepository) FindCategories() ([]*domain.Category, error) {
	var categories []*domain.Category
	err := r.db.Where("is_active = true").Order("sort_order ASC, created_at ASC").Find(&categories).Error
	return categories, err
}

// UpdateCategory 更新分类
func (r *GormProductRepository) UpdateCategory(category *domain.Category) error {
	return r.db.Save(category).Error
}

// DeleteCategory 删除分类
func (r *GormProductRepository) DeleteCategory(id string) error {
	return r.db.Delete(&domain.Category{}, "id = ?", id).Error
}

// Create 创建商品
func (r *GormProductRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

// FindByID 根据ID查找商品
func (r *GormProductRepository) FindByID(id string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Category").Preload("SKUs", "is_active = true").
		Where("id = ? AND deleted_at IS NULL", id).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, err
}

// FindBySKU 根据SKU查找商品
func (r *GormProductRepository) FindBySKU(sku string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Where("sku = ? AND deleted_at IS NULL", sku).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, err
}

// FindBySlug 根据Slug查找商品
func (r *GormProductRepository) FindBySlug(slug string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Category").Preload("SKUs", "is_active = true").
		Where("slug = ? AND deleted_at IS NULL", slug).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, err
}

// FindProducts 查找商品列表
func (r *GormProductRepository) FindProducts(query string, conditions map[string]interface{}, priceRange map[string]interface{}, page, pageSize int, sortBy, sortOrder string) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64

	// 构建查询
	dbQuery := r.db.Model(&domain.Product{}).Where("deleted_at IS NULL")

	// 关键词搜索
	if query != "" {
		searchQuery := fmt.Sprintf("%%%s%%", strings.TrimSpace(query))
		dbQuery = dbQuery.Where("name LIKE ? OR description LIKE ? OR short_description LIKE ?", 
			searchQuery, searchQuery, searchQuery)
	}

	// 应用条件
	for key, value := range conditions {
		dbQuery = dbQuery.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 价格范围
	if minPrice, ok := priceRange["min"]; ok {
		dbQuery = dbQuery.Where("price >= ?", minPrice)
	}
	if maxPrice, ok := priceRange["max"]; ok {
		dbQuery = dbQuery.Where("price <= ?", maxPrice)
	}

	// 计算总数
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderBy := fmt.Sprintf("%s %s", sortBy, sortOrder)
	if sortBy == "" {
		orderBy = "created_at DESC"
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := dbQuery.Preload("Category").
		Offset(offset).Limit(pageSize).
		Order(orderBy).
		Find(&products).Error

	return products, total, err
}

// Update 更新商品
func (r *GormProductRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

// Delete 删除商品（软删除）
func (r *GormProductRepository) Delete(id string) error {
	return r.db.Model(&domain.Product{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

// UpdateStock 更新商品库存
func (r *GormProductRepository) UpdateStock(productID string, quantity int) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", productID).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", quantity)).Error
}

// UpdateStockWithTx 在事务中更新商品库存
func (r *GormProductRepository) UpdateStockWithTx(tx *gorm.DB, productID string, quantity int) error {
	return tx.Model(&domain.Product{}).
		Where("id = ?", productID).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", quantity)).Error
}

// IncrementViewCount 增加商品浏览量
func (r *GormProductRepository) IncrementViewCount(productID string) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", productID).
		Update("view_count", gorm.Expr("view_count + 1")).Error
}

// UpdateRating 更新商品评分
func (r *GormProductRepository) UpdateRating(productID string, rating float64, reviewCount int) error {
	return r.db.Model(&domain.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"rating":       rating,
			"review_count": reviewCount,
		}).Error
}

// CreateSKU 创建SKU
func (r *GormProductRepository) CreateSKU(sku *domain.ProductSKU) error {
	return r.db.Create(sku).Error
}

// FindSKUByID 根据ID查找SKU
func (r *GormProductRepository) FindSKUByID(id string) (*domain.ProductSKU, error) {
	var sku domain.ProductSKU
	err := r.db.Where("id = ?", id).First(&sku).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &sku, err
}

// FindSKUsByProductID 根据商品ID查找SKU列表
func (r *GormProductRepository) FindSKUsByProductID(productID string) ([]*domain.ProductSKU, error) {
	var skus []*domain.ProductSKU
	err := r.db.Where("product_id = ? AND is_active = true", productID).Find(&skus).Error
	return skus, err
}

// UpdateSKU 更新SKU
func (r *GormProductRepository) UpdateSKU(sku *domain.ProductSKU) error {
	return r.db.Save(sku).Error
}

// DeleteSKU 删除SKU
func (r *GormProductRepository) DeleteSKU(id string) error {
	return r.db.Delete(&domain.ProductSKU{}, "id = ?", id).Error
}

// UpdateSKUStock 更新SKU库存
func (r *GormProductRepository) UpdateSKUStock(skuID string, quantity int) error {
	return r.db.Model(&domain.ProductSKU{}).
		Where("id = ?", skuID).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", quantity)).Error
}

// UpdateSKUStockWithTx 在事务中更新SKU库存
func (r *GormProductRepository) UpdateSKUStockWithTx(tx *gorm.DB, skuID string, quantity int) error {
	return tx.Model(&domain.ProductSKU{}).
		Where("id = ?", skuID).
		Update("stock_quantity", gorm.Expr("stock_quantity + ?", quantity)).Error
}

// CreateReview 创建评价
func (r *GormProductRepository) CreateReview(review *domain.ProductReview) error {
	return r.db.Create(review).Error
}

// FindReviewByID 根据ID查找评价
func (r *GormProductRepository) FindReviewByID(id string) (*domain.ProductReview, error) {
	var review domain.ProductReview
	err := r.db.Preload("User").Where("id = ?", id).First(&review).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &review, err
}

// FindReviewByOrderItemID 根据订单商品ID查找评价
func (r *GormProductRepository) FindReviewByOrderItemID(orderItemID string) (*domain.ProductReview, error) {
	var review domain.ProductReview
	err := r.db.Where("order_item_id = ?", orderItemID).First(&review).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &review, err
}

// FindReviewsByProductID 根据商品ID查找评价列表
func (r *GormProductRepository) FindReviewsByProductID(productID string, page, pageSize int) ([]*domain.ProductReview, int64, error) {
	var reviews []*domain.ProductReview
	var total int64

	query := r.db.Model(&domain.ProductReview{}).
		Where("product_id = ? AND status = ?", productID, domain.ReviewStatusPublished).
		Preload("User")

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&reviews).Error

	return reviews, total, err
}

// FindAllReviewsByProductID 根据商品ID查找所有评价
func (r *GormProductRepository) FindAllReviewsByProductID(productID string) ([]*domain.ProductReview, error) {
	var reviews []*domain.ProductReview
	err := r.db.Where("product_id = ?", productID).Find(&reviews).Error
	return reviews, err
}

// UpdateReview 更新评价
func (r *GormProductRepository) UpdateReview(review *domain.ProductReview) error {
	return r.db.Save(review).Error
}

// DeleteReview 删除评价
func (r *GormProductRepository) DeleteReview(id string) error {
	return r.db.Delete(&domain.ProductReview{}, "id = ?", id).Error
}