package tests

import (
	"context"
	auth2 "dz/auth/internal/api/auth"
	"dz/auth/internal/converter"
	"dz/auth/internal/model"
	"dz/auth/internal/service"
	serviceMock "dz/auth/internal/service/mocks"
	desc "dz/auth/pkg/w1"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id     = gofakeit.Int64()
		name   = gofakeit.Name()
		email  = gofakeit.Email()
		resErr = fmt.Errorf("service error")

		req = &desc.GetRequest{Id: id}

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
				req: req,
			},
			want: res,
			err:  nil,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMock.NewAuthServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(res, nil)
				return mock
			},
		}, {
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  resErr,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMock.NewAuthServiceMock(mc)
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
			api := auth2.NewImplementation(serviceAMock)

			resu, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, converter.ToGetResponse(tt.want), resu)
		})
	}
}
