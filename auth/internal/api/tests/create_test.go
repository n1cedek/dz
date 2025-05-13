package tests

import (
	"context"
	auth2 "dz/auth/internal/api/auth"
	"dz/auth/internal/repo/auth/model"
	"dz/auth/internal/service"
	serviceMock "dz/auth/internal/service/mocks"
	desc "dz/auth/pkg/w1"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService
	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name  = gofakeit.BeerName()
		email = gofakeit.Email()
		pas   = gofakeit.Letter()

		//roles = []string{"user", "admin"}
		//role  = model.Role{Value: roles[gofakeit.Number(0, 1)]}

		errorSer = fmt.Errorf("service error")

		req = &desc.CreateRequest{
			Name:     name,
			Email:    email,
			Password: pas,
		}

		// Объект пользователя, ожидаемый в вызове метода Create сервиса
		info = &model.User{
			ID:       0, // На входе всегда 0! Всегда 0 для нового пользователя
			Name:     name,
			Email:    email,
			Password: pas,
			Role:     model.RoleUser,
		}

		res = &desc.CreateResponse{Id: id}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
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
				mock.CreateMock.Expect(ctx, info).Return(id, nil)
				return mock
			},
		}, {
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  errorSer,
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMock.NewAuthServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(0, errorSer)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			authServMock := tt.authServiceMock(mc)
			api := auth2.NewImplementation(authServMock)

			nId, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, nId)
		})
	}

}
