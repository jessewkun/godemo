package service

import (
	"context"
	"time"

	"godemo/internal/dto"
	"godemo/internal/model"
	"godemo/internal/repository"
	"godemo/internal/wire/provider"
)

// UserService 用户服务
type UserService struct {
	repo  repository.UserRepository // 用户仓储
	cache provider.MainCache        // 缓存连接
}

// NewUserService 创建用户服务
func NewUserService(repo repository.UserRepository, cache provider.MainCache) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, req *dto.UserCreateRequest) (*dto.UserCreateResponse, error) {
	user := &model.User{
		Username: req.Username,
		Password: req.Password, // 注意：实际应用中需要对密码进行加密
		Email:    req.Email,
	}

	// 使用仓储创建用户
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 使用缓存连接缓存用户信息
	key := "user:" + user.Username
	if err := s.cache.Set(ctx, key, user, time.Hour); err != nil {
		// 缓存失败不影响主流程，只记录日志
		// TODO: 记录日志
	}

	return &dto.UserCreateResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		CreateAt: user.CreatedAt.String(),
	}, nil
}

// List 获取用户列表
func (s *UserService) List(ctx context.Context, req *dto.UserListRequest) (*dto.UserListResponse, error) {
	// 使用仓储获取用户列表
	users, total, err := s.repo.List(ctx, (req.Page-1)*req.PageSize, req.PageSize, req.Keyword)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]dto.UserCreateResponse, 0, len(users))
	for _, user := range users {
		list = append(list, dto.UserCreateResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			CreateAt: user.CreatedAt.String(),
		})
	}

	return &dto.UserListResponse{
		Total: total,
		List:  list,
	}, nil
}
