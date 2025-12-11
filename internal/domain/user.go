package domain

import "time"

type User struct {
	TGUserID    int64     `json:"tg_user_id" db:"tg_user_id"`
	Name        string    `json:"name" db:"name" validate:"required"`
	PhoneNumber *string   `json:"phone_number" db:"phone_number" validate:"omitempty,e164"`
	TGLink      *string   `json:"tg_link" db:"tg_link"`
	RoleID      int       `json:"role_id" db:"role_id"`
	Role        Role      `json:"role" db:"role_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
