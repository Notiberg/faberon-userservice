package user

import (
	"context"
	"errors"

	"github.com/m04kA/SMC-UserService/internal/domain"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this telegram id already exists")
	ErrCarNotFound       = errors.New("car not found")
	ErrCarAccessDenied   = errors.New("access denied to this car")
)

// UserRepository определяет контракт для работы с хранилищем пользователей.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByTGID(ctx context.Context, tgID int64) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, tgID int64) error
	GetSuperUsers(ctx context.Context) ([]int64, error)
}

// CarRepository определяет контракт для работы с хранилищем автомобилей.
type CarRepository interface {
	Create(ctx context.Context, car *domain.Car) (*domain.Car, error)
	GetByID(ctx context.Context, carID int64) (*domain.Car, error)
	GetByUserID(ctx context.Context, userID int64) ([]*domain.Car, error)
	GetSelectedByUserID(ctx context.Context, userID int64) (*domain.Car, error)
	Update(ctx context.Context, car *domain.Car) error
	Delete(ctx context.Context, carID int64) error
	UnselectAllByUserID(ctx context.Context, userID int64) error
}
