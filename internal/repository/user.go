package repository

import (
	"context"
	"godemo/internal/model"
	"godemo/internal/wire/provider"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	List(ctx context.Context, offset, limit int, keyword string) ([]*model.User, int64, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	db provider.MainDB // 主库
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db provider.MainDB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create 创建用户
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.DB.Create(user).Error
}

// FindByID 根据ID查询用户
func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查询用户
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表
func (r *userRepository) List(ctx context.Context, offset, limit int, keyword string) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.DB.Model(&model.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?",
			"%"+keyword+"%",
			"%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
