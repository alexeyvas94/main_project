package repository

import (
	"context"
	"github.com/alexeyvas94/main_project/internal/models"
	pb "github.com/alexeyvas94/main_project/pkg/car"
)

type CarRepository interface {
	Create(ctx context.Context, req *pb.CreateRequest) (int64, error)
	Get(ctx context.Context, id int64) (*models.Car, error)
	Update(ctx context.Context, req *pb.UpdateRequest) error
	Delete(ctx context.Context, id int64) error
}
