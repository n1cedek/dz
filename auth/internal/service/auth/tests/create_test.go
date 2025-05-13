package tests

import (
	"context"
	"dz/auth/internal/repo"
	"dz/auth/internal/repo/auth/model"
	repoMock "dz/auth/internal/repo/mocks"
	"dz/auth/internal/service/auth"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type authRepoMockFunc func(mc *minimock.Controller) repo.AuthRepo

	type args struct {
		ctx context.Context
		req *model.User
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.BeerName()
		email = gofakeit.Email()
		pas   = gofakeit.Letter()

		roles = []string{"user", "admin"}
		role  = model.Role{Value: roles[gofakeit.Number(0, 1)]}

		errorSer = fmt.Errorf("service error")

		req = &model.User{
			ID:       0, // На входе всегда 0! Всегда 0 для нового пользователя
			Name:     name,
			Email:    email,
			Password: pas,
			Role:     role,
		}
	)

	tests := []struct {
		name         string
		args         args
		want         int64
		err          error
		authRepoMock authRepoMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			authRepoMock: func(mc *minimock.Controller) repo.AuthRepo {
				mock := repoMock.NewAuthRepoMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
		}, {
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  errorSer,
			authRepoMock: func(mc *minimock.Controller) repo.AuthRepo {
				mock := repoMock.NewAuthRepoMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, errorSer)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServMock := tt.authRepoMock(mc)
			api := auth1.NewMockService(authServMock)

			nId, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, nId)
		})
	}

}
