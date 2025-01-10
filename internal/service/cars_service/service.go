package cars_service

import (
	"context"
	"fmt"
	"github.com/alexeyvas94/main_project/internal/models"
	"github.com/alexeyvas94/main_project/internal/repository"
	"github.com/alexeyvas94/main_project/internal/service"
	pb "github.com/alexeyvas94/main_project/pkg/car"
)

type CarService struct {
	carRepository repository.CarRepository
}

func NewCarService(carRepo repository.CarRepository) service.CarService {
	return &CarService{carRepository: carRepo}
}

func (s *CarService) CreateCar(ctx context.Context, req *pb.CreateRequest) (int64, error) {
	if req.Marka == "" {
		return 0, fmt.Errorf("marka is required")
	}
	return s.carRepository.Create(ctx, req)
}

func (s *CarService) GetCar(ctx context.Context, id int64) (*models.Car, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid users ID")
	}
	return s.carRepository.Get(ctx, id)
}

func (s *CarService) UpdateCar(ctx context.Context, req *pb.UpdateRequest) error {
	if req.Marka == nil || req.Model == nil || req.VIN == nil || req.Year == nil {
		return fmt.Errorf("marka, model, vin and year are required")
	}
	return s.carRepository.Update(ctx, req)
}
func (s *CarService) DeleteCar(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid users ID")
	}
	return s.carRepository.Delete(ctx, id)
}
