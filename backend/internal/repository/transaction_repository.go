package repository

import (
	"github.com/huabtc/polygame/backend/internal/model"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create 创建交易记录
func (r *TransactionRepository) Create(tx *model.Transaction) error {
	return r.db.Create(tx).Error
}

// FindByUserID 根据用户 ID 查找交易记录
func (r *TransactionRepository) FindByUserID(userID uint, page, pageSize int) ([]model.Transaction, int64, error) {
	var transactions []model.Transaction
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Transaction{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&transactions).Error

	return transactions, total, err
}
