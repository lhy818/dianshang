package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/service"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService *service.UserService
	jwtSecret   string
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(userService *service.UserService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User         *domain.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	existingUser, err := h.userService.GetUserByUsername(req.Username)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	existingUser, err = h.userService.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "邮箱已存在",
		})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "密码加密失败",
			"error":   err.Error(),
		})
		return
	}

	// 创建用户
	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword),
		Status:       domain.UserStatusActive,
	}

	if err := h.userService.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "用户创建失败",
			"error":   err.Error(),
		})
		return
	}

	// 生成JWT令牌
	accessToken, refreshToken, expiresAt, err := h.generateTokens(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "令牌生成失败",
			"error":   err.Error(),
		})
		return
	}

	// 隐藏敏感信息
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, gin.H{
		"code": 201,
		"message": "注册成功",
		"data": AuthResponse{
			User:         user,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    expiresAt,
		},
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 查找用户（支持用户名或邮箱登录）
	var user *domain.User
	var err error
	
	// 先尝试按用户名查找
	user, err = h.userService.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		// 如果用户名找不到，尝试按邮箱查找
		user, err = h.userService.GetUserByEmail(req.Username)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "用户名或密码错误",
			})
			return
		}
	}

	// 检查用户状态
	if user.Status != domain.UserStatusActive {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "用户账号已被禁用",
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户名或密码错误",
		})
		return
	}

	// 更新最后登录时间
	user.LastLoginAt = time.Now()
	if err := h.userService.UpdateUser(user); err != nil {
		// 登录成功但更新登录时间失败，不影响主要功能
	}

	// 生成JWT令牌
	accessToken, refreshToken, expiresAt, err := h.generateTokens(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "令牌生成失败",
			"error":   err.Error(),
		})
		return
	}

	// 隐藏敏感信息
	user.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "登录成功",
		"data": AuthResponse{
			User:         user,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    expiresAt,
		},
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 在实际应用中，这里可以将令牌加入黑名单
	// 由于JWT是无状态的，客户端需要自己删除令牌
	
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登出成功",
	})
}

// RefreshToken 刷新访问令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "刷新令牌不能为空",
		})
		return
	}

	// 验证刷新令牌
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "刷新令牌无效",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "令牌格式错误",
		})
		return
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "令牌无效",
		})
		return
	}

	// 查找用户
	user, err := h.userService.GetUserByID(userID)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户不存在",
		})
		return
	}

	// 生成新的访问令牌
	accessToken, _, expiresAt, err := h.generateTokens(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "令牌生成失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "令牌刷新成功",
		"data": gin.H{
			"access_token": accessToken,
			"expires_at":   expiresAt,
		},
	})
}

// GetCurrentUser 获取当前用户信息
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}

	user, err := h.userService.GetUserByID(userID.(string))
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	// 隐藏敏感信息
	user.PasswordHash = ""

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    user,
	})
}

// generateTokens 生成JWT令牌
func (h *AuthHandler) generateTokens(userID, username string) (accessToken, refreshToken string, expiresAt time.Time, err error) {
	// 访问令牌有效期：2小时
	accessExpiresAt := time.Now().Add(2 * time.Hour)
	accessClaims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      accessExpiresAt.Unix(),
		"iat":      time.Now().Unix(),
		"type":     "access",
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessJWT.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	// 刷新令牌有效期：7天
	refreshExpiresAt := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      refreshExpiresAt.Unix(),
		"iat":      time.Now().Unix(),
		"type":     "refresh",
	}

	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshJWT.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken, refreshToken, accessExpiresAt, nil
}

// AuthMiddleware JWT认证中间件
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌不能为空",
			})
			c.Abort()
			return
		}

		// 移除Bearer前缀
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		// 解析令牌
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证令牌无效",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "令牌格式错误",
			})
			c.Abort()
			return
		}

		// 检查令牌类型
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "令牌类型错误",
			})
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		userID, ok := claims["user_id"].(string)
		if ok {
			c.Set("user_id", userID)
		}

		username, ok := claims["username"].(string)
		if ok {
			c.Set("username", username)
		}

		c.Next()
	}
}