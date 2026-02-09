package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ecommerce-backend/internal/api"
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/repository"
	"ecommerce-backend/internal/service"
)

func main() {
	// 加载配置
	loadConfig()

	// 初始化数据库
	db := initDatabase()
	
	// 自动迁移数据库表
	if err := autoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化仓库
	userRepo := repository.NewGormUserRepository(db)
	productRepo := repository.NewGormProductRepository(db)
	cartRepo := repository.NewGormCartRepository(db)
	orderRepo := repository.NewGormOrderRepository(db)

	// 初始化服务
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo)

	// 初始化处理器
	jwtSecret := viper.GetString("jwt.secret")
	authHandler := api.NewAuthHandler(userService, jwtSecret)
	productHandler := api.NewProductHandler(productService)
	cartHandler := api.NewCartHandler(cartService)
	orderHandler := api.NewOrderHandler(orderService)

	// 创建Gin引擎
	r := gin.Default()

	// 配置中间件
	configureMiddleware(r)

	// 注册路由
	registerRoutes(r, authHandler, productHandler, cartHandler, orderHandler)

	// 启动服务器
	port := viper.GetString("app.port")
	log.Printf("服务器启动在 http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// loadConfig 加载配置
func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// 设置默认值
	viper.SetDefault("app.port", "8080")
	viper.SetDefault("app.env", "development")
	viper.SetDefault("database.postgres.host", "localhost")
	viper.SetDefault("database.postgres.port", 5432)
	viper.SetDefault("database.postgres.dbname", "ecommerce")
	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("配置文件未找到，使用默认配置")
		} else {
			log.Fatalf("读取配置文件失败: %v", err)
		}
	}

	// 从环境变量覆盖配置
	viper.AutomaticEnv()
}

// initDatabase 初始化数据库连接
func initDatabase() *gorm.DB {
	dsn := buildDSN()
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(viper.GetInt("database.postgres.max_idle_conns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("database.postgres.max_open_conns"))
	sqlDB.SetConnMaxLifetime(viper.GetDuration("database.postgres.conn_max_lifetime"))

	log.Println("数据库连接成功")
	return db
}

// buildDSN 构建数据库连接字符串
func buildDSN() string {
	host := viper.GetString("database.postgres.host")
	port := viper.GetInt("database.postgres.port")
	user := viper.GetString("database.postgres.user")
	password := viper.GetString("database.postgres.password")
	dbname := viper.GetString("database.postgres.dbname")
	sslmode := viper.GetString("database.postgres.sslmode")

	return "host=" + host + " port=" + string(rune(port)) + " user=" + user + 
		" password=" + password + " dbname=" + dbname + " sslmode=" + sslmode
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	// 只在开发环境自动迁移
	if viper.GetString("app.env") != "development" {
		log.Println("非开发环境，跳过自动迁移")
		return nil
	}

	log.Println("开始自动迁移数据库表...")

	// 迁移所有表
	err := db.AutoMigrate(
		&domain.User{},
		&domain.UserProfile{},
		&domain.UserAddress{},
		&domain.Category{},
		&domain.Product{},
		&domain.ProductSKU{},
		&domain.ProductReview{},
		&domain.ShoppingCart{},
		&domain.CartItem{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Payment{},
	)

	if err != nil {
		return err
	}

	log.Println("数据库表迁移完成")
	return nil
}

// configureMiddleware 配置中间件
func configureMiddleware(r *gin.Engine) {
	// 恢复中间件
	r.Use(gin.Recovery())

	// 日志中间件
	if viper.GetString("app.env") == "development" {
		r.Use(gin.Logger())
	}

	// CORS中间件
	r.Use(corsMiddleware())

	// 请求ID中间件
	r.Use(requestIDMiddleware())
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// requestIDMiddleware 请求ID中间件
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// 生成简单的请求ID
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 这里应该使用更复杂的ID生成算法
	// 为了简单起见，使用时间戳
	return "req_" + string(rune(os.Getpid())) + "_" + string(rune(os.Getppid()))
}

// registerRoutes 注册路由
func registerRoutes(r *gin.Engine, authHandler *api.AuthHandler, productHandler *api.ProductHandler, cartHandler *api.CartHandler, orderHandler *api.OrderHandler) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	// API版本前缀
	apiV1 := r.Group("/api/v1")
	{
		// 认证相关路由
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 需要认证的路由
		authenticated := apiV1.Group("")
		authenticated.Use(authHandler.AuthMiddleware())
		{
			// 用户相关
			users := authenticated.Group("/users")
			{
				users.GET("/me", authHandler.GetCurrentUser)
			}

			// 认证相关
			auth := authenticated.Group("/auth")
			{
				auth.GET("/profile", authHandler.GetCurrentUser)
			}

			// 购物车相关
			cart := authenticated.Group("/cart")
			{
				cart.GET("", cartHandler.GetCart)
				cart.POST("/items", cartHandler.AddCartItem)
				cart.PUT("/items/:id", cartHandler.UpdateCartItem)
				cart.DELETE("/items/:id", cartHandler.DeleteCartItem)
			}

			// 订单相关
			orders := authenticated.Group("/orders")
			{
				orders.GET("", orderHandler.GetOrders)
				orders.POST("", orderHandler.CreateOrder)
				orders.GET("/:id", orderHandler.GetOrderByID)
				orders.PUT("/:id/cancel", orderHandler.CancelOrder)
			}
		}

		// 公开路由（不需要认证）
		// 商品相关
		products := apiV1.Group("/products")
		{
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProductByID)
			products.GET("/:id/skus", productHandler.GetProductSKUs)
		}

		// 分类相关
		categories := apiV1.Group("/categories")
		{
			categories.GET("", productHandler.GetCategories)
		}
	}

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code": 404,
			"message": "接口不存在",
			"path": c.Request.URL.Path,
		})
	})
}

// 导入time包
import "time"