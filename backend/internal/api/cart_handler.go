package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/service"
)

// CartHandler 购物车处理器
type CartHandler struct {
	cartService *service.CartService
}

// NewCartHandler 创建购物车处理器
func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// GetCart 获取购物车
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	cart, err := h.cartService.GetCartByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果购物车不存在，创建一个空的购物车
			cart = &domain.ShoppingCart{
				UserID: userID,
				Items:  []domain.CartItem{},
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取购物车失败",
				"error":   err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    cart,
	})
}

// AddCartItemRequest 添加购物车商品请求
type AddCartItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	SKUID     string `json:"sku_id"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// AddCartItem 添加商品到购物车
func (h *CartHandler) AddCartItem(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	var req AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 添加商品到购物车
	cartItem, err := h.cartService.AddCartItem(userID, req.ProductID, req.SKUID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "添加商品到购物车失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "商品已添加到购物车",
		"data":    cartItem,
	})
}

// UpdateCartItemRequest 更新购物车商品请求
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// UpdateCartItem 更新购物车商品
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	cartItemID := c.Param("id")
	if cartItemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "购物车商品ID不能为空",
		})
		return
	}

	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 更新购物车商品
	cartItem, err := h.cartService.UpdateCartItem(userID, cartItemID, req.Quantity)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "购物车商品不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新购物车商品失败",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "购物车商品更新成功",
		"data":    cartItem,
	})
}

// DeleteCartItem 删除购物车商品
func (h *CartHandler) DeleteCartItem(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	cartItemID := c.Param("id")
	if cartItemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "购物车商品ID不能为空",
		})
		return
	}

	// 删除购物车商品
	err := h.cartService.DeleteCartItem(userID, cartItemID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "购物车商品不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除购物车商品失败",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "购物车商品删除成功",
	})
}