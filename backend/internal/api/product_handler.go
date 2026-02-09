package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/service"
)

// ProductHandler 商品处理器
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler 创建商品处理器
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetCategories 获取分类列表
func (h *ProductHandler) GetCategories(c *gin.Context) {
	categories, err := h.productService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取分类列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    categories,
	})
}

// GetProducts 获取商品列表
func (h *ProductHandler) GetProducts(c *gin.Context) {
	// 解析查询参数
	query := c.Query("q")
	categoryID := c.Query("category_id")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")
	isActive := c.Query("is_active")
	isRecommended := c.Query("is_recommended")
	isHot := c.Query("is_hot")
	isNew := c.Query("is_new")

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 构建查询条件
	filter := &service.ProductFilter{
		Query:          query,
		CategoryID:     categoryID,
		MinPrice:       parseFloat(minPrice),
		MaxPrice:       parseFloat(maxPrice),
		IsActive:       parseBool(isActive),
		IsRecommended:  parseBool(isRecommended),
		IsHot:          parseBool(isHot),
		IsNew:          parseBool(isNew),
		Page:           page,
		PageSize:       pageSize,
		SortBy:         sortBy,
		SortOrder:      strings.ToUpper(sortOrder),
	}

	// 获取商品列表
	products, total, err := h.productService.GetProducts(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取商品列表失败",
			"error":   err.Error(),
		})
		return
	}

	// 计算分页信息
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"products":    products,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": totalPages,
			},
		},
	})
}

// GetProductByID 获取商品详情
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "商品ID不能为空",
		})
		return
	}

	product, err := h.productService.GetProductByID(productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "商品不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取商品详情失败",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    product,
	})
}

// GetProductSKUs 获取商品SKU列表
func (h *ProductHandler) GetProductSKUs(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "商品ID不能为空",
		})
		return
	}

	skus, err := h.productService.GetProductSKUs(productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "商品不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取商品SKU列表失败",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    skus,
	})
}

// parseFloat 解析浮点数
func parseFloat(s string) *float64 {
	if s == "" {
		return nil
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &val
}

// parseBool 解析布尔值
func parseBool(s string) *bool {
	if s == "" {
		return nil
	}
	val, err := strconv.ParseBool(s)
	if err != nil {
		return nil
	}
	return &val
}