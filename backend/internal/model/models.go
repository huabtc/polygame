package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	Username       string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email          string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	PasswordHash   string         `gorm:"size:255;not null" json:"-"`
	VirtualBalance float64        `gorm:"type:decimal(20,2);default:10000" json:"virtual_balance"` // 初始虚拟积分 10000
	Avatar         string         `gorm:"size:255" json:"avatar"`
	IsAdmin        bool           `gorm:"default:false" json:"is_admin"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// Market 市场模型
type Market struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	Title          string         `gorm:"size:255;not null" json:"title"`
	Description    string         `gorm:"type:text" json:"description"`
	Category       string         `gorm:"size:50;not null;index" json:"category"` // sports, esports, entertainment, tech
	ImageURL       string         `gorm:"size:500" json:"image_url"`
	StartTime      *time.Time     `json:"start_time"`
	EndTime        *time.Time     `json:"end_time"`
	ResolutionTime *time.Time     `json:"resolution_time"`
	Status         string         `gorm:"size:20;not null;default:'pending';index" json:"status"` // pending, active, closed, resolved, cancelled
	TotalVolume    float64        `gorm:"type:decimal(20,2);default:0" json:"total_volume"`
	CreatedBy      uint           `gorm:"not null" json:"created_by"`
	ResolvedBy     *uint          `json:"resolved_by"`
	WinningOutcome *uint          `json:"winning_outcome"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Outcomes       []Outcome      `gorm:"foreignKey:MarketID" json:"outcomes,omitempty"`
}

// Outcome 市场结果选项模型
type Outcome struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	MarketID     uint           `gorm:"not null;index" json:"market_id"`
	OutcomeName  string         `gorm:"size:100;not null" json:"outcome_name"`
	CurrentPrice float64        `gorm:"type:decimal(10,4);default:0.5" json:"current_price"` // 0-1 之间
	TotalShares  float64        `gorm:"type:decimal(20,2);default:0" json:"total_shares"`
	TotalVolume  float64        `gorm:"type:decimal(20,2);default:0" json:"total_volume"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Order 订单模型
type Order struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	MarketID   uint           `gorm:"not null;index" json:"market_id"`
	OutcomeID  uint           `gorm:"not null;index" json:"outcome_id"`
	OrderType  string         `gorm:"size:10;not null" json:"order_type"` // buy, sell
	Shares     float64        `gorm:"type:decimal(20,4);not null" json:"shares"`
	Price      float64        `gorm:"type:decimal(10,4);not null" json:"price"`
	TotalCost  float64        `gorm:"type:decimal(20,2);not null" json:"total_cost"`
	Status     string         `gorm:"size:20;not null;default:'pending';index" json:"status"` // pending, filled, partially_filled, cancelled
	FilledAt   *time.Time     `json:"filled_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Market     Market         `gorm:"foreignKey:MarketID" json:"market,omitempty"`
	Outcome    Outcome        `gorm:"foreignKey:OutcomeID" json:"outcome,omitempty"`
}

// Position 持仓模型
type Position struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index:idx_user_market_outcome,priority:1" json:"user_id"`
	MarketID  uint           `gorm:"not null;index:idx_user_market_outcome,priority:2" json:"market_id"`
	OutcomeID uint           `gorm:"not null;index:idx_user_market_outcome,priority:3" json:"outcome_id"`
	Shares    float64        `gorm:"type:decimal(20,4);not null" json:"shares"`
	AvgPrice  float64        `gorm:"type:decimal(10,4);not null" json:"avg_price"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Market    Market         `gorm:"foreignKey:MarketID" json:"market,omitempty"`
	Outcome   Outcome        `gorm:"foreignKey:OutcomeID" json:"outcome,omitempty"`
}

// Transaction 交易记录模型
type Transaction struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	Type         string    `gorm:"size:20;not null" json:"type"` // register_bonus, trade_buy, trade_sell, settlement_win, settlement_loss
	Amount       float64   `gorm:"type:decimal(20,2);not null" json:"amount"`
	BalanceAfter float64   `gorm:"type:decimal(20,2);not null" json:"balance_after"`
	OrderID      *uint     `json:"order_id"`
	MarketID     *uint     `json:"market_id"`
	Description  string    `gorm:"size:255" json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

// MarketStatistics 市场统计模型
type MarketStatistics struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	MarketID        uint      `gorm:"uniqueIndex;not null" json:"market_id"`
	TotalTrades     int64     `gorm:"default:0" json:"total_trades"`
	TotalVolume     float64   `gorm:"type:decimal(20,2);default:0" json:"total_volume"`
	UniqueTraders   int64     `gorm:"default:0" json:"unique_traders"`
	LastTradePrice  float64   `gorm:"type:decimal(10,4)" json:"last_trade_price"`
	PriceChange24h  float64   `gorm:"type:decimal(10,4)" json:"price_change_24h"`
	VolumeChange24h float64   `gorm:"type:decimal(20,2)" json:"volume_change_24h"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value     string    `gorm:"type:text;not null" json:"value"`
	Type      string    `gorm:"size:20;not null" json:"type"` // string, number, boolean, json
	Category  string    `gorm:"size:50;not null" json:"category"`
	Label     string    `gorm:"size:100;not null" json:"label"`
	UpdatedAt time.Time `json:"updated_at"`
}
