package models

import (
	"time"

	"github.com/m04kA/SMC-UserService/internal/domain"
)

// User DTOs

type CreateUserInputDTO struct {
	TGUserID    int64       `json:"tg_user_id" validate:"required"`
	Name        string      `json:"name" validate:"required"`
	PhoneNumber *string     `json:"phone_number" validate:"omitempty,e164"`
	TGLink      *string     `json:"tg_link"`
	Role        domain.Role `json:"role" validate:"required,oneof=client manager superuser"`
}

type UpdateUserInputDTO struct {
	Name        *string `json:"name" validate:"omitempty"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,e164"`
	TGLink      *string `json:"tg_link"`
}

type UserDTO struct {
	TGUserID    int64       `json:"tg_user_id"`
	Name        string      `json:"name"`
	PhoneNumber *string     `json:"phone_number,omitempty"`
	TGLink      *string     `json:"tg_link,omitempty"`
	Role        domain.Role `json:"role"`
	CreatedAt   time.Time   `json:"created_at"`
}

type UserWithCarsDTO struct {
	TGUserID    int64       `json:"tg_user_id"`
	Name        string      `json:"name"`
	PhoneNumber *string     `json:"phone_number,omitempty"`
	TGLink      *string     `json:"tg_link,omitempty"`
	Role        domain.Role `json:"role"`
	CreatedAt   time.Time   `json:"created_at"`
	Cars        []CarDTO    `json:"cars"`
}

// Car DTOs

type CreateCarInputDTO struct {
	Brand        string  `json:"brand" validate:"required"`
	Model        string  `json:"model" validate:"required"`
	LicensePlate string  `json:"license_plate" validate:"required"`
	Color        *string `json:"color"`
	Size         *string `json:"size"`
}

type UpdateCarInputDTO struct {
	Brand        *string `json:"brand"`
	Model        *string `json:"model"`
	LicensePlate *string `json:"license_plate"`
	Color        *string `json:"color"`
	Size         *string `json:"size"`
}

type CarDTO struct {
	ID           int64   `json:"id"`
	UserID       int64   `json:"user_id"`
	Brand        string  `json:"brand"`
	Model        string  `json:"model"`
	LicensePlate string  `json:"license_plate"`
	Color        *string `json:"color,omitempty"`
	Size         *string `json:"size,omitempty"`
	IsSelected   bool    `json:"is_selected"`
}
