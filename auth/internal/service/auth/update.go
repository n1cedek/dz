package auth1

import (
	"context"
	"dz/auth/internal/model"
	"log"
)

func (s *serv) Update(ctx context.Context, id int64, info *model.User) error {
	err := s.authRepo.Update(ctx, id, info)
	if err != nil {
		log.Printf("failed to update user with id %d: %v", id, err)
		return err
	}
	return nil
}
