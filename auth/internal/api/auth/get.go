package auth2

import (
	"context"
	desc "dz/auth/pkg/w1"
	"log"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest)(*desc.GetResponse, error) {
	noteObj, err := i.authService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	log.Printf("id: %d, name: %s, email: %s", noteObj.ID, noteObj.Name, noteObj.Email)
	return &desc.GetResponse{
		Id:    noteObj.ID,
		Name:  noteObj.Name,
		Email: noteObj.Email,
	}, nil
}
