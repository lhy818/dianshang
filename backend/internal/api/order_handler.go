package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/service"
)

// OrderHandler 订单处理器
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// GetOrders 获取订单列表
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	// 解析查询参数
	status := c.Query("status")
	paymentStatus := c.Query("payment_status")

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
	filter := &service.OrderFilter{
		UserID:        userID,
		Status:        status,
		PaymentStatus: paymentStatus,
		Page:          page,
		PageSize:      pageSize,
	}

	// 获取订单列表
	orders, total, err := h.orderService.GetOrders(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取订单列表失败",
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
			"orders": orders,
			"pagination": gin.H{
				"page":        page,
				"page_size":   pageSize,
				"total":       total,
				"total_pages": totalPages,
			},
		},
	})
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	ShippingAddress domain.JSONB `json:"shipping_address" binding:"required"`
	ShippingMethod  string       `json:"shipping_method"`
	BuyerRemark     string       `json:"buyer_remark"`
}

// CreateOrder 创建订单
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 创建订单
	order, err := h.orderService.CreateOrder(userID, req.ShippingAddress, req.ShippingMethod, req.BuyerRemark)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建订单失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "订单创建成功",
		"data":    order,
	})
}

// GetOrderByID 获取订单详情
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "订单ID不能为空",
		})
		return
	}

	// 获取订单详情
	order, err := h.orderService.GetOrderByID(userID, orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "订单不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取订单详情失败",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    order,
	})
}

// CancelOrder 取消订单
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "订单ID不能为空",
		})
		return
	}

	// 取消订单
	order, err := h.orderService.CancelOrder(userID, orderID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "订单不存在",
			})
		} else if err.Error() == "订单状态不允许取消" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "订单状态不允许取消",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "取消订单失败",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "订单取消成功",
		"data":    order,
	})
}