package repository

import (
	"github.com/huabtc/polygame/backend/internal/model"
	"gorm.io/gorm"
)

type PositionRepository struct {
	db *gorm.DB
}

func NewPositionRepository(db *gorm.DB) *PositionRepository {
	return &PositionRepository{db: db}
}

// FindByUserAndOutcome 根据用户和结果查找持仓
func (r *PositionRepository) FindByUserAndOutcome(userID, marketID, outcomeID uint) (*model.Position, error) {
	var position model.Position
	err := r.db.Where("user_id = ? AND market_id = ? AND outcome_id = ?", userID, marketID, outcomeID).
		First(&position).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &position, nil
}

// Create 创建持仓
func (r *PositionRepository) Create(position *model.Position) error {
	return r.db.Create(position).Error
}

// Update 更新持仓
func (r *PositionRepository) Update(position *model.Position) error {
	return r.db.Save(position).Error
}

// FindByUserID 根据用户 ID 查找所有持仓
func (r *PositionRepository) FindByUserID(userID uint) ([]model.Position, error) {
	var positions []model.Position
	err := r.db.Preload("Market").Preload("Outcome").
		Where("user_id = ? AND shares > 0", userID).
		Find(&positions).Error
	return positions, err
}

// FindByMarketID 根据市场 ID 查找所有持仓
func (r *PositionRepository) FindByMarketID(marketID uint) ([]model.Position, error) {
	var positions []model.Position
	err := r.db.Where("market_id = ? AND shares > 0", marketID).Find(&positions).Error
	return positions, err
}

// Delete 删除持仓（份额为 0 时）
func (r *PositionRepository) Delete(positionID uint) error {
	return r.db.Delete(&model.Position{}, positionID).Error
}
