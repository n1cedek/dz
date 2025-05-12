package auth1

import (
	"context"
	modelRepo "dz/auth/repo/auth/model"
)

func (s *serv) Create(ctx context.Context, info *modelRepo.User) (int64, error){
	id, err := s.authRepo.Create(ctx, info)
	if err != nil {
		return 0, err
	}
	return id,nil
}
