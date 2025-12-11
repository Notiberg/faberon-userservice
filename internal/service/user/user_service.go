package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/m04kA/SMC-UserService/internal/domain"
	"github.com/m04kA/SMC-UserService/internal/service/user/models"
)

var (
	ErrServiceCreateUser = errors.New("service: failed to create user")
	ErrServiceGetUser    = errors.New("service: failed to get user")
	ErrServiceUpdateUser = errors.New("service: failed to update user")
	ErrServiceDeleteUser = errors.New("service: failed to delete user")
	ErrServiceCreateCar  = errors.New("service: failed to create car")
	ErrServiceGetCar     = errors.New("service: failed to get car")
	ErrServiceUpdateCar  = errors.New("service: failed to update car")
	ErrServiceDeleteCar  = errors.New("service: failed to delete car")
)

type Service struct {
	userRepo UserRepository
	carRepo  CarRepository
}

func NewUserService(ur UserRepository, cr CarRepository) *Service {
	return &Service{userRepo: ur, carRepo: cr}
}

// CreateUser создает нового пользователя
func (s *Service) CreateUser(ctx context.Context, input models.CreateUserInputDTO) (*models.UserDTO, error) {
	_, err := s.userRepo.GetByTGID(ctx, input.TGUserID)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}
	if !errors.Is(err, ErrUserNotFound) {
		return nil, fmt.Errorf("%w: %v", ErrServiceGetUser, err)
	}

	// Маппинг роли в role_id
	roleID := roleToID(input.Role)

	user := &domain.User{
		TGUserID:    input.TGUserID,
		Name:        input.Name,
		PhoneNumber: input.PhoneNumber,
		TGLink:      input.TGLink,
		RoleID:      roleID,
		Role:        input.Role,
		CreatedAt:   time.Now(),
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrServiceCreateUser, err)
	}

	response := &models.UserDTO{
		TGUserID:    user.TGUserID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		TGLink:      user.TGLink,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
	}

	return response, nil
}

// UpdateUser обновляет данные пользователя (частичное обновление)
func (s *Service) UpdateUser(ctx context.Context, tgID int64, input models.UpdateUserInputDTO) (*models.UserDTO, error) {
	user, err := s.userRepo.GetByTGID(ctx, tgID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrServiceGetUser, err)
	}

	// Обновляем только те поля, которые переданы
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.PhoneNumber != nil {
		user.PhoneNumber = input.PhoneNumber
	}
	if input.TGLink != nil {
		user.TGLink = input.TGLink
	}

	if err = s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrServiceUpdateUser, err)
	}

	response := &models.UserDTO{
		TGUserID:    user.TGUserID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		TGLink:      user.TGLink,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
	}

	return response, nil
}

// DeleteUser удаляет пользователя
func (s *Service) DeleteUser(ctx context.Context, tgID int64) error {
	err := s.userRepo.Delete(ctx, tgID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return err
		}
		return fmt.Errorf("%w: %v", ErrServiceDeleteUser, err)
	}
	return nil
}

// GetUserByID получает пользователя по ID
func (s *Service) GetUserByID(ctx context.Context, tgID int64) (*models.UserDTO, error) {
	user, err := s.userRepo.GetByTGID(ctx, tgID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrServiceGetUser, err)
	}

	response := &models.UserDTO{
		TGUserID:    user.TGUserID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		TGLink:      user.TGLink,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
	}

	return response, nil
}

// GetUserWithCars получает пользователя со всеми его автомобилями
func (s *Service) GetUserWithCars(ctx context.Context, tgID int64) (*models.UserWithCarsDTO, error) {
	user, err := s.userRepo.GetByTGID(ctx, tgID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrServiceGetUser, err)
	}

	cars, err := s.carRepo.GetByUserID(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrServiceGetCar, err)
	}

	carDTOs := make([]models.CarDTO, 0, len(cars))
	for _, car := range cars {
		carDTOs = append(carDTOs, models.CarDTO{
			ID:           car.ID,
			UserID:       car.UserID,
			Brand:        car.Brand,
			Model:        car.Model,
			LicensePlate: car.LicensePlate,
			Color:        car.Color,
			Size:         car.Size,
			IsSelected:   car.IsSelected,
		})
	}

	response := &models.UserWithCarsDTO{
		TGUserID:    user.TGUserID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		TGLink:      user.TGLink,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		Cars:        carDTOs,
	}

	return response, nil
}

// GetSuperUsers возвращает список tg_user_id всех суперпользователей
func (s *Service) GetSuperUsers(ctx context.Context) ([]int64, error) {
	userIDs, err := s.userRepo.GetSuperUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrServiceGetUser, err)
	}
	return userIDs, nil
}

// roleToID маппит роль в ID для БД
func roleToID(role domain.Role) int {
	switch role {
	case domain.RoleClient:
		return domain.RoleIDClient
	case domain.RoleManager:
		return domain.RoleIDManager
	case domain.RoleSuperUser:
		return domain.RoleIDSuperUser
	default:
		return domain.RoleIDClient // default client
	}
}

// CreateCar создает новый автомобиль
func (s *Service) CreateCar(ctx context.Context, tgID int64, input models.CreateCarInputDTO) (*models.CarDTO, error) {
	_, err := s.userRepo.GetByTGID(ctx, tgID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrServiceGetUser, err)
	}

	// Проверяем, есть ли уже автомобили у пользователя
	existingCars, err := s.carRepo.GetByUserID(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrServiceGetCar, err)
	}

	// Если это первый автомобиль, он автоматически становится выбранным
	isSelected := len(existingCars) == 0

	car := &domain.Car{
		UserID:       tgID,
		Brand:        input.Brand,
		Model:        input.Model,
		LicensePlate: input.LicensePlate,
		Color:        input.Color,
		Size:         input.Size,
		IsSelected:   isSelected,
	}

	createdCar, err := s.carRepo.Create(ctx, car)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrServiceCreateCar, err)
	}

	response := &models.CarDTO{
		ID:           createdCar.ID,
		UserID:       createdCar.UserID,
		Brand:        createdCar.Brand,
		Model:        createdCar.Model,
		LicensePlate: createdCar.LicensePlate,
		Color:        createdCar.Color,
		Size:         createdCar.Size,
		IsSelected:   createdCar.IsSelected,
	}

	return response, nil
}

// UpdateCar обновляет автомобиль (PATCH) с проверкой роли
func (s *Service) UpdateCar(ctx context.Context, tgID int64, carID int64, input models.UpdateCarInputDTO, role domain.Role) (*models.CarDTO, error) {
	car, err := s.carRepo.GetByID(ctx, carID)
	if err != nil {
		if errors.Is(err, ErrCarNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrServiceGetCar, err)
	}

	// Проверка доступа: владелец может изменять свою машину, superuser - любую
	if !role.CanModifyUser(car.UserID, tgID) {
		return nil, ErrCarAccessDenied
	}

	if input.Brand != nil {
		car.Brand = *input.Brand
	}
	if input.Model != nil {
		car.Model = *input.Model
	}
	if input.LicensePlate != nil {
		car.LicensePlate = *input.LicensePlate
	}
	if input.Color != nil {
		car.Color = input.Color
	}
	if input.Size != nil {
		car.Size = input.Size
	}

	err = s.carRepo.Update(ctx, car)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrServiceUpdateCar, err)
	}

	response := &models.CarDTO{
		ID:           car.ID,
		UserID:       car.UserID,
		Brand:        car.Brand,
		Model:        car.Model,
		LicensePlate: car.LicensePlate,
		Color:        car.Color,
		Size:         car.Size,
		IsSelected:   car.IsSelected,
	}

	return response, nil
}

// DeleteCar удаляет автомобиль с проверкой роли
func (s *Service) DeleteCar(ctx context.Context, tgID int64, carID int64, role domain.Role) error {
	car, err := s.carRepo.GetByID(ctx, carID)
	if err != nil {
		if errors.Is(err, ErrCarNotFound) {
			return err
		}
		return fmt.Errorf("%w: %v", ErrServiceGetCar, err)
	}

	// Проверка доступа: владелец может удалять свою машину, superuser - любую
	if !role.CanModifyUser(car.UserID, tgID) {
		return ErrCarAccessDenied
	}

	wasSelected := car.IsSelected

	err = s.carRepo.Delete(ctx, carID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrServiceDeleteCar, err)
	}

	// Если удалили выбранный автомобиль, выбираем первый из оставшихся
	if wasSelected {
		remainingCars, err := s.carRepo.GetByUserID(ctx, car.UserID)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrServiceGetCar, err)
		}

		if len(remainingCars) > 0 {
			remainingCars[0].IsSelected = true
			if err := s.carRepo.Update(ctx, remainingCars[0]); err != nil {
				return fmt.Errorf("%w: %v", ErrServiceUpdateCar, err)
			}
		}
	}

	return nil
}

// GetSelectedCar получает текущий выбранный автомобиль пользователя
func (s *Service) GetSelectedCar(ctx context.Context, tgID int64) (*models.CarDTO, error) {
	car, err := s.carRepo.GetSelectedByUserID(ctx, tgID)
	if err != nil {
		if errors.Is(err, ErrCarNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrServiceGetCar, err)
	}

	response := &models.CarDTO{
		ID:           car.ID,
		UserID:       car.UserID,
		Brand:        car.Brand,
		Model:        car.Model,
		LicensePlate: car.LicensePlate,
		Color:        car.Color,
		Size:         car.Size,
		IsSelected:   car.IsSelected,
	}

	return response, nil
}

// SetSelectedCar устанавливает автомобиль как выбранный
func (s *Service) SetSelectedCar(ctx context.Context, tgID int64, carID int64, role domain.Role) (*models.CarDTO, error) {
	car, err := s.carRepo.GetByID(ctx, carID)
	if err != nil {
		if errors.Is(err, ErrCarNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrServiceGetCar, err)
	}

	// Проверка доступа
	if !role.CanModifyUser(car.UserID, tgID) {
		return nil, ErrCarAccessDenied
	}

	// Если уже выбран, просто возвращаем
	if car.IsSelected {
		response := &models.CarDTO{
			ID:           car.ID,
			UserID:       car.UserID,
			Brand:        car.Brand,
			Model:        car.Model,
			LicensePlate: car.LicensePlate,
			Color:        car.Color,
			Size:         car.Size,
			IsSelected:   car.IsSelected,
		}
		return response, nil
	}

	// Сначала снимаем выбор со всех автомобилей пользователя
	if err := s.carRepo.UnselectAllByUserID(ctx, car.UserID); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrServiceUpdateCar, err)
	}

	// Устанавливаем текущий автомобиль как выбранный
	car.IsSelected = true
	if err := s.carRepo.Update(ctx, car); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrServiceUpdateCar, err)
	}

	response := &models.CarDTO{
		ID:           car.ID,
		UserID:       car.UserID,
		Brand:        car.Brand,
		Model:        car.Model,
		LicensePlate: car.LicensePlate,
		Color:        car.Color,
		Size:         car.Size,
		IsSelected:   car.IsSelected,
	}

	return response, nil
}
