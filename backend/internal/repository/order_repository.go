package repository

import (
	"errors"

	"github.com/huabtc/polygame/backend/internal/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create 创建订单
func (r *OrderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

// FindByID 根据 ID 查找订单
func (r *OrderRepository) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Market").Preload("Outcome").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

// FindByUserID 根据用户 ID 查找订单列表
func (r *OrderRepository) FindByUserID(userID uint, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Market").Preload("Outcome").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

// Update 更新订单
func (r *OrderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

// UpdateStatus 更新订单状态
func (r *OrderRepository) UpdateStatus(orderID uint, status string) error {
	return r.db.Model(&model.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}

// Cancel 取消订单
func (r *OrderRepository) Cancel(orderID uint, userID uint) error {
	return r.db.Model(&model.Order{}).
		Where("id = ? AND user_id = ? AND status = ?", orderID, userID, "pending").
		Update("status", "cancelled").Error
}
