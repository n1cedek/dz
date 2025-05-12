package tests

import (
	"context"
	"dz/auth/model"
	"dz/auth/repo"
	repoMock "dz/auth/repo/mocks"
	auth1 "dz/auth/service/auth"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	type authServiceMockFunc func(mc *minimock.Controller) repo.AuthRepo

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id     = gofakeit.Int64()
		name   = gofakeit.Name()
		email  = gofakeit.Email()
		resErr = fmt.Errorf("service error")

		res = &model.PublicInfo{
			ID:    id,
			Name:  name,
			Email: email,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *model.PublicInfo
		err             error
		authServiceMock authServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: res,
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) repo.AuthRepo {
				mock := repoMock.NewAuthRepoMock(mc)
				mock.GetMock.Expect(ctx, id).Return(res, nil)
				return mock
			},
		}, {
			name: "error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: nil,
			err:  resErr,
			authServiceMock: func(mc *minimock.Controller) repo.AuthRepo {
				mock := repoMock.NewAuthRepoMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, resErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			serviceAMock := tt.authServiceMock(mc)
			api := auth1.NewMockService(serviceAMock)

			resu, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resu)
		})
	}
}
