package auth2

import (
	"context"
	desc "dz/auth/pkg/w1"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	err := i.authService.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	log.Printf("deleted user id: %d", req.Id)
	return &empty.Empty{}, nil
}
