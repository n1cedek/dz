package tests

import (
	"context"
	"dz/auth/internal/model"
	"dz/auth/internal/repo"
	repoMock "dz/auth/internal/repo/mocks"
	auth1 "dz/auth/internal/service/auth"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdate(t *testing.T) {
	type authServiceMockFunc func(mc *minimock.Controller) repo.AuthRepo

	type args struct {
		ctx  context.Context
		id   int64
		info *model.User
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)
		id  = gofakeit.Int64()

		info = &model.User{
			Name:     gofakeit.Name(),
			Email:    gofakeit.Email(),
			Password: gofakeit.BeerName(),
			Role:     model.RoleUser,
		}
		resErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		arg             args
		err             error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "success case",
			arg: args{
				ctx:  ctx,
				id:   id,
				info: info,
			},
			err: nil,
			authServiceMock: func(mc *minimock.Controller) repo.AuthRepo {
				mock := repoMock.NewAuthRepoMock(mc)
				mock.UpdateMock.Expect(ctx, id, info).Return(nil)
				return mock
			},
		}, {
			name: "error case",
			arg: args{
				ctx:  ctx,
				id:   id,
				info: info,
			},
			err: resErr,
			authServiceMock: func(mc *minimock.Controller) repo.AuthRepo {
				mock := repoMock.NewAuthRepoMock(mc)
				mock.UpdateMock.Expect(ctx, id, info).Return(resErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			serviceMock := tt.authServiceMock(mc)
			api := auth1.NewServ(serviceMock)

			err := api.Update(tt.arg.ctx, tt.arg.id, tt.arg.info)
			require.Equal(t, tt.err, err)
		})
	}
}
