package auth1

import (
	"context"
	"log"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.authRepo.Delete(ctx, id)
	if err != nil {
		log.Printf("failed to delete user with id %d: %v", id, err)
		return err
	}
	return nil
}
