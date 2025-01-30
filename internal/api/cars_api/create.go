package cars_api

import (
	"context"
	pb "github.com/alexeyvas94/main_project/pkg/car"
)

func (c *CarServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	id, err := c.carService.CreateCar(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateResponse{
		Id: id,
	}, nil
}
