package tests

import (
	"context"
	"fmt"
	"github.com/alexeyvas94/main_project/internal/cache"
	cacheMocks "github.com/alexeyvas94/main_project/internal/cache/mocks"
	"github.com/alexeyvas94/main_project/internal/models"
	"github.com/alexeyvas94/main_project/internal/repository"
	repoMocks "github.com/alexeyvas94/main_project/internal/repository/mocks"
	"github.com/alexeyvas94/main_project/internal/service/cars_service"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type carRepoMockFunc func(mc *minimock.Controller) repository.CarRepository
	type carCacheMockFunc func(mc *minimock.Controller) cache.Cache

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
		mc.Finish() // Завершение контроллера моков
	})

	tests := []struct {
		name         string
		args         args
		want         int64
		err          error
		carRepoMock  carRepoMockFunc
		carCacheMock carCacheMockFunc
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
				mock := repoMocks.NewCarRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
			carCacheMock: func(mc *minimock.Controller) cache.Cache {
				mock := cacheMocks.NewCacheMock(mc)
				mock.SetCarInCacheMock.Expect(&models.Car{
					ID:        id,
					Marka:     req.Marka,
					Model:     req.Model,
					VIN:       req.VIN,
					Year:      req.Year,
					Role:      req.Role.String(), // Убедитесь, что Role преобразован в строку
					CreatedAt: time.Now().UTC(),  // Устанавливаем текущую метку времени
					UpdatedAt: time.Now().UTC(),  // Устанавливаем текущую метку времени
				}).Return(nil)
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
				mock := repoMocks.NewCarRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, req).Return(0, serviceErr)
				return mock
			},
			carCacheMock: func(mc *minimock.Controller) cache.Cache {
				return cacheMocks.NewCacheMock(mc)
			},
		},
		{
			name: "validation error case",
			args: args{
				ctx: ctx,
				req: &pb.CreateRequest{}, // Пустой запрос
			},
			want: 0,
			err:  fmt.Errorf("marka is required"),
			carRepoMock: func(mc *minimock.Controller) repository.CarRepository {
				return repoMocks.NewCarRepositoryMock(mc)
			},
			carCacheMock: func(mc *minimock.Controller) cache.Cache {
				return cacheMocks.NewCacheMock(mc)
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

			newID, err := carService.CreateCar(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
