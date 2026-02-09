package service

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/repository"
)

// OrderService 订单服务
type OrderService struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

// NewOrderService 创建订单服务
func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// OrderFilter 订单查询过滤器
type OrderFilter struct {
	UserID        string
	Status        string
	PaymentStatus string
	Page          int
	PageSize      int
}

// GetOrders 获取订单列表
func (s *OrderService) GetOrders(filter *OrderFilter) ([]*domain.Order, int64, error) {
	if filter.UserID == "" {
		return nil, 0, errors.New("用户ID不能为空")
	}
	
	// 构建查询条件
	conditions := make(map[string]interface{})
	conditions["user_id"] = filter.UserID
	
	if filter.Status != "" {
		conditions["status"] = filter.Status
	}
	
	if filter.PaymentStatus != "" {
		conditions["payment_status"] = filter.PaymentStatus
	}
	
	return s.orderRepo.FindOrders(conditions, filter.Page, filter.PageSize)
}

// GetOrderByID 根据ID获取订单
func (s *OrderService) GetOrderByID(userID, orderID string) (*domain.Order, error) {
	if userID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	
	if orderID == "" {
		return nil, errors.New("订单ID不能为空")
	}
	
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	
	// 验证订单属于该用户
	if order.UserID != userID {
		return nil, errors.New("订单不属于该用户")
	}
	
	return order, nil
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(userID string, shippingAddress domain.JSONB, shippingMethod, buyerRemark string) (*domain.Order, error) {
	if userID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	
	// 验证收货地址
	if len(shippingAddress) == 0 {
		return nil, errors.New("收货地址不能为空")
	}
	
	// 获取购物车
	cart, err := s.cartRepo.FindByUserID(userID)
	if err != nil || cart == nil {
		return nil, errors.New("购物车为空")
	}
	
	// 获取购物车商品
	cartItems, err := s.cartRepo.FindItemsByCartID(cart.ID)
	if err != nil || len(cartItems) == 0 {
		return nil, errors.New("购物车为空")
	}
	
	// 开始事务
	tx := s.orderRepo.Begin()
	
	// 创建订单
	order := &domain.Order{
		OrderNo:         generateOrderNo(),
		UserID:          userID,
		ShippingAddress: shippingAddress,
		ShippingMethod:  shippingMethod,
		BuyerRemark:     buyerRemark,
		Status:          domain.OrderStatusPending,
		PaymentStatus:   domain.PaymentStatusUnpaid,
	}
	
	// 计算订单金额
	var totalAmount float64
	var orderItems []*domain.OrderItem
	
	for _, cartItem := range cartItems {
		// 获取商品信息
		product, err := s.productRepo.FindByID(cartItem.ProductID)
		if err != nil || product == nil || !product.IsActive {
			tx.Rollback()
			return nil, fmt.Errorf("商品不存在或已下架: %s", cartItem.ProductID)
		}
		
		// 检查库存
		var stockQuantity int
		var sku *domain.ProductSKU
		
		if cartItem.SKUID != "" {
			sku, err = s.productRepo.FindSKUByID(cartItem.SKUID)
			if err != nil || sku == nil || !sku.IsActive {
				tx.Rollback()
				return nil, fmt.Errorf("商品SKU不存在或已下架: %s", cartItem.SKUID)
			}
			stockQuantity = sku.StockQuantity
		} else {
			stockQuantity = product.StockQuantity
		}
		
		// 检查库存是否充足
		if stockQuantity < cartItem.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("商品库存不足: %s", product.Name)
		}
		
		// 扣减库存
		if cartItem.SKUID != "" {
			if err := s.productRepo.UpdateSKUStockWithTx(tx, cartItem.SKUID, -cartItem.Quantity); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("扣减SKU库存失败: %v", err)
			}
		} else {
			if err := s.productRepo.UpdateStockWithTx(tx, cartItem.ProductID, -cartItem.Quantity); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("扣减商品库存失败: %v", err)
			}
		}
		
		// 创建订单商品项
		orderItem := &domain.OrderItem{
			OrderID:       order.ID,
			ProductID:     cartItem.ProductID,
			SKUID:         cartItem.SKUID,
			ProductName:   product.Name,
			SKUAttributes: make(domain.JSONB),
			Quantity:      cartItem.Quantity,
			UnitPrice:     cartItem.UnitPrice,
			TotalPrice:    cartItem.UnitPrice * float64(cartItem.Quantity),
		}
		
		if sku != nil {
			orderItem.SKUAttributes = sku.Attributes
		}
		
		orderItems = append(orderItems, orderItem)
		totalAmount += orderItem.TotalPrice
	}
	
	// 设置订单金额
	order.TotalAmount = totalAmount
	order.FinalAmount = totalAmount // 暂时不考虑优惠和运费
	
	// 保存订单
	if err := s.orderRepo.CreateWithTx(tx, order); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建订单失败: %v", err)
	}
	
	// 保存订单商品项
	for _, orderItem := range orderItems {
		orderItem.OrderID = order.ID
		if err := s.orderRepo.CreateItemWithTx(tx, orderItem); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("创建订单商品项失败: %v", err)
		}
	}
	
	// 清空购物车
	if err := s.cartRepo.ClearItemsWithTx(tx, cart.ID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("清空购物车失败: %v", err)
	}
	
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}
	
	return order, nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(userID, orderID string) (*domain.Order, error) {
	if userID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	
	if orderID == "" {
		return nil, errors.New("订单ID不能为空")
	}
	
	// 获取订单
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	
	// 验证订单属于该用户
	if order.UserID != userID {
		return nil, errors.New("订单不属于该用户")
	}
	
	// 检查订单状态是否可以取消
	if order.Status != domain.OrderStatusPending && order.Status != domain.OrderStatusPaid {
		return nil, errors.New("订单状态不允许取消")
	}
	
	// 开始事务
	tx := s.orderRepo.Begin()
	
	// 恢复库存
	orderItems, err := s.orderRepo.FindItemsByOrderID(orderID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	
	for _, orderItem := range orderItems {
		if orderItem.SKUID != "" {
			if err := s.productRepo.UpdateSKUStockWithTx(tx, orderItem.SKUID, orderItem.Quantity); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("恢复SKU库存失败: %v", err)
			}
		} else {
			if err := s.productRepo.UpdateStockWithTx(tx, orderItem.ProductID, orderItem.Quantity); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("恢复商品库存失败: %v", err)
			}
		}
	}
	
	// 更新订单状态
	order.Status = domain.OrderStatusCancelled
	order.CancelledTime = time.Now()
	
	if err := s.orderRepo.UpdateWithTx(tx, order); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新订单状态失败: %v", err)
	}
	
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}
	
	return order, nil
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(orderID, status string) error {
	if orderID == "" {
		return errors.New("订单ID不能为空")
	}
	
	if status == "" {
		return errors.New("订单状态不能为空")
	}
	
	// 验证状态
	validStatuses := map[string]bool{
		domain.OrderStatusPending:   true,
		domain.OrderStatusPaid:      true,
		domain.OrderStatusShipped:   true,
		domain.OrderStatusDelivered: true,
		domain.OrderStatusCompleted: true,
		domain.OrderStatusCancelled: true,
		domain.OrderStatusRefunded:  true,
	}
	
	if !validStatuses[status] {
		return errors.New("订单状态无效")
	}
	
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	
	order.Status = status
	
	// 根据状态设置相应时间
	now := time.Now()
	switch status {
	case domain.OrderStatusPaid:
		order.PaymentTime = now
	case domain.OrderStatusShipped:
		order.ShippingTime = now
	case domain.OrderStatusDelivered:
		order.DeliveredTime = now
	case domain.OrderStatusCancelled:
		order.CancelledTime = now
	}
	
	return s.orderRepo.Update(order)
}

// UpdatePaymentStatus 更新支付状态
func (s *OrderService) UpdatePaymentStatus(orderID, paymentStatus string) error {
	if orderID == "" {
		return errors.New("订单ID不能为空")
	}
	
	if paymentStatus == "" {
		return errors.New("支付状态不能为空")
	}
	
	// 验证状态
	validStatuses := map[string]bool{
		domain.PaymentStatusUnpaid:  true,
		domain.PaymentStatusPaid:    true,
		domain.PaymentStatusRefunded: true,
		domain.PaymentStatusFailed:  true,
	}
	
	if !validStatuses[paymentStatus] {
		return errors.New("支付状态无效")
	}
	
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}
	
	order.PaymentStatus = paymentStatus
	
	// 如果支付成功，更新订单状态为已支付
	if paymentStatus == domain.PaymentStatusPaid {
		order.Status = domain.OrderStatusPaid
		order.PaymentTime = time.Now()
	}
	
	return s.orderRepo.Update(order)
}

// generateOrderNo 生成订单号
func generateOrderNo() string {
	// 格式: YYYYMMDDHHMMSS + 6位随机数
	now := time.Now()
	timestamp := now.Format("20060102150405")
	random := fmt.Sprintf("%06d", now.Nanosecond()%1000000)
	return timestamp + random
}