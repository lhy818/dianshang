package repository

import (
	"errors"

	"gorm.io/gorm"

	"ecommerce-backend/internal/domain"
)

// OrderRepository 订单仓库接口
type OrderRepository interface {
	// 订单相关
	Create(order *domain.Order) error
	CreateWithTx(tx *gorm.DB, order *domain.Order) error
	FindByID(id string) (*domain.Order, error)
	FindByOrderNo(orderNo string) (*domain.Order, error)
	FindOrders(conditions map[string]interface{}, page, pageSize int) ([]*domain.Order, int64, error)
	Update(order *domain.Order) error
	UpdateWithTx(tx *gorm.DB, order *domain.Order) error
	Delete(id string) error
	Begin() *gorm.DB

	// 订单商品项相关
	CreateItem(item *domain.OrderItem) error
	CreateItemWithTx(tx *gorm.DB, item *domain.OrderItem) error
	FindItemByID(id string) (*domain.OrderItem, error)
	FindItemsByOrderID(orderID string) ([]*domain.OrderItem, error)
	UpdateItem(item *domain.OrderItem) error
	DeleteItem(id string) error

	// 支付记录相关
	CreatePayment(payment *domain.Payment) error
	FindPaymentByID(id string) (*domain.Payment, error)
	FindPaymentByPaymentNo(paymentNo string) (*domain.Payment, error)
	FindPaymentsByOrderID(orderID string) ([]*domain.Payment, error)
	UpdatePayment(payment *domain.Payment) error
	DeletePayment(id string) error
}

// GormOrderRepository GORM实现的订单仓库
type GormOrderRepository struct {
	db *gorm.DB
}

// NewGormOrderRepository 创建GORM订单仓库
func NewGormOrderRepository(db *gorm.DB) OrderRepository {
	return &GormOrderRepository{db: db}
}

// Begin 开始事务
func (r *GormOrderRepository) Begin() *gorm.DB {
	return r.db.Begin()
}

// Create 创建订单
func (r *GormOrderRepository) Create(order *domain.Order) error {
	return r.db.Create(order).Error
}

// CreateWithTx 在事务中创建订单
func (r *GormOrderRepository) CreateWithTx(tx *gorm.DB, order *domain.Order) error {
	return tx.Create(order).Error
}

// FindByID 根据ID查找订单
func (r *GormOrderRepository) FindByID(id string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items").Preload("Payments").
		Where("id = ?", id).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &order, err
}

// FindByOrderNo 根据订单号查找订单
func (r *GormOrderRepository) FindByOrderNo(orderNo string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Preload("Items").Preload("Payments").
		Where("order_no = ?", orderNo).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &order, err
}

// FindOrders 查找订单列表
func (r *GormOrderRepository) FindOrders(conditions map[string]interface{}, page, pageSize int) ([]*domain.Order, int64, error) {
	var orders []*domain.Order
	var total int64

	// 构建查询
	dbQuery := r.db.Model(&domain.Order{})

	// 应用条件
	for key, value := range conditions {
		dbQuery = dbQuery.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 计算总数
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := dbQuery.Preload("Items").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

// Update 更新订单
func (r *GormOrderRepository) Update(order *domain.Order) error {
	return r.db.Save(order).Error
}

// UpdateWithTx 在事务中更新订单
func (r *GormOrderRepository) UpdateWithTx(tx *gorm.DB, order *domain.Order) error {
	return tx.Save(order).Error
}

// Delete 删除订单
func (r *GormOrderRepository) Delete(id string) error {
	return r.db.Delete(&domain.Order{}, "id = ?", id).Error
}

// CreateItem 创建订单商品项
func (r *GormOrderRepository) CreateItem(item *domain.OrderItem) error {
	return r.db.Create(item).Error
}

// CreateItemWithTx 在事务中创建订单商品项
func (r *GormOrderRepository) CreateItemWithTx(tx *gorm.DB, item *domain.OrderItem) error {
	return tx.Create(item).Error
}

// FindItemByID 根据ID查找订单商品项
func (r *GormOrderRepository) FindItemByID(id string) (*domain.OrderItem, error) {
	var item domain.OrderItem
	err := r.db.Preload("Product").Preload("SKU").
		Where("id = ?", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

// FindItemsByOrderID 根据订单ID查找所有商品项
func (r *GormOrderRepository) FindItemsByOrderID(orderID string) ([]*domain.OrderItem, error) {
	var items []*domain.OrderItem
	err := r.db.Preload("Product").Preload("SKU").
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

// UpdateItem 更新订单商品项
func (r *GormOrderRepository) UpdateItem(item *domain.OrderItem) error {
	return r.db.Save(item).Error
}

// DeleteItem 删除订单商品项
func (r *GormOrderRepository) DeleteItem(id string) error {
	return r.db.Delete(&domain.OrderItem{}, "id = ?", id).Error
}

// CreatePayment 创建支付记录
func (r *GormOrderRepository) CreatePayment(payment *domain.Payment) error {
	return r.db.Create(payment).Error
}

// FindPaymentByID 根据ID查找支付记录
func (r *GormOrderRepository) FindPaymentByID(id string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Where("id = ?", id).First(&payment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &payment, err
}

// FindPaymentByPaymentNo 根据支付单号查找支付记录
func (r *GormOrderRepository) FindPaymentByPaymentNo(paymentNo string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Where("payment_no = ?", paymentNo).First(&payment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &payment, err
}

// FindPaymentsByOrderID 根据订单ID查找支付记录
func (r *GormOrderRepository) FindPaymentsByOrderID(orderID string) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	err := r.db.Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&payments).Error
	return payments, err
}

// UpdatePayment 更新支付记录
func (r *GormOrderRepository) UpdatePayment(payment *domain.Payment) error {
	return r.db.Save(payment).Error
}

// DeletePayment 删除支付记录
func (r *GormOrderRepository) DeletePayment(id string) error {
	return r.db.Delete(&domain.Payment{}, "id = ?", id).Error
}

// fmt包需要导入
import "fmt"