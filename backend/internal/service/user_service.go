package service

import (
	"errors"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/repository"
)

// UserService 用户服务
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *domain.User) error {
	// 验证用户数据
	if err := s.validateUser(user); err != nil {
		return err
	}

	return s.userRepo.Create(user)
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id string) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return s.userRepo.FindByUsername(username)
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(email)
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(user *domain.User) error {
	if user.ID == "" {
		return errors.New("用户ID不能为空")
	}

	// 验证用户数据
	if err := s.validateUser(user); err != nil {
		return err
	}

	return s.userRepo.Update(user)
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}

// ListUsers 用户列表
func (s *UserService) ListUsers(page, pageSize int, filters map[string]interface{}) ([]*domain.User, int64, error) {
	return s.userRepo.List(page, pageSize, filters)
}

// CreateUserProfile 创建用户资料
func (s *UserService) CreateUserProfile(profile *domain.UserProfile) error {
	if profile.UserID == "" {
		return errors.New("用户ID不能为空")
	}

	// 检查用户是否存在
	user, err := s.userRepo.FindByID(profile.UserID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	// 检查是否已存在资料
	existingProfile, err := s.userRepo.FindProfileByUserID(profile.UserID)
	if err == nil && existingProfile != nil {
		return errors.New("用户资料已存在")
	}

	return s.userRepo.CreateProfile(profile)
}

// GetUserProfile 获取用户资料
func (s *UserService) GetUserProfile(userID string) (*domain.UserProfile, error) {
	return s.userRepo.FindProfileByUserID(userID)
}

// UpdateUserProfile 更新用户资料
func (s *UserService) UpdateUserProfile(profile *domain.UserProfile) error {
	if profile.ID == "" {
		return errors.New("资料ID不能为空")
	}

	return s.userRepo.UpdateProfile(profile)
}

// CreateUserAddress 创建用户地址
func (s *UserService) CreateUserAddress(address *domain.UserAddress) error {
	if address.UserID == "" {
		return errors.New("用户ID不能为空")
	}

	// 检查用户是否存在
	user, err := s.userRepo.FindByID(address.UserID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	// 如果是默认地址，取消其他默认地址
	if address.IsDefault {
		if err := s.userRepo.ClearDefaultAddress(address.UserID); err != nil {
			return err
		}
	}

	return s.userRepo.CreateAddress(address)
}

// GetUserAddress 获取用户地址
func (s *UserService) GetUserAddress(id string) (*domain.UserAddress, error) {
	return s.userRepo.FindAddressByID(id)
}

// GetUserAddresses 获取用户所有地址
func (s *UserService) GetUserAddresses(userID string) ([]*domain.UserAddress, error) {
	return s.userRepo.FindAddressesByUserID(userID)
}

// UpdateUserAddress 更新用户地址
func (s *UserService) UpdateUserAddress(address *domain.UserAddress) error {
	if address.ID == "" {
		return errors.New("地址ID不能为空")
	}

	// 如果是默认地址，取消其他默认地址
	if address.IsDefault {
		if err := s.userRepo.ClearDefaultAddress(address.UserID); err != nil {
			return err
		}
	}

	return s.userRepo.UpdateAddress(address)
}

// DeleteUserAddress 删除用户地址
func (s *UserService) DeleteUserAddress(id string) error {
	return s.userRepo.DeleteAddress(id)
}

// SetDefaultAddress 设置默认地址
func (s *UserService) SetDefaultAddress(userID, addressID string) error {
	// 检查地址是否存在且属于该用户
	address, err := s.userRepo.FindAddressByID(addressID)
	if err != nil || address == nil {
		return errors.New("地址不存在")
	}

	if address.UserID != userID {
		return errors.New("地址不属于该用户")
	}

	// 取消其他默认地址
	if err := s.userRepo.ClearDefaultAddress(userID); err != nil {
		return err
	}

	// 设置该地址为默认
	address.IsDefault = true
	return s.userRepo.UpdateAddress(address)
}

// validateUser 验证用户数据
func (s *UserService) validateUser(user *domain.User) error {
	if user.Username == "" {
		return errors.New("用户名不能为空")
	}

	if len(user.Username) < 3 || len(user.Username) > 50 {
		return errors.New("用户名长度必须在3-50个字符之间")
	}

	if user.Email == "" {
		return errors.New("邮箱不能为空")
	}

	// 简单的邮箱格式验证
	if !isValidEmail(user.Email) {
		return errors.New("邮箱格式不正确")
	}

	if user.PasswordHash == "" {
		return errors.New("密码不能为空")
	}

	// 验证状态
	if user.Status != "" {
		validStatuses := map[string]bool{
			domain.UserStatusActive:   true,
			domain.UserStatusInactive: true,
			domain.UserStatusBanned:   true,
		}
		if !validStatuses[user.Status] {
			return errors.New("用户状态无效")
		}
	}

	return nil
}

// isValidEmail 简单的邮箱格式验证
func isValidEmail(email string) bool {
	// 这里应该使用更严格的邮箱验证
	// 为了简单起见，只做基本检查
	if len(email) < 3 || len(email) > 100 {
		return false
	}
	
	// 检查是否包含@和.
	hasAt := false
	hasDot := false
	for _, ch := range email {
		if ch == '@' {
			hasAt = true
		}
		if ch == '.' && hasAt {
			hasDot = true
		}
	}
	
	return hasAt && hasDot
}