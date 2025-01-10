package tests

import (
	"context"
	"fmt"
	"github.com/alexeyvas94/main_project/internal/repository"
	serviceMocks "github.com/alexeyvas94/main_project/internal/repository/mocks"
	"github.com/alexeyvas94/main_project/internal/service/cars_service"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type carServiceMockFunc func(mc *minimock.Controller) repository.CarRepository

	type args struct {
		ctx context.Context
		req *pb.CreateRequest
	}

	var (
		ctx   = context.Background()
		mc    = minimock.NewController(t)
		id    = gofakeit.Int64()
		marka = gofakeit.CarMaker()
		mode  = gofakeit.CarModel()
		vin   = gofakeit.Numerify("###-###-###")
		years = gofakeit.Year()
		role  = pb.Role_ON_SALE

		serviceErr = fmt.Errorf("repo error")

		req = &pb.CreateRequest{
			Marka: marka,
			Model: mode,
			VIN:   vin,
			Year:  int64(years),
			Role:  role,
		}
	)
	t.Cleanup(func() {
	})

	tests := []struct {
		name        string
		args        args
		want        int64
		err         error
		carRepoMock carServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			carRepoMock: func(mc *minimock.Controller) repository.CarRepository {
				mock := serviceMocks.NewCarRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
		},
		{
			name: "repo error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  serviceErr,
			carRepoMock: func(mc *minimock.Controller) repository.CarRepository {
				mock := serviceMocks.NewCarRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, serviceErr)
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

			newID, err := carService.CreateCar(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
