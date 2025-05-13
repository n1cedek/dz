package auth1

import (
	"context"
	"dz/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.PublicInfo, error) {
	noteObj, err := s.authRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return noteObj, nil
}
