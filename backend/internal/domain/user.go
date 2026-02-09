package domain

import (
	"time"
)

// User 用户实体
type User struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Username          string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email             string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Phone             string    `gorm:"type:varchar(20)" json:"phone"`
	PasswordHash      string    `gorm:"type:varchar(255);not null" json:"-"`
	AvatarURL         string    `gorm:"type:varchar(500)" json:"avatar_url"`
	Status            string    `gorm:"type:varchar(20);default:'active';index" json:"status"`
	IsEmailVerified   bool      `gorm:"default:false" json:"is_email_verified"`
	IsPhoneVerified   bool      `gorm:"default:false" json:"is_phone_verified"`
	LastLoginAt       time.Time `json:"last_login_at"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         time.Time `gorm:"index" json:"deleted_at,omitempty"`
	
	// 关联关系
	Profile   *UserProfile   `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	Addresses []UserAddress  `gorm:"foreignKey:UserID" json:"addresses,omitempty"`
}

// UserProfile 用户资料实体
type UserProfile struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID       string    `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	RealName     string    `gorm:"type:varchar(100)" json:"real_name"`
	Gender       string    `gorm:"type:varchar(10)" json:"gender"`
	BirthDate    time.Time `json:"birth_date"`
	IDCardNumber string    `gorm:"type:varchar(20)" json:"id_card_number"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// UserAddress 用户地址实体
type UserAddress struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID        string    `gorm:"type:uuid;not null;index" json:"user_id"`
	RecipientName string    `gorm:"type:varchar(100);not null" json:"recipient_name"`
	Phone         string    `gorm:"type:varchar(20);not null" json:"phone"`
	Province      string    `gorm:"type:varchar(50);not null" json:"province"`
	City          string    `gorm:"type:varchar(50);not null" json:"city"`
	District      string    `gorm:"type:varchar(50);not null" json:"district"`
	StreetAddress string    `gorm:"type:varchar(200);not null" json:"street_address"`
	PostalCode    string    `gorm:"type:varchar(10)" json:"postal_code"`
	IsDefault     bool      `gorm:"default:false" json:"is_default"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// UserStatus 用户状态常量
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusBanned   = "banned"
)

// Gender 性别常量
const (
	GenderMale   = "male"
	GenderFemale = "female"
	GenderOther  = "other"
)

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

func (UserAddress) TableName() string {
	return "user_addresses"
}