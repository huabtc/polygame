package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huabtc/polygame/backend/internal/service"
)

type TradingHandler struct {
	tradingService *service.TradingService
}

func NewTradingHandler(tradingService *service.TradingService) *TradingHandler {
	return &TradingHandler{tradingService: tradingService}
}

// PlaceOrder 下单
func (h *TradingHandler) PlaceOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		MarketID   uint    `json:"market_id" binding:"required"`
		OutcomeID  uint    `json:"outcome_id" binding:"required"`
		OrderType  string  `json:"order_type" binding:"required,oneof=buy sell"`
		Shares     float64 `json:"shares" binding:"required,gt=0"`
		Price      float64 `json:"price" binding:"required,gt=0,lte=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.tradingService.PlaceOrder(
		userID,
		req.MarketID,
		req.OutcomeID,
		req.OrderType,
		req.Shares,
		req.Price,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order": order})
}

// GetUserOrders 获取用户订单列表
func (h *TradingHandler) GetUserOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	var pageInt, pageSizeInt int
	fmt.Sscanf(page, "%d", &pageInt)
	fmt.Sscanf(pageSize, "%d", &pageSizeInt)

	if pageInt < 1 {
		pageInt = 1
	}
	if pageSizeInt < 1 || pageSizeInt > 100 {
		pageSizeInt = 20
	}

	orders, total, err := h.tradingService.GetUserOrders(userID, pageInt, pageSizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  total,
		"page":   pageInt,
	})
}

// GetUserPositions 获取用户持仓
func (h *TradingHandler) GetUserPositions(c *gin.Context) {
	userID := c.GetUint("user_id")

	positions, err := h.tradingService.GetUserPositions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"positions": positions})
}

// CancelOrder 取消订单
func (h *TradingHandler) CancelOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.tradingService.CancelOrder(uri.ID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}
