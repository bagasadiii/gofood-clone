package model

import (
	"github.com/google/uuid"
)

type Driver struct {
	DriverID uuid.UUID
	Name     string
	Rating   float64
	License  string
	Area     string
	Income   int
	UserID   uuid.UUID
	Username string
}
type DriverRes struct {
	Name     string
	Rating   float64
	License  string
	Area     string
	Income   int
	Username string
}
