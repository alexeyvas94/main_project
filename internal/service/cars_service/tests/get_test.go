package tests

import (
	"context"
	"fmt"
	"github.com/alexeyvas94/main_project/internal/cache"
	"github.com/alexeyvas94/main_project/internal/models"
	"github.com/alexeyvas94/main_project/internal/repository"
	serviceMocks "github.com/alexeyvas94/main_project/internal/repository/mocks"
	"github.com/alexeyvas94/main_project/internal/service/cars_service"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type carServiceMockFunc func(mc *minimock.Controller) repository.CarRepository
	type carCacheMockFunc func(mc *minimock.Controller) cache.Cache

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx   = context.Background()
		mc    = minimock.NewController(t)
		id    = gofakeit.Int64()
		marka = gofakeit.CarMaker()
		model = gofakeit.CarModel()
		vin   = gofakeit.Numerify("###-###-###")
		years = gofakeit.Year()
		role  = pb.Role_ON_SALE

		serviceErr = fmt.Errorf("repo error")
		car        = models.Car{
			ID:    id,
			Marka: marka,
			Model: model,
			VIN:   vin,
			Year:  int64(years),
			Role:  role.String(),
		}
	)

	t.Cleanup(func() {
	})

	tests := []struct {
		name         string
		args         args
		want         *models.Car
		err          error
		carRepoMock  carServiceMockFunc
		carCacheMock carCacheMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: &car,
			err:  nil,
			carRepoMock: func(mc *minimock.Controller) repository.CarRepository {
				mock := serviceMocks.NewCarRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(&car, nil)
				return mock
			},
		},
		{
			name: "repo error case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: &models.Car{},
			err:  serviceErr,
			carRepoMock: func(mc *minimock.Controller) repository.CarRepository {
				mock := serviceMocks.NewCarRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(&models.Car{}, serviceErr)
				return mock
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

			model, err := carService.GetCar(tt.args.ctx, id)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, model)
		})
	}
}
