package tests

import (
	"context"
	"fmt"
	"github.com/alexeyvas94/main_project/internal/cache"
	"github.com/alexeyvas94/main_project/internal/repository"
	serviceMocks "github.com/alexeyvas94/main_project/internal/repository/mocks"
	"github.com/alexeyvas94/main_project/internal/service/cars_service"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type carServiceMockFunc func(mc *minimock.Controller) repository.CarRepository
	type carCacheMockFunc func(mc *minimock.Controller) cache.Cache

	type args struct {
		ctx context.Context
		req *pb.UpdateRequest
	}

	var (
		ctx   = context.Background()
		mc    = minimock.NewController(t)
		id    = gofakeit.Int64()
		marka = gofakeit.CarMaker()
		model = gofakeit.CarModel()
		vin   = gofakeit.Numerify("###-###-###")
		years = int64(gofakeit.Year())
		role  = pb.Role_ON_SALE

		serviceErr = fmt.Errorf("repo error")

		req = &pb.UpdateRequest{
			Id:    id,
			Marka: wrapperspb.String(marka),
			Model: wrapperspb.String(model),
			VIN:   wrapperspb.String(vin),
			Year:  wrapperspb.Int64(years),
			Role:  role,
		}
	)

	t.Cleanup(func() {
	})

	tests := []struct {
		name         string
		args         args
		err          error
		carRepoMock  carServiceMockFunc
		carCacheMock carCacheMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			carRepoMock: func(mc *minimock.Controller) repository.CarRepository {
				mock := serviceMocks.NewCarRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, req).Return(nil)
				return mock
			},
		},
		{
			name: "repo error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: serviceErr,
			carRepoMock: func(mc *minimock.Controller) repository.CarRepository {
				mock := serviceMocks.NewCarRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, req).Return(serviceErr)
				return mock
			},
		},
		{
			name: "validation error case",
			args: args{
				ctx: ctx,
				req: &pb.UpdateRequest{
					Id:    id,
					Marka: nil, // Некорректные данные
				},
			},
			err: fmt.Errorf("marka, model, vin and year are required"),
			carRepoMock: func(mc *minimock.Controller) repository.CarRepository {
				mock := serviceMocks.NewCarRepositoryMock(mc)
				return mock // Мок не вызывается, так как проверка завершится раньше
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			carRepoMock := tt.carRepoMock(mc)
			carCacheMock := tt.carCacheMock(mc)
			carService := cars_service.NewCarService(carRepoMock, carCacheMock)

			err := carService.UpdateCar(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
		})
	}
}
