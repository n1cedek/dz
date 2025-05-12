package service

import (
	"context"
	"dz/auth/model"
	modelRepo "dz/auth/repo/auth/model"
)

type AuthService interface {
	Create(ctx context.Context, info *modelRepo.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.PublicInfo, error)
}