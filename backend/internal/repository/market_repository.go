package repository

import (
	"errors"

	"github.com/huabtc/polygame/backend/internal/model"
	"gorm.io/gorm"
)

type MarketRepository struct {
	db *gorm.DB
}

func NewMarketRepository(db *gorm.DB) *MarketRepository {
	return &MarketRepository{db: db}
}

// Create 创建市场
func (r *MarketRepository) Create(market *model.Market) error {
	return r.db.Create(market).Error
}

// FindByID 根据 ID 查找市场
func (r *MarketRepository) FindByID(id uint) (*model.Market, error) {
	var market model.Market
	err := r.db.Preload("Outcomes").First(&market, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("market not found")
		}
		return nil, err
	}
	return &market, nil
}

// List 获取市场列表
func (r *MarketRepository) List(category string, status string, page, pageSize int) ([]model.Market, int64, error) {
	var markets []model.Market
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Market{})

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Outcomes").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&markets).Error

	return markets, total, err
}

// Update 更新市场
func (r *MarketRepository) Update(market *model.Market) error {
	return r.db.Save(market).Error
}

// UpdateStatus 更新市场状态
func (r *MarketRepository) UpdateStatus(marketID uint, status string) error {
	return r.db.Model(&model.Market{}).
		Where("id = ?", marketID).
		Update("status", status).Error
}

// Resolve 结算市场
func (r *MarketRepository) Resolve(marketID uint, winningOutcomeID uint, resolvedBy uint) error {
	return r.db.Model(&model.Market{}).
		Where("id = ?", marketID).
		Updates(map[string]interface{}{
			"status":          "resolved",
			"winning_outcome": winningOutcomeID,
			"resolved_by":     resolvedBy,
		}).Error
}

// GetTrending 获取热门市场
func (r *MarketRepository) GetTrending(limit int) ([]model.Market, error) {
	var markets []model.Market
	err := r.db.Preload("Outcomes").
		Where("status = ?", "active").
		Order("total_volume DESC").
		Limit(limit).
		Find(&markets).Error
	return markets, err
}

// Search 搜索市场
func (r *MarketRepository) Search(keyword string, page, pageSize int) ([]model.Market, int64, error) {
	var markets []model.Market
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Market{}).
		Where("title ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Outcomes").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&markets).Error

	return markets, total, err
}
