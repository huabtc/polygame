package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huabtc/polygame/backend/internal/model"
	"github.com/huabtc/polygame/backend/internal/service"
)

type MarketHandler struct {
	marketService *service.MarketService
}

func NewMarketHandler(marketService *service.MarketService) *MarketHandler {
	return &MarketHandler{marketService: marketService}
}

// CreateMarket 创建市场（管理员）
func (h *MarketHandler) CreateMarket(c *gin.Context) {
	var req struct {
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description"`
		Category    string   `json:"category" binding:"required"`
		ImageURL    string   `json:"image_url"`
		StartTime   *string  `json:"start_time"`
		EndTime     *string  `json:"end_time"`
		Outcomes    []string `json:"outcomes" binding:"required,min=2"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy := c.GetUint("user_id")

	market := &model.Market{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		Status:      "active",
		CreatedBy:   createdBy,
	}

	if req.StartTime != nil {
		t, _ := time.Parse(time.RFC3339, *req.StartTime)
		market.StartTime = &t
	}
	if req.EndTime != nil {
		t, _ := time.Parse(time.RFC3339, *req.EndTime)
		market.EndTime = &t
	}

	if err := h.marketService.CreateMarket(market, req.Outcomes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"market": market})
}

// GetMarket 获取市场详情
func (h *MarketHandler) GetMarket(c *gin.Context) {
	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	market, err := h.marketService.GetMarket(uri.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"market": market})
}

// ListMarkets 获取市场列表
func (h *MarketHandler) ListMarkets(c *gin.Context) {
	category := c.Query("category")
	status := c.Query("status")
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

	markets, total, err := h.marketService.ListMarkets(category, status, pageInt, pageSizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"markets": markets,
		"total":   total,
		"page":    pageInt,
	})
}

// UpdateMarket 更新市场（管理员）
func (h *MarketHandler) UpdateMarket(c *gin.Context) {
	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	market, err := h.marketService.GetMarket(uri.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
		ImageURL    *string `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != nil {
		market.Title = *req.Title
	}
	if req.Description != nil {
		market.Description = *req.Description
	}
	if req.Status != nil {
		market.Status = *req.Status
	}
	if req.ImageURL != nil {
		market.ImageURL = *req.ImageURL
	}

	if err := h.marketService.UpdateMarket(market); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"market": market})
}

// ResolveMarket 结算市场（管理员）
func (h *MarketHandler) ResolveMarket(c *gin.Context) {
	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		WinningOutcomeID uint `json:"winning_outcome_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resolvedBy := c.GetUint("user_id")

	if err := h.marketService.ResolveMarket(uri.ID, req.WinningOutcomeID, resolvedBy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Market resolved successfully"})
}

// GetTrendingMarkets 获取热门市场
func (h *MarketHandler) GetTrendingMarkets(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	var limitInt int
	fmt.Sscanf(limit, "%d", &limitInt)

	if limitInt < 1 || limitInt > 50 {
		limitInt = 10
	}

	markets, err := h.marketService.GetTrendingMarkets(limitInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"markets": markets})
}

// SearchMarkets 搜索市场
func (h *MarketHandler) SearchMarkets(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search keyword required"})
		return
	}

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

	markets, total, err := h.marketService.SearchMarkets(keyword, pageInt, pageSizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"markets": markets,
		"total":   total,
		"page":    pageInt,
	})
}
