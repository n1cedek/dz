package repo

import (
	"context"
	"dz/auth/internal/model"
	modelRepo "dz/auth/internal/repo/auth/model"
)

type AuthRepo interface {
	Create(ctx context.Context, info *modelRepo.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.PublicInfo, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, info *model.User) error
}
