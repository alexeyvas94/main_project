package cars_service

import (
	"context"
	"fmt"
	"github.com/alexeyvas94/main_project/internal/cache"
	"github.com/alexeyvas94/main_project/internal/models"
	"github.com/alexeyvas94/main_project/internal/repository"
	"github.com/alexeyvas94/main_project/internal/service"
	pb "github.com/alexeyvas94/main_project/pkg/car"
	"time"
)

type CarService struct {
	carRepository repository.CarRepository
	cache         cache.Cache // Интегрируем интерфейс кэша
}

func NewCarService(carRepo repository.CarRepository, carCache cache.Cache) service.CarService {
	return &CarService{
		carRepository: carRepo,
		cache:         carCache,
	}
}

func (s *CarService) CreateCar(ctx context.Context, req *pb.CreateRequest) (int64, error) {
	if req.Marka == "" {
		return 0, fmt.Errorf("marka is required")
	}

	// Создаем автомобиль в базе
	carID, err := s.carRepository.Create(ctx, req)
	if err != nil {
		return 0, err
	}

	// Формируем объект автомобиля с правильными значениями
	car := &models.Car{
		ID:        carID,
		Marka:     req.Marka,
		Model:     req.Model,
		VIN:       req.VIN,
		Year:      req.Year,
		Role:      req.Role.String(), // Убедитесь, что Role передано как строка
		CreatedAt: time.Now().UTC(),  // Текущая дата и время
		UpdatedAt: time.Now().UTC(),  // Текущая дата и время
	}

	// Сохраняем автомобиль в кэш
	_ = s.cache.SetCarInCache(car)

	return carID, nil
}

func (s *CarService) GetCar(ctx context.Context, id int64) (*models.Car, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid car ID")
	}

	// Проверяем кэш
	car, err := s.cache.GetCarFromCache(id)
	if err == nil && car != nil {
		// Если нашли в кэше, возвращаем
		return car, nil
	}

	// Если данных нет в кэше, обращаемся к базе данных
	car, err = s.carRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Добавляем данные в кэш
	_ = s.cache.SetCarInCache(car) // Ошибки кэша не прерывают выполнение

	return car, nil
}

func (s *CarService) UpdateCar(ctx context.Context, req *pb.UpdateRequest) error {
	if req.Marka == nil || req.Model == nil || req.VIN == nil || req.Year == nil {
		return fmt.Errorf("marka, model, vin and year are required")
	}

	// Обновляем данные в базе
	err := s.carRepository.Update(ctx, req)
	if err != nil {
		return err
	}

	// Удаляем устаревшие данные из кэша (или можно обновить кэш сразу)
	_ = s.cache.SetCarInCache(&models.Car{
		ID:    req.Id,
		Marka: req.Marka.Value,
		Model: req.Model.Value,
		VIN:   req.VIN.Value,
		Year:  req.Year.Value,
	}) // Обновляем кэш после изменения

	return nil
}

func (s *CarService) DeleteCar(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid car ID")
	}

	// Удаляем из базы данных
	err := s.carRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Удаляем данные из кэша
	_ = s.cache.SetCarInCache(nil) // Инвалидация кэша

	return nil
}
