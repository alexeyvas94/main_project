package converter

import (
	"github.com/alexeyvas94/main_project/internal/models"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Role(role string) pb.Role {
	switch role {
	case "ON_SALE":
		return pb.Role_ON_SALE
	case "BOOKED":
		return pb.Role_BOOKED
	default:
		return pb.Role_ON_SALE // Добавьте значение UNKNOWN в вашем protobuf определении, если его нет.
	}
}
func ToProto(c *models.Car) *pb.IsCar {
	return &pb.IsCar{
		Id:        c.ID,
		Marka:     c.Marka,
		Model:     c.Model,
		VIN:       c.VIN,
		Year:      c.Year,
		Role:      Role(c.Role),
		CreatedAt: timestamppb.New(c.CreatedAt), // Конвертация time.Time в *timestamppb.Timestamp.
		UpdatedAt: timestamppb.New(c.UpdatedAt),
	}
}
