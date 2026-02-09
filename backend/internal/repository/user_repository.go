package repository

import (
	"errors"

	"gorm.io/gorm"

	"ecommerce-backend/internal/domain"
)

// UserRepository 用户仓库接口
type UserRepository interface {
	// 用户相关
	Create(user *domain.User) error
	FindByID(id string) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
	List(page, pageSize int, filters map[string]interface{}) ([]*domain.User, int64, error)

	// 用户资料相关
	CreateProfile(profile *domain.UserProfile) error
	FindProfileByID(id string) (*domain.UserProfile, error)
	FindProfileByUserID(userID string) (*domain.UserProfile, error)
	UpdateProfile(profile *domain.UserProfile) error
	DeleteProfile(id string) error

	// 用户地址相关
	CreateAddress(address *domain.UserAddress) error
	FindAddressByID(id string) (*domain.UserAddress, error)
	FindAddressesByUserID(userID string) ([]*domain.UserAddress, error)
	UpdateAddress(address *domain.UserAddress) error
	DeleteAddress(id string) error
	ClearDefaultAddress(userID string) error
}

// GormUserRepository GORM实现的用户仓库
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository 创建GORM用户仓库
func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

// Create 创建用户
func (r *GormUserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// FindByID 根据ID查找用户
func (r *GormUserRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// FindByUsername 根据用户名查找用户
func (r *GormUserRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// FindByEmail 根据邮箱查找用户
func (r *GormUserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Update 更新用户
func (r *GormUserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *GormUserRepository) Delete(id string) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

// List 用户列表
func (r *GormUserRepository) List(page, pageSize int, filters map[string]interface{}) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	query := r.db.Model(&domain.User{}).Where("deleted_at IS NULL")

	// 应用过滤器
	for key, value := range filters {
		switch key {
		case "username":
			query = query.Where("username LIKE ?", "%"+value.(string)+"%")
		case "email":
			query = query.Where("email LIKE ?", "%"+value.(string)+"%")
		case "status":
			query = query.Where("status = ?", value)
		case "start_date":
			query = query.Where("created_at >= ?", value)
		case "end_date":
			query = query.Where("created_at <= ?", value)
		}
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error

	return users, total, err
}

// CreateProfile 创建用户资料
func (r *GormUserRepository) CreateProfile(profile *domain.UserProfile) error {
	return r.db.Create(profile).Error
}

// FindProfileByID 根据ID查找用户资料
func (r *GormUserRepository) FindProfileByID(id string) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	err := r.db.Where("id = ?", id).First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, err
}

// FindProfileByUserID 根据用户ID查找用户资料
func (r *GormUserRepository) FindProfileByUserID(userID string) (*domain.UserProfile, error) {
	var profile domain.UserProfile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, err
}

// UpdateProfile 更新用户资料
func (r *GormUserRepository) UpdateProfile(profile *domain.UserProfile) error {
	return r.db.Save(profile).Error
}

// DeleteProfile 删除用户资料
func (r *GormUserRepository) DeleteProfile(id string) error {
	return r.db.Delete(&domain.UserProfile{}, "id = ?", id).Error
}

// CreateAddress 创建用户地址
func (r *GormUserRepository) CreateAddress(address *domain.UserAddress) error {
	return r.db.Create(address).Error
}

// FindAddressByID 根据ID查找用户地址
func (r *GormUserRepository) FindAddressByID(id string) (*domain.UserAddress, error) {
	var address domain.UserAddress
	err := r.db.Where("id = ?", id).First(&address).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &address, err
}

// FindAddressesByUserID 根据用户ID查找所有地址
func (r *GormUserRepository) FindAddressesByUserID(userID string) ([]*domain.UserAddress, error) {
	var addresses []*domain.UserAddress
	err := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error
	return addresses, err
}

// UpdateAddress 更新用户地址
func (r *GormUserRepository) UpdateAddress(address *domain.UserAddress) error {
	return r.db.Save(address).Error
}

// DeleteAddress 删除用户地址
func (r *GormUserRepository) DeleteAddress(id string) error {
	return r.db.Delete(&domain.UserAddress{}, "id = ?", id).Error
}

// ClearDefaultAddress 清除用户的默认地址
func (r *GormUserRepository) ClearDefaultAddress(userID string) error {
	return r.db.Model(&domain.UserAddress{}).
		Where("user_id = ? AND is_default = true", userID).
		Update("is_default", false).Error
}