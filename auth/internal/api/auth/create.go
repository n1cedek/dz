package auth2

import (
	"context"
	converter1 "dz/auth/internal/converter"
	desc "dz/auth/pkg/w1"
	"log"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.authService.Create(ctx, converter1.ToRepoUserFromCreateRequest(req))
	if err != nil {
		return nil, err
	}
	log.Printf("inserted user id: %d", id)
	return &desc.CreateResponse{Id: id}, nil
}
