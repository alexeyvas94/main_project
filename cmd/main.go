package main

import (
	"context"
	"github.com/alexeyvas94/main_project/internal/api/cars_api"
	"github.com/alexeyvas94/main_project/internal/repository/cars"
	"github.com/alexeyvas94/main_project/internal/service/cars_service"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const dbDSN = "host=localhost port=54321 dbname=car_base user=car-user password=car-password sslmode=disable"

func main() {

	ctx := context.Background()
	// тут должно быть переменное окружение
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// создаем пул соеденений с базой данных
	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}
	defer pool.Close()
	// создаем репозиторий
	carRepo := cars.NewPostgresCarRepository(pool)
	// создаем сервисный слой
	carService := cars_service.NewCarService(carRepo)
	// создаем API слой
	carApi := cars_api.NewCarServer(carService)
	s := grpc.NewServer()
	pb.RegisterCarServer(s, carApi)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
