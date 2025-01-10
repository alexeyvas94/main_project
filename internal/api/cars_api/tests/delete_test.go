package tests

import (
	"context"
	"fmt"
	"github.com/alexeyvas94/main_project/internal/api/cars_api"
	"github.com/alexeyvas94/main_project/internal/service"
	serviceMocks "github.com/alexeyvas94/main_project/internal/service/mocks"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestDelete(t *testing.T) {
	// Тип для генерации мока
	type carServiceMockFunc func(mc *minimock.Controller) service.CarService
	type args struct {
		ctx context.Context
		req *pb.DeleteRequest
	}
	var (
		id               = gofakeit.Int64()
		mc               = minimock.NewController(t)
		ctx              = context.Background()
		serviceErr       = fmt.Errorf("service error")
		req              = &pb.DeleteRequest{Id: id}
		expectedResponse = &emptypb.Empty{}
	)
	defer mc.Finish()
	tests := []struct {
		name           string
		args           args
		expectedResult *emptypb.Empty
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
				mock.DeleteCarMock.Expect(ctx, id).Return(nil)
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
				mock.DeleteCarMock.Expect(ctx, id).Return(serviceErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			carServiceMock := tt.carServiceMock(mc)
			carApi := cars_api.NewCarServer(carServiceMock)

			result, err := carApi.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.expectedError, err)
			require.Equal(t, tt.expectedResult, result)
		})
	}
}
