package tests

import (
	"context"
	"fmt"
	"github.com/alexeyvas94/main_project/internal/repository"
	serviceMocks "github.com/alexeyvas94/main_project/internal/repository/mocks"
	"github.com/alexeyvas94/main_project/internal/service/cars_service"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDelete(t *testing.T) {
	t.Parallel()
	type carServiceMockFunc func(mc *minimock.Controller) repository.CarRepository

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)
		id  = gofakeit.Int64()

		serviceErr = fmt.Errorf("repo error")
	)
	t.Cleanup(func() {
	})

	tests := []struct {
		name        string
		args        args
		err         error
		carRepoMock carServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: id,
			},
			err: nil,
			carRepoMock: func(mc *minimock.Controller) repository.CarRepository {
				mock := serviceMocks.NewCarRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "repo error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			err: serviceErr,
			carRepoMock: func(mc *minimock.Controller) repository.CarRepository {
				mock := serviceMocks.NewCarRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			carServiceMock := tt.carRepoMock(mc)
			carService := cars_service.NewCarService(carServiceMock)
			err := carService.DeleteCar(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
