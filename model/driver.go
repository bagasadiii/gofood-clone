package model

import (
	"github.com/google/uuid"
)

type Driver struct {
	DriverID uuid.UUID `json:"driver_id,omitempty"`
	Name     string    `json:"name" validate:"required"`
	Rating   float64   `json:"rating,omitempty"`
	License  string    `json:"license" validate:"required"`
	Area     string    `json:"area" validate:"required"`
	Income   int       `json:"income,omitempty"`
	UserID   uuid.UUID `json:"user_id,omitempty"`
	Username string    `json:"username,omitempty"`
}

type DriverRes struct {
	Name     string  `json:"name"`
	Rating   float64 `json:"rating"`
	License  string  `json:"license"`
	Area     string  `json:"area"`
	Income   int     `json:"income"`
	Username string  `json:"username"`
}
