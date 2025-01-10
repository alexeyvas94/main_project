package cars_api

import (
	"context"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *CarServer) Update(ctx context.Context, req *pb.UpdateRequest) (*emptypb.Empty, error) {
	err := c.carService.UpdateCar(ctx, req)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
