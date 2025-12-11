package domain

type Car struct {
	ID           int64   `json:"id" db:"id"`
	UserID       int64   `json:"user_id" db:"user_id"`
	Brand        string  `json:"brand" db:"brand" validate:"required"`
	Model        string  `json:"model" db:"model" validate:"required"`
	LicensePlate string  `json:"license_plate" db:"license_plate" validate:"required"`
	Color        *string `json:"color,omitempty" db:"color"`
	Size         *string `json:"size,omitempty" db:"size"`
	IsSelected   bool    `json:"is_selected" db:"is_selected"`
}
