package cars_api

import (
	"github.com/alexeyvas94/main_project/internal/service"
	pb "github.com/alexeyvas94/main_project/pkg/car"
)

type CarServer struct {
	pb.UnimplementedCarServer
	carService service.CarService
}

func NewCarServer(carService service.CarService) *CarServer {
	return &CarServer{
		carService: carService,
	}
}
