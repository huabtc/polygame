package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/huabtc/polygame/backend/config"
	"github.com/huabtc/polygame/backend/internal/api"
	"github.com/huabtc/polygame/backend/internal/middleware"
	"github.com/huabtc/polygame/backend/internal/repository"
	"github.com/huabtc/polygame/backend/internal/service"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	if err := repository.InitDatabase(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	db := repository.GetDB()

	// 初始化仓储层
	userRepo := repository.NewUserRepository(db)
	marketRepo := repository.NewMarketRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	positionRepo := repository.NewPositionRepository(db)
	txRepo := repository.NewTransactionRepository(db)

	// 初始化服务层
	userService := service.NewUserService(userRepo, txRepo, cfg)
	marketService := service.NewMarketService(marketRepo, positionRepo, userRepo, txRepo, db)
	tradingService := service.NewTradingService(orderRepo, positionRepo, userRepo, marketRepo, txRepo, db)

	// 初始化处理器
	userHandler := api.NewUserHandler(userService)
	marketHandler := api.NewMarketHandler(marketService)
	tradingHandler := api.NewTradingHandler(tradingService)

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由组
	apiV1 := r.Group("/api/v1")

	// 公开路由
	auth := apiV1.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
	}

	// 需要认证的路由
	authenticated := apiV1.Group("")
	authenticated.Use(middleware.AuthMiddleware(cfg))
	{
		// 用户相关
		user := authenticated.Group("/user")
		{
			user.GET("/profile", userHandler.GetProfile)
			user.PUT("/profile", userHandler.UpdateProfile)
			user.GET("/balance", userHandler.GetBalance)
		}

		// 市场相关
		markets := authenticated.Group("/markets")
		{
			markets.GET("", marketHandler.ListMarkets)
			markets.GET("/trending", marketHandler.GetTrendingMarkets)
			markets.GET("/search", marketHandler.SearchMarkets)
			markets.GET("/:id", marketHandler.GetMarket)
		}

		// 交易相关
		trading := authenticated.Group("/trading")
		{
			trading.POST("/orders", tradingHandler.PlaceOrder)
			trading.GET("/orders", tradingHandler.GetUserOrders)
			trading.DELETE("/orders/:id", tradingHandler.CancelOrder)
			trading.GET("/positions", tradingHandler.GetUserPositions)
		}

		// 管理员路由
		admin := authenticated.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/users", userHandler.ListUsers)
			admin.POST("/markets", marketHandler.CreateMarket)
			admin.PUT("/markets/:id", marketHandler.UpdateMarket)
			admin.POST("/markets/:id/resolve", marketHandler.ResolveMarket)
		}
	}

	// 启动服务器
	addr := ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
