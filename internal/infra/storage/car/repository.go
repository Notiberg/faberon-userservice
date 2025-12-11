package car

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/m04kA/SMC-UserService/internal/domain"
	userservice "github.com/m04kA/SMC-UserService/internal/service/user"
	"github.com/m04kA/SMC-UserService/pkg/psqlbuilder"
)

var (
	ErrCreateCar  = errors.New("failed to create car in database")
	ErrGetCar     = errors.New("failed to get car from database")
	ErrUpdateCar  = errors.New("failed to update car in database")
	ErrDeleteCar  = errors.New("failed to delete car from database")
	ErrBuildQuery = errors.New("failed to build SQL query")
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(executor *sqlx.DB) *Repository {
	return &Repository{
		db: executor,
	}
}

// Create создает новый автомобиль и возвращает его с присвоенным ID
func (r *Repository) Create(ctx context.Context, car *domain.Car) (*domain.Car, error) {
	query, args, err := psqlbuilder.Insert("cars").
		Columns("user_id", "brand", "model", "license_plate", "color", "size", "is_selected").
		Values(car.UserID, car.Brand, car.Model, car.LicensePlate, car.Color, car.Size, car.IsSelected).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBuildQuery, err)
	}

	var carID int64
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&carID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreateCar, err)
	}

	car.ID = carID
	return car, nil
}

// GetByID получает автомобиль по ID
func (r *Repository) GetByID(ctx context.Context, carID int64) (*domain.Car, error) {
	query, args, err := psqlbuilder.Select("id", "user_id", "brand", "model", "license_plate", "color", "size", "is_selected").
		From("cars").
		Where(squirrel.Eq{"id": carID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBuildQuery, err)
	}

	var car domain.Car
	err = r.db.GetContext(ctx, &car, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userservice.ErrCarNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrGetCar, err)
	}

	return &car, nil
}

// GetByUserID получает все автомобили пользователя
func (r *Repository) GetByUserID(ctx context.Context, userID int64) ([]*domain.Car, error) {
	query, args, err := psqlbuilder.Select("id", "user_id", "brand", "model", "license_plate", "color", "size", "is_selected").
		From("cars").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBuildQuery, err)
	}

	var cars []*domain.Car
	err = r.db.SelectContext(ctx, &cars, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetCar, err)
	}

	if cars == nil {
		cars = []*domain.Car{}
	}

	return cars, nil
}

// Update обновляет данные автомобиля
func (r *Repository) Update(ctx context.Context, car *domain.Car) error {
	query, args, err := psqlbuilder.Update("cars").
		Set("brand", car.Brand).
		Set("model", car.Model).
		Set("license_plate", car.LicensePlate).
		Set("color", car.Color).
		Set("size", car.Size).
		Set("is_selected", car.IsSelected).
		Where(squirrel.Eq{"id": car.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBuildQuery, err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateCar, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: failed to get rows affected: %v", ErrUpdateCar, err)
	}

	if rowsAffected == 0 {
		return userservice.ErrCarNotFound
	}

	return nil
}

// Delete удаляет автомобиль
func (r *Repository) Delete(ctx context.Context, carID int64) error {
	query, args, err := psqlbuilder.Delete("cars").
		Where(squirrel.Eq{"id": carID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBuildQuery, err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteCar, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: failed to get rows affected: %v", ErrDeleteCar, err)
	}

	if rowsAffected == 0 {
		return userservice.ErrCarNotFound
	}

	return nil
}

// GetSelectedByUserID получает выбранный автомобиль пользователя
func (r *Repository) GetSelectedByUserID(ctx context.Context, userID int64) (*domain.Car, error) {
	query, args, err := psqlbuilder.Select("id", "user_id", "brand", "model", "license_plate", "color", "size", "is_selected").
		From("cars").
		Where(squirrel.Eq{"user_id": userID, "is_selected": true}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBuildQuery, err)
	}

	var car domain.Car
	err = r.db.GetContext(ctx, &car, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userservice.ErrCarNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrGetCar, err)
	}

	return &car, nil
}

// UnselectAllByUserID снимает выбор со всех автомобилей пользователя
func (r *Repository) UnselectAllByUserID(ctx context.Context, userID int64) error {
	query, args, err := psqlbuilder.Update("cars").
		Set("is_selected", false).
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrBuildQuery, err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateCar, err)
	}

	return nil
}
