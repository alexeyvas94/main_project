package tests

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"testing"

	"github.com/alexeyvas94/main_project/internal/api/cars_api"
	"github.com/alexeyvas94/main_project/internal/service"
	serviceMocks "github.com/alexeyvas94/main_project/internal/service/mocks"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	t.Parallel()
	type carServiceMockFunc func(mc *minimock.Controller) service.CarService

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
		years = gofakeit.Year()
		role  = pb.Role_ON_SALE

		serviceErr = fmt.Errorf("service error")
		req        = &pb.UpdateRequest{
			Id:    id,
			Marka: wrapperspb.String(marka),
			Model: wrapperspb.String(model),
			VIN:   wrapperspb.String(vin),
			Year:  wrapperspb.Int64(int64(years)),
			Role:  role,
		}
		res = &emptypb.Empty{}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name           string
		args           args
		want           *emptypb.Empty
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
				mock.UpdateCarMock.Expect(ctx, req).Return(nil)
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
				mock.UpdateCarMock.Expect(ctx, req).Return(serviceErr)
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

			result, err := carApi.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, result)
		})
	}
}
