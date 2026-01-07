package service

import (
	"errors"

	"github.com/huabtc/polygame/backend/internal/model"
	"github.com/huabtc/polygame/backend/internal/repository"
	"gorm.io/gorm"
)

type MarketService struct {
	marketRepo   *repository.MarketRepository
	positionRepo *repository.PositionRepository
	userRepo     *repository.UserRepository
	txRepo       *repository.TransactionRepository
	db           *gorm.DB
}

func NewMarketService(
	marketRepo *repository.MarketRepository,
	positionRepo *repository.PositionRepository,
	userRepo *repository.UserRepository,
	txRepo *repository.TransactionRepository,
	db *gorm.DB,
) *MarketService {
	return &MarketService{
		marketRepo:   marketRepo,
		positionRepo: positionRepo,
		userRepo:     userRepo,
		txRepo:       txRepo,
		db:           db,
	}
}

// CreateMarket 创建市场（管理员）
func (s *MarketService) CreateMarket(market *model.Market, outcomes []string) error {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建市场
	if err := tx.Create(market).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建结果选项
	for _, outcomeName := range outcomes {
		outcome := model.Outcome{
			MarketID:     market.ID,
			OutcomeName:  outcomeName,
			CurrentPrice: 0.5, // 初始价格 0.5
		}
		if err := tx.Create(&outcome).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetMarket 获取市场详情
func (s *MarketService) GetMarket(marketID uint) (*model.Market, error) {
	return s.marketRepo.FindByID(marketID)
}

// ListMarkets 获取市场列表
func (s *MarketService) ListMarkets(category, status string, page, pageSize int) ([]model.Market, int64, error) {
	return s.marketRepo.List(category, status, page, pageSize)
}

// UpdateMarket 更新市场（管理员）
func (s *MarketService) UpdateMarket(market *model.Market) error {
	return s.marketRepo.Update(market)
}

// ResolveMarket 结算市场（管理员）
func (s *MarketService) ResolveMarket(marketID, winningOutcomeID, resolvedBy uint) error {
	// 验证市场状态
	market, err := s.marketRepo.FindByID(marketID)
	if err != nil {
		return err
	}
	if market.Status == "resolved" {
		return errors.New("market already resolved")
	}

	// 验证获胜结果是否存在
	validOutcome := false
	for _, outcome := range market.Outcomes {
		if outcome.ID == winningOutcomeID {
			validOutcome = true
			break
		}
	}
	if !validOutcome {
		return errors.New("invalid winning outcome")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新市场状态
	if err := tx.Model(&model.Market{}).
		Where("id = ?", marketID).
		Updates(map[string]interface{}{
			"status":          "resolved",
			"winning_outcome": winningOutcomeID,
			"resolved_by":     resolvedBy,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 获取所有持仓
	positions, err := s.positionRepo.FindByMarketID(marketID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 结算持仓
	for _, position := range positions {
		if position.OutcomeID == winningOutcomeID {
			// 获胜方：每份额获得 1.0 虚拟积分
			payout := position.Shares * 1.0

			if err := tx.Model(&model.User{}).
				Where("id = ?", position.UserID).
				UpdateColumn("virtual_balance", gorm.Expr("virtual_balance + ?", payout)).
				Error; err != nil {
				tx.Rollback()
				return err
			}

			// 记录交易
			user, _ := s.userRepo.FindByID(position.UserID)
			txRecord := &model.Transaction{
				UserID:       position.UserID,
				Type:         "settlement_win",
				Amount:       payout,
				BalanceAfter: user.VirtualBalance + payout,
				MarketID:     &marketID,
				Description:  "Market settlement - win",
			}
			tx.Create(txRecord)
		} else {
			// 失败方：记录损失
			user, _ := s.userRepo.FindByID(position.UserID)
			txRecord := &model.Transaction{
				UserID:       position.UserID,
				Type:         "settlement_loss",
				Amount:       0,
				BalanceAfter: user.VirtualBalance,
				MarketID:     &marketID,
				Description:  "Market settlement - loss",
			}
			tx.Create(txRecord)
		}
	}

	return tx.Commit().Error
}

// GetTrendingMarkets 获取热门市场
func (s *MarketService) GetTrendingMarkets(limit int) ([]model.Market, error) {
	return s.marketRepo.GetTrending(limit)
}

// SearchMarkets 搜索市场
func (s *MarketService) SearchMarkets(keyword string, page, pageSize int) ([]model.Market, int64, error) {
	return s.marketRepo.Search(keyword, page, pageSize)
}
