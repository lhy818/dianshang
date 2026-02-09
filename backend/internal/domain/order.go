package domain

import (
	"time"
)

// Order 订单实体
type Order struct {
	ID              string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderNo         string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_no"`
	UserID          string    `gorm:"type:uuid;not null;index" json:"user_id"`
	TotalAmount     float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	DiscountAmount  float64   `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	ShippingFee     float64   `gorm:"type:decimal(10,2);default:0" json:"shipping_fee"`
	FinalAmount     float64   `gorm:"type:decimal(10,2);not null" json:"final_amount"`
	Status          string    `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	PaymentStatus   string    `gorm:"type:varchar(20);default:'unpaid';index" json:"payment_status"`
	PaymentMethod   string    `gorm:"type:varchar(50)" json:"payment_method"`
	PaymentTime     time.Time `json:"payment_time"`
	ShippingAddress JSONB     `gorm:"type:jsonb;not null" json:"shipping_address"`
	ShippingMethod  string    `gorm:"type:varchar(50)" json:"shipping_method"`
	ShippingNo      string    `gorm:"type:varchar(100)" json:"shipping_no"`
	ShippingTime    time.Time `json:"shipping_time"`
	DeliveredTime   time.Time `json:"delivered_time"`
	BuyerRemark     string    `json:"buyer_remark"`
	SellerRemark    string    `json:"seller_remark"`
	CancelledReason string    `json:"cancelled_reason"`
	CancelledTime   time.Time `json:"cancelled_time"`
	CreatedAt       time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// 关联关系
	User      *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items     []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Payments  []Payment   `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
}

// OrderItem 订单商品项实体
type OrderItem struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID       string    `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID     string    `gorm:"type:uuid;not null" json:"product_id"`
	SKUID         string    `gorm:"type:uuid" json:"sku_id"`
	ProductName   string    `gorm:"type:varchar(200);not null" json:"product_name"`
	SKUAttributes JSONB     `gorm:"type:jsonb;default:'{}'" json:"sku_attributes"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	UnitPrice     float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	TotalPrice    float64   `gorm:"type:decimal(10,2);not null" json:"total_price"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	
	// 关联关系
	Order   *Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SKU     *ProductSKU `gorm:"foreignKey:SKUID" json:"sku,omitempty"`
}

// Payment 支付记录实体
type Payment struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID       string    `gorm:"type:uuid;not null;index" json:"order_id"`
	PaymentNo     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"payment_no"`
	PaymentMethod string    `gorm:"type:varchar(50);not null" json:"payment_method"`
	PaymentChannel string   `gorm:"type:varchar(50)" json:"payment_channel"`
	Amount        float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        string    `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	TransactionID string    `gorm:"type:varchar(100)" json:"transaction_id"`
	PayerInfo     JSONB     `gorm:"type:jsonb" json:"payer_info"`
	RawResponse   JSONB     `gorm:"type:jsonb" json:"raw_response"`
	PaidAt        time.Time `json:"paid_at"`
	CreatedAt     time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// 关联关系
	Order *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// OrderStatus 订单状态常量
const (
	OrderStatusPending   = "pending"   // 待支付
	OrderStatusPaid      = "paid"      // 已支付
	OrderStatusShipped   = "shipped"   // 已发货
	OrderStatusDelivered = "delivered" // 已送达
	OrderStatusCompleted = "completed" // 已完成
	OrderStatusCancelled = "cancelled" // 已取消
	OrderStatusRefunded  = "refunded"  // 已退款
)

// PaymentStatus 支付状态常量
const (
	PaymentStatusUnpaid  = "unpaid"   // 未支付
	PaymentStatusPaid    = "paid"     // 已支付
	PaymentStatusRefunded = "refunded" // 已退款
	PaymentStatusFailed  = "failed"   // 支付失败
)

// PaymentStatus 支付记录状态常量
const (
	PaymentRecordStatusPending = "pending" // 待支付
	PaymentRecordStatusSuccess = "success" // 支付成功
	PaymentRecordStatusFailed  = "failed"  // 支付失败
	PaymentRecordStatusRefunded = "refunded" // 已退款
)

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

func (OrderItem) TableName() string {
	return "order_items"
}

func (Payment) TableName() string {
	return "payments"
}