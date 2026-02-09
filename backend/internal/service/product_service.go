package service

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/repository"
)

// ProductService 商品服务
type ProductService struct {
	productRepo repository.ProductRepository
}

// NewProductService 创建商品服务
func NewProductService(productRepo repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// ProductFilter 商品查询过滤器
type ProductFilter struct {
	Query         string
	CategoryID    string
	MinPrice      *float64
	MaxPrice      *float64
	IsActive      *bool
	IsRecommended *bool
	IsHot         *bool
	IsNew         *bool
	Page          int
	PageSize      int
	SortBy        string
	SortOrder     string
}

// GetCategories 获取分类列表
func (s *ProductService) GetCategories() ([]*domain.Category, error) {
	return s.productRepo.FindCategories()
}

// GetProducts 获取商品列表
func (s *ProductService) GetProducts(filter *ProductFilter) ([]*domain.Product, int64, error) {
	// 构建查询条件
	conditions := make(map[string]interface{})
	
	if filter.CategoryID != "" {
		conditions["category_id"] = filter.CategoryID
	}
	
	if filter.IsActive != nil {
		conditions["is_active"] = *filter.IsActive
	}
	
	if filter.IsRecommended != nil {
		conditions["is_recommended"] = *filter.IsRecommended
	}
	
	if filter.IsHot != nil {
		conditions["is_hot"] = *filter.IsHot
	}
	
	if filter.IsNew != nil {
		conditions["is_new"] = *filter.IsNew
	}
	
	// 价格范围
	priceRange := make(map[string]interface{})
	if filter.MinPrice != nil {
		priceRange["min"] = *filter.MinPrice
	}
	if filter.MaxPrice != nil {
		priceRange["max"] = *filter.MaxPrice
	}
	
	// 排序
	sortBy := "created_at"
	sortOrder := "DESC"
	
	if filter.SortBy != "" {
		validSortFields := map[string]bool{
			"price":          true,
			"created_at":     true,
			"sold_quantity":  true,
			"rating":         true,
			"view_count":     true,
		}
		if validSortFields[filter.SortBy] {
			sortBy = filter.SortBy
		}
	}
	
	if filter.SortOrder == "ASC" {
		sortOrder = "ASC"
	}
	
	return s.productRepo.FindProducts(
		filter.Query,
		conditions,
		priceRange,
		filter.Page,
		filter.PageSize,
		sortBy,
		sortOrder,
	)
}

// GetProductByID 根据ID获取商品
func (s *ProductService) GetProductByID(id string) (*domain.Product, error) {
	if id == "" {
		return nil, errors.New("商品ID不能为空")
	}
	
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	
	// 增加浏览量
	if product != nil {
		go s.incrementViewCount(id)
	}
	
	return product, nil
}

// GetProductSKUs 获取商品SKU列表
func (s *ProductService) GetProductSKUs(productID string) ([]*domain.ProductSKU, error) {
	if productID == "" {
		return nil, errors.New("商品ID不能为空")
	}
	
	return s.productRepo.FindSKUsByProductID(productID)
}

// CreateProduct 创建商品
func (s *ProductService) CreateProduct(product *domain.Product) error {
	// 验证商品数据
	if err := s.validateProduct(product); err != nil {
		return err
	}
	
	// 检查SKU是否已存在
	if product.SKU != "" {
		existingProduct, err := s.productRepo.FindBySKU(product.SKU)
		if err == nil && existingProduct != nil {
			return errors.New("商品SKU已存在")
		}
	}
	
	return s.productRepo.Create(product)
}

// UpdateProduct 更新商品
func (s *ProductService) UpdateProduct(product *domain.Product) error {
	if product.ID == "" {
		return errors.New("商品ID不能为空")
	}
	
	// 验证商品数据
	if err := s.validateProduct(product); err != nil {
		return err
	}
	
	// 检查SKU是否与其他商品冲突
	if product.SKU != "" {
		existingProduct, err := s.productRepo.FindBySKU(product.SKU)
		if err == nil && existingProduct != nil && existingProduct.ID != product.ID {
			return errors.New("商品SKU已存在")
		}
	}
	
	return s.productRepo.Update(product)
}

// DeleteProduct 删除商品（软删除）
func (s *ProductService) DeleteProduct(id string) error {
	return s.productRepo.Delete(id)
}

// UpdateStock 更新商品库存
func (s *ProductService) UpdateStock(productID string, quantity int) error {
	if productID == "" {
		return errors.New("商品ID不能为空")
	}
	
	return s.productRepo.UpdateStock(productID, quantity)
}

// UpdateSKUStock 更新SKU库存
func (s *ProductService) UpdateSKUStock(skuID string, quantity int) error {
	if skuID == "" {
		return errors.New("SKU ID不能为空")
	}
	
	return s.productRepo.UpdateSKUStock(skuID, quantity)
}

// CreateProductReview 创建商品评价
func (s *ProductService) CreateProductReview(review *domain.ProductReview) error {
	// 验证评价数据
	if err := s.validateReview(review); err != nil {
		return err
	}
	
	// 检查是否已评价
	existingReview, err := s.productRepo.FindReviewByOrderItemID(review.OrderItemID)
	if err == nil && existingReview != nil {
		return errors.New("该订单商品已评价")
	}
	
	// 创建评价
	err = s.productRepo.CreateReview(review)
	if err != nil {
		return err
	}
	
	// 更新商品评分
	go s.updateProductRating(review.ProductID)
	
	return nil
}

// GetProductReviews 获取商品评价列表
func (s *ProductService) GetProductReviews(productID string, page, pageSize int) ([]*domain.ProductReview, int64, error) {
	if productID == "" {
		return nil, 0, errors.New("商品ID不能为空")
	}
	
	return s.productRepo.FindReviewsByProductID(productID, page, pageSize)
}

// validateProduct 验证商品数据
func (s *ProductService) validateProduct(product *domain.Product) error {
	if product.Name == "" {
		return errors.New("商品名称不能为空")
	}
	
	if len(product.Name) > 200 {
		return errors.New("商品名称不能超过200个字符")
	}
	
	if product.SKU == "" {
		return errors.New("商品SKU不能为空")
	}
	
	if len(product.SKU) > 50 {
		return errors.New("商品SKU不能超过50个字符")
	}
	
	if product.Price < 0 {
		return errors.New("商品价格不能为负数")
	}
	
	if product.OriginalPrice != nil && *product.OriginalPrice < 0 {
		return errors.New("商品原价不能为负数")
	}
	
	if product.CostPrice != nil && *product.CostPrice < 0 {
		return errors.New("商品成本价不能为负数")
	}
	
	if product.StockQuantity < 0 {
		return errors.New("商品库存不能为负数")
	}
	
	if product.SoldQuantity < 0 {
		return errors.New("商品销量不能为负数")
	}
	
	// 验证状态
	if product.IsActive && product.CategoryID != "" {
		// 检查分类是否存在且激活
		category, err := s.productRepo.FindCategoryByID(product.CategoryID)
		if err != nil || category == nil || !category.IsActive {
			return errors.New("商品分类不存在或未激活")
		}
	}
	
	return nil
}

// validateReview 验证评价数据
func (s *ProductService) validateReview(review *domain.ProductReview) error {
	if review.OrderItemID == "" {
		return errors.New("订单商品ID不能为空")
	}
	
	if review.UserID == "" {
		return errors.New("用户ID不能为空")
	}
	
	if review.ProductID == "" {
		return errors.New("商品ID不能为空")
	}
	
	if review.Rating < 1 || review.Rating > 5 {
		return errors.New("评分必须在1-5之间")
	}
	
	if review.Title != "" && len(review.Title) > 200 {
		return errors.New("评价标题不能超过200个字符")
	}
	
	// 验证状态
	if review.Status != "" {
		validStatuses := map[string]bool{
			domain.ReviewStatusPending:   true,
			domain.ReviewStatusPublished: true,
			domain.ReviewStatusHidden:    true,
		}
		if !validStatuses[review.Status] {
			return errors.New("评价状态无效")
		}
	}
	
	return nil
}

// incrementViewCount 增加商品浏览量
func (s *ProductService) incrementViewCount(productID string) {
	_ = s.productRepo.IncrementViewCount(productID)
}

// updateProductRating 更新商品评分
func (s *ProductService) updateProductRating(productID string) {
	// 获取商品的所有评价
	reviews, err := s.productRepo.FindAllReviewsByProductID(productID)
	if err != nil {
		return
	}
	
	if len(reviews) == 0 {
		return
	}
	
	// 计算平均评分
	var totalRating int
	var publishedCount int
	
	for _, review := range reviews {
		if review.Status == domain.ReviewStatusPublished {
			totalRating += review.Rating
			publishedCount++
		}
	}
	
	if publishedCount > 0 {
		averageRating := float64(totalRating) / float64(publishedCount)
		_ = s.productRepo.UpdateRating(productID, averageRating, publishedCount)
	}
}