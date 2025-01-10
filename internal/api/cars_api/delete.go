package cars_api

import (
	"context"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *CarServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	err := c.carService.DeleteCar(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
