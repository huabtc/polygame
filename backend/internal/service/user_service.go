package service

import (
	"errors"

	"github.com/huabtc/polygame/backend/config"
	"github.com/huabtc/polygame/backend/internal/middleware"
	"github.com/huabtc/polygame/backend/internal/model"
	"github.com/huabtc/polygame/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
	txRepo   *repository.TransactionRepository
	cfg      *config.Config
}

func NewUserService(userRepo *repository.UserRepository, txRepo *repository.TransactionRepository, cfg *config.Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		txRepo:   txRepo,
		cfg:      cfg,
	}
}

// Register 用户注册
func (s *UserService) Register(username, email, password string) (*model.User, string, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return nil, "", errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, "", errors.New("email already exists")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// 创建用户
	user := &model.User{
		Username:       username,
		Email:          email,
		PasswordHash:   string(hashedPassword),
		VirtualBalance: 10000, // 初始虚拟积分 10000
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	// 记录注册奖励交易
	tx := &model.Transaction{
		UserID:       user.ID,
		Type:         "register_bonus",
		Amount:       10000,
		BalanceAfter: 10000,
		Description:  "Registration bonus",
	}
	_ = s.txRepo.Create(tx)

	// 生成 JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.IsAdmin, s.cfg)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (*model.User, string, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	// 生成 JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.IsAdmin, s.cfg)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetProfile 获取用户信息
func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

// UpdateProfile 更新用户信息
func (s *UserService) UpdateProfile(userID uint, avatar string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if avatar != "" {
		user.Avatar = avatar
	}

	return s.userRepo.Update(user)
}

// GetBalance 获取用户余额
func (s *UserService) GetBalance(userID uint) (float64, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return 0, err
	}
	return user.VirtualBalance, nil
}

// ListUsers 获取用户列表（管理员）
func (s *UserService) ListUsers(page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.List(page, pageSize)
}
