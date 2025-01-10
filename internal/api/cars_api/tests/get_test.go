package tests

import (
	"context"
	"fmt"
	"github.com/alexeyvas94/main_project/internal/api/cars_api"
	"github.com/alexeyvas94/main_project/internal/converter"
	"github.com/alexeyvas94/main_project/internal/models"
	"github.com/alexeyvas94/main_project/internal/service"
	serviceMocks "github.com/alexeyvas94/main_project/internal/service/mocks"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	t.Parallel()

	// Тип для генерации мока
	type carServiceMockFunc func(mc *minimock.Controller) service.CarService
	type args struct {
		ctx context.Context
		req *pb.GetRequest
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

		serviceErr = fmt.Errorf("service error")
		car        = models.Car{
			Marka:     marka,
			Model:     mode,
			VIN:       vin,
			Year:      int64(years),
			Role:      string(role),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		req              = &pb.GetRequest{Id: id}
		expectedResponse = &pb.GetResponse{User: converter.ToProto(&car)}
	)

	t.Cleanup(func() {
	})
	// Cleanup для завершения минимок-контроллера

	tests := []struct {
		name           string
		args           args
		expectedResult *pb.GetResponse
		expectedError  error
		carServiceMock carServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			expectedResult: expectedResponse,
			expectedError:  nil,
			carServiceMock: func(mc *minimock.Controller) service.CarService {
				mock := serviceMocks.NewCarServiceMock(mc)
				mock.GetCarMock.Expect(ctx, id).Return(&car, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			expectedResult: nil,
			expectedError:  serviceErr,
			carServiceMock: func(mc *minimock.Controller) service.CarService {
				mock := serviceMocks.NewCarServiceMock(mc)
				mock.GetCarMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt // фиксируем tt для каждого `t.Run`
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Создаем мок сервиса
			carServiceMock := tt.carServiceMock(mc)
			carApi := cars_api.NewCarServer(carServiceMock)

			// Вызываем метод Get
			actualResult, actualError := carApi.Get(tt.args.ctx, tt.args.req)

			// Проверяем ошибки
			require.Equal(t, tt.expectedError, actualError, "ошибка не совпадает с ожидаемой")
			// Проверяем результаты
			require.Equal(t, tt.expectedResult, actualResult, "ответ не совпадает с ожидаемым")
		})
	}
}
