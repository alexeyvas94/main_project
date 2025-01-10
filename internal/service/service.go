package service

import (
	"context"
	"github.com/alexeyvas94/main_project/internal/models"
	pb "github.com/alexeyvas94/main_project/pkg/car"
)

type CarService interface {
	CreateCar(ctx context.Context, req *pb.CreateRequest) (int64, error)
	GetCar(ctx context.Context, id int64) (*models.Car, error)
	UpdateCar(ctx context.Context, req *pb.UpdateRequest) error
	DeleteCar(ctx context.Context, id int64) error
}
