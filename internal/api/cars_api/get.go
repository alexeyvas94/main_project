package cars_api

import (
	"context"
	"github.com/alexeyvas94/main_project/internal/converter"
	pb "github.com/alexeyvas94/main_project/pkg/car"
)

func (c *CarServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	carObj, err := c.carService.GetCar(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetResponse{
		User: converter.ToProto(carObj),
	}, nil
}
