package domain

import (
	"time"
)

// Category 商品分类实体
type Category struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ParentID    string    `gorm:"type:uuid;index" json:"parent_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description string    `json:"description"`
	ImageURL    string    `gorm:"type:varchar(500)" json:"image_url"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// 关联关系
	Parent   *Category `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// Product 商品实体
type Product struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CategoryID       string    `gorm:"type:uuid;index" json:"category_id"`
	Name             string    `gorm:"type:varchar(200);not null" json:"name"`
	Slug             string    `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug"`
	Description      string    `json:"description"`
	ShortDescription string    `gorm:"type:varchar(500)" json:"short_description"`
	SKU              string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"sku"`
	Price            float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice    float64   `gorm:"type:decimal(10,2)" json:"original_price"`
	CostPrice        float64   `gorm:"type:decimal(10,2)" json:"cost_price"`
	StockQuantity    int       `gorm:"default:0" json:"stock_quantity"`
	SoldQuantity     int       `gorm:"default:0" json:"sold_quantity"`
	Weight           float64   `gorm:"type:decimal(8,3);default:0" json:"weight"`
	Volume           float64   `gorm:"type:decimal(8,3);default:0" json:"volume"`
	MainImageURL     string    `gorm:"type:varchar(500)" json:"main_image_url"`
	ImageURLs        JSONB     `gorm:"type:jsonb;default:'[]'" json:"image_urls"`
	Attributes       JSONB     `gorm:"type:jsonb;default:'{}'" json:"attributes"`
	IsActive         bool      `gorm:"default:true;index" json:"is_active"`
	IsRecommended    bool      `gorm:"default:false" json:"is_recommended"`
	IsHot            bool      `gorm:"default:false" json:"is_hot"`
	IsNew            bool      `gorm:"default:false" json:"is_new"`
	ViewCount        int       `gorm:"default:0" json:"view_count"`
	Rating           float64   `gorm:"type:decimal(3,2);default:0;index" json:"rating"`
	ReviewCount      int       `gorm:"default:0" json:"review_count"`
	CreatedAt        time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        time.Time `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关联关系
	Category *Category    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	SKUs     []ProductSKU `gorm:"foreignKey:ProductID" json:"skus,omitempty"`
	Reviews  []ProductReview `gorm:"foreignKey:ProductID" json:"reviews,omitempty"`
}

// ProductSKU 商品SKU实体
type ProductSKU struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProductID    string    `gorm:"type:uuid;not null;index" json:"product_id"`
	SKUCode      string    `gorm:"type:varchar(50);not null" json:"sku_code"`
	Attributes   JSONB     `gorm:"type:jsonb;not null;default:'{}'" json:"attributes"`
	Price        float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice float64  `gorm:"type:decimal(10,2)" json:"original_price"`
	StockQuantity int      `gorm:"default:0" json:"stock_quantity"`
	SoldQuantity  int      `gorm:"default:0" json:"sold_quantity"`
	ImageURL     string    `gorm:"type:varchar(500)" json:"image_url"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// 关联关系
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// ProductReview 商品评价实体
type ProductReview struct {
	ID                 string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderItemID        string    `gorm:"type:uuid;uniqueIndex;not null" json:"order_item_id"`
	UserID             string    `gorm:"type:uuid;not null;index" json:"user_id"`
	ProductID          string    `gorm:"type:uuid;not null;index" json:"product_id"`
	Rating             int       `gorm:"not null" json:"rating"`
	Title              string    `gorm:"type:varchar(200)" json:"title"`
	Content            string    `json:"content"`
	ImageURLs          JSONB     `gorm:"type:jsonb;default:'[]'" json:"image_urls"`
	IsAnonymous        bool      `gorm:"default:false" json:"is_anonymous"`
	IsVerifiedPurchase bool      `gorm:"default:true" json:"is_verified_purchase"`
	LikeCount          int       `gorm:"default:0" json:"like_count"`
	ReplyCount         int       `gorm:"default:0" json:"reply_count"`
	Status             string    `gorm:"type:varchar(20);default:'published'" json:"status"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
	// 关联关系
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// JSONB 类型包装
type JSONB map[string]interface{}

// Scan 实现Scanner接口
func (j *JSONB) Scan(value interface{}) error {
	// 实现JSONB扫描逻辑
	return nil
}

// Value 实现Valuer接口
func (j JSONB) Value() (interface{}, error) {
	// 实现JSONB值转换逻辑
	return j, nil
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

func (Product) TableName() string {
	return "products"
}

func (ProductSKU) TableName() string {
	return "product_skus"
}

func (ProductReview) TableName() string {
	return "product_reviews"
}

// ReviewStatus 评价状态常量
const (
	ReviewStatusPending   = "pending"
	ReviewStatusPublished = "published"
	ReviewStatusHidden    = "hidden"
)