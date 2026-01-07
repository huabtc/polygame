package service

import (
	"errors"
	"time"

	"github.com/huabtc/polygame/backend/internal/model"
	"github.com/huabtc/polygame/backend/internal/repository"
	"gorm.io/gorm"
)

type TradingService struct {
	orderRepo    *repository.OrderRepository
	positionRepo *repository.PositionRepository
	userRepo     *repository.UserRepository
	marketRepo   *repository.MarketRepository
	txRepo       *repository.TransactionRepository
	db           *gorm.DB
}

func NewTradingService(
	orderRepo *repository.OrderRepository,
	positionRepo *repository.PositionRepository,
	userRepo *repository.UserRepository,
	marketRepo *repository.MarketRepository,
	txRepo *repository.TransactionRepository,
	db *gorm.DB,
) *TradingService {
	return &TradingService{
		orderRepo:    orderRepo,
		positionRepo: positionRepo,
		userRepo:     userRepo,
		marketRepo:   marketRepo,
		txRepo:       txRepo,
		db:           db,
	}
}

// PlaceOrder 下单
func (s *TradingService) PlaceOrder(userID, marketID, outcomeID uint, orderType string, shares, price float64) (*model.Order, error) {
	// 验证市场状态
	market, err := s.marketRepo.FindByID(marketID)
	if err != nil {
		return nil, err
	}
	if market.Status != "active" {
		return nil, errors.New("market is not active")
	}

	// 计算总成本
	totalCost := shares * price

	// 验证用户余额（买入时）
	if orderType == "buy" {
		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			return nil, err
		}
		if user.VirtualBalance < totalCost {
			return nil, errors.New("insufficient balance")
		}
	}

	// 验证持仓（卖出时）
	if orderType == "sell" {
		position, err := s.positionRepo.FindByUserAndOutcome(userID, marketID, outcomeID)
		if err != nil {
			return nil, err
		}
		if position == nil || position.Shares < shares {
			return nil, errors.New("insufficient shares")
		}
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建订单
	order := &model.Order{
		UserID:     userID,
		MarketID:   marketID,
		OutcomeID:  outcomeID,
		OrderType:  orderType,
		Shares:     shares,
		Price:      price,
		TotalCost:  totalCost,
		Status:     "filled", // 简化处理，直接成交
		FilledAt:   &time.Time{},
	}
	*order.FilledAt = time.Now()

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 更新用户余额
	if orderType == "buy" {
		if err := tx.Model(&model.User{}).
			Where("id = ?", userID).
			UpdateColumn("virtual_balance", gorm.Expr("virtual_balance - ?", totalCost)).
			Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// 记录交易
		user, _ := s.userRepo.FindByID(userID)
		txRecord := &model.Transaction{
			UserID:       userID,
			Type:         "trade_buy",
			Amount:       -totalCost,
			BalanceAfter: user.VirtualBalance - totalCost,
			OrderID:      &order.ID,
			MarketID:     &marketID,
			Description:  "Buy shares",
		}
		tx.Create(txRecord)

	} else {
		if err := tx.Model(&model.User{}).
			Where("id = ?", userID).
			UpdateColumn("virtual_balance", gorm.Expr("virtual_balance + ?", totalCost)).
			Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// 记录交易
		user, _ := s.userRepo.FindByID(userID)
		txRecord := &model.Transaction{
			UserID:       userID,
			Type:         "trade_sell",
			Amount:       totalCost,
			BalanceAfter: user.VirtualBalance + totalCost,
			OrderID:      &order.ID,
			MarketID:     &marketID,
			Description:  "Sell shares",
		}
		tx.Create(txRecord)
	}

	// 更新持仓
	if err := s.updatePosition(tx, userID, marketID, outcomeID, orderType, shares, price); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return order, nil
}

// updatePosition 更新持仓
func (s *TradingService) updatePosition(tx *gorm.DB, userID, marketID, outcomeID uint, orderType string, shares, price float64) error {
	var position model.Position
	err := tx.Where("user_id = ? AND market_id = ? AND outcome_id = ?", userID, marketID, outcomeID).
		First(&position).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新持仓
		if orderType == "buy" {
			position = model.Position{
				UserID:    userID,
				MarketID:  marketID,
				OutcomeID: outcomeID,
				Shares:    shares,
				AvgPrice:  price,
			}
			return tx.Create(&position).Error
		}
		return errors.New("no position to sell")
	}

	if err != nil {
		return err
	}

	// 更新持仓
	if orderType == "buy" {
		totalCost := position.Shares*position.AvgPrice + shares*price
		position.Shares += shares
		position.AvgPrice = totalCost / position.Shares
	} else {
		position.Shares -= shares
		if position.Shares < 0 {
			return errors.New("insufficient shares")
		}
	}

	return tx.Save(&position).Error
}

// GetUserOrders 获取用户订单列表
func (s *TradingService) GetUserOrders(userID uint, page, pageSize int) ([]model.Order, int64, error) {
	return s.orderRepo.FindByUserID(userID, page, pageSize)
}

// GetUserPositions 获取用户持仓
func (s *TradingService) GetUserPositions(userID uint) ([]model.Position, error) {
	return s.positionRepo.FindByUserID(userID)
}

// CancelOrder 取消订单
func (s *TradingService) CancelOrder(orderID, userID uint) error {
	return s.orderRepo.Cancel(orderID, userID)
}
