package cars

import (
	"context"
	"fmt"
	"github.com/alexeyvas94/main_project/internal/models"
	"github.com/alexeyvas94/main_project/internal/repository"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type postgresCarRepository struct {
	db *pgxpool.Pool
}

// Конструктор репозитория
func NewPostgresCarRepository(db *pgxpool.Pool) repository.CarRepository {
	return &postgresCarRepository{db: db}
}

func (r *postgresCarRepository) Create(ctx context.Context, req *pb.CreateRequest) (int64, error) {
	car := models.Car{
		Marka: req.Marka,
		Model: req.Model,
		VIN:   req.VIN,
		Year:  req.Year,
		Role:  req.Role.String(), // Преобразование к типу модели
	}

	var carID int64
	err := r.db.QueryRow(ctx,
		"INSERT INTO cars (marka, model, vin, year, role) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		car.Marka, car.Model, car.VIN, car.Year, car.Role).Scan(&carID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert cars: %v", err)
	}
	return carID, nil
}
func (r *postgresCarRepository) Get(ctx context.Context, id int64) (*models.Car, error) {

	var car models.Car
	err := r.db.QueryRow(ctx,
		"SELECT id, marka, model, vin, year, role, created_at, updated_at FROM cars WHERE id = $1",
		id).Scan(&car.ID, &car.Marka, &car.Model, &car.VIN, &car.Year, &car.Role, &car.CreatedAt, &car.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("car with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get car: %v", err)
	}
	return &car, nil
	//return &pb.GetResponse{
	//	User: converter.ToProto(&car),
	//}, nil
}
func (r *postgresCarRepository) Update(ctx context.Context, req *pb.UpdateRequest) error {

	now := time.Now()
	// Динамическое формирование запроса
	query := "UPDATE cars SET updated_at = $1"
	args := []interface{}{now}
	index := 2 // Следующий индекс для SQL-запроса

	if req.Marka != nil {
		query += fmt.Sprintf(", marka = $%d", index)
		args = append(args, req.Marka.Value)
		index++
	}
	if req.Model != nil {
		query += fmt.Sprintf(", model = $%d", index)
		args = append(args, req.Model.Value)
		index++
	}
	if req.VIN != nil {
		query += fmt.Sprintf(", vin = $%d", index)
		args = append(args, req.VIN.Value)
		index++
	}
	if req.Year != nil {
		query += fmt.Sprintf(", year = $%d", index)
		args = append(args, req.Year.Value)
		index++
	}
	query += fmt.Sprintf(" WHERE id = $%d", index)
	args = append(args, req.Id)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update car: %v", err)
	}

	// Проверка на количество изменённых строк
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no car found with id: %d", req.Id)
	}
	return nil
}

func (r *postgresCarRepository) Delete(ctx context.Context, id int64) error {

	_, err := r.db.Exec(ctx, "DELETE FROM cars WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete cars: %v", err)
	}
	return nil
}
