package domain

import (
	"time"
)

// ShoppingCart 购物车实体
type ShoppingCart struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    string    `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	SessionID string    `gorm:"type:varchar(100)" json:"session_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// 关联关系
	Items []CartItem `gorm:"foreignKey:CartID" json:"items,omitempty"`
	User  *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// CartItem 购物车商品项实体
type CartItem struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CartID    string    `gorm:"type:uuid;not null;index" json:"cart_id"`
	ProductID string    `gorm:"type:uuid;not null;index" json:"product_id"`
	SKUID     string    `gorm:"type:uuid;index" json:"sku_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	UnitPrice float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// 关联关系
	Cart    *ShoppingCart `gorm:"foreignKey:CartID" json:"cart,omitempty"`
	Product *Product      `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SKU     *ProductSKU   `gorm:"foreignKey:SKUID" json:"sku,omitempty"`
}

// TableName 指定表名
func (ShoppingCart) TableName() string {
	return "shopping_carts"
}

func (CartItem) TableName() string {
	return "cart_items"
}