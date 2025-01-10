package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/alexeyvas94/main_project/internal/api/cars_api"
	"github.com/alexeyvas94/main_project/internal/service"
	serviceMocks "github.com/alexeyvas94/main_project/internal/service/mocks"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type carServiceMockFunc func(mc *minimock.Controller) service.CarService

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

		serviceErr = fmt.Errorf("service error")

		req = &pb.CreateRequest{
			Marka: marka,
			Model: mode,
			VIN:   vin,
			Year:  int64(years),
			Role:  role,
		}
		res = &pb.CreateResponse{
			Id: id,
		}
	)
	t.Cleanup(func() {
	})

	tests := []struct {
		name           string
		args           args
		want           *pb.CreateResponse
		err            error
		carServiceMock carServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			carServiceMock: func(mc *minimock.Controller) service.CarService {
				mock := serviceMocks.NewCarServiceMock(mc)
				mock.CreateCarMock.Expect(ctx, req).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			carServiceMock: func(mc *minimock.Controller) service.CarService {
				mock := serviceMocks.NewCarServiceMock(mc)
				mock.CreateCarMock.Expect(ctx, req).Return(0, serviceErr)
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

			newID, err := carApi.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
