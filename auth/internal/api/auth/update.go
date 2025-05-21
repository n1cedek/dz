package auth2

import (
	"context"
	converter1 "dz/auth/internal/converter"
	desc "dz/auth/pkg/w1"
	"github.com/golang/protobuf/ptypes/empty"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*empty.Empty, error) {
	err := i.authService.Update(ctx, req.Id, converter1.ToInfoFromUpdate(req))
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
