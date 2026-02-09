package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 健康检查响应
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Services  struct {
		Database  string `json:"database"`
		Redis     string `json:"redis"`
		RabbitMQ  string `json:"rabbitmq"`
	} `json:"services"`
}

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("警告: .env文件未找到，使用环境变量")
	}

	// 初始化数据库连接
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// 创建Gin路由器
	r := gin.Default()

	// 中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		response := HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
		}

		// 检查数据库连接
		sqlDB, err := db.DB()
		if err == nil && sqlDB.Ping() == nil {
			response.Services.Database = "connected"
		} else {
			response.Services.Database = "disconnected"
		}

		// TODO: 添加Redis和RabbitMQ健康检查
		response.Services.Redis = "not_implemented"
		response.Services.RabbitMQ = "not_implemented"

		c.JSON(200, response)
	})

	// API版本前缀
	api := r.Group("/api/v1")

	// 认证路由
	auth := api.Group("/auth")
	{
		auth.POST("/register", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "注册功能待实现",
				"status":  "development",
			})
		})
		auth.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "登录功能待实现",
				"status":  "development",
			})
		})
		auth.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "用户信息功能待实现",
				"status":  "development",
			})
		})
	}

	// 商品路由
	products := api.Group("/products")
	{
		products.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "商品列表功能待实现",
				"status":  "development",
			})
		})
		products.GET("/:id", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "商品详情功能待实现",
				"status":  "development",
			})
		})
		products.POST("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "创建商品功能待实现",
				"status":  "development",
			})
		})
	}

	// 分类路由
	categories := api.Group("/categories")
	{
		categories.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "分类列表功能待实现",
				"status":  "development",
			})
		})
	}

	// 购物车路由
	cart := api.Group("/cart")
	{
		cart.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "购物车功能待实现",
				"status":  "development",
			})
		})
		cart.POST("/add", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "添加商品到购物车功能待实现",
				"status":  "development",
			})
		})
	}

	// 订单路由
	orders := api.Group("/orders")
	{
		orders.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "订单列表功能待实现",
				"status":  "development",
			})
		})
		orders.POST("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "创建订单功能待实现",
				"status":  "development",
			})
		})
	}

	// 根路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":   "电商系统API",
			"version":   "1.0.0",
			"status":    "development",
			"timestamp": time.Now(),
			"endpoints": []string{
				"GET  /health",
				"POST /api/v1/auth/register",
				"POST /api/v1/auth/login",
				"GET  /api/v1/auth/profile",
				"GET  /api/v1/products",
				"GET  /api/v1/products/:id",
				"POST /api/v1/products",
				"GET  /api/v1/categories",
				"GET  /api/v1/cart",
				"POST /api/v1/cart/add",
				"GET  /api/v1/orders",
				"POST /api/v1/orders",
			},
		})
	})

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 服务器启动在端口 %s", port)
	log.Printf("📡 健康检查: http://localhost:%s/health", port)
	log.Printf("📚 API文档: http://localhost:%s/", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// 初始化数据库连接
func initDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=ecommerce port=5432 sslmode=disable"
	}

	// 配置GORM日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	log.Println("✅ 数据库连接成功")
	return db, nil
}

// CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}