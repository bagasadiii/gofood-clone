package model

import "github.com/google/uuid"

type Menu struct {
	MenuID      uuid.UUID `json:"menu_id"`
	Name        string    `json:"name" validate:"required"`
	Price       int64     `json:"price" validate:"required"`
	Description string    `json:"description"`
	Category    string    `json:"category" validate:"required"`
	Rating      float64   `json:"rating,omitempty"`
	Stock       int       `json:"stock,omitempty"`
	MerchantID  uuid.UUID `json:"merchant_id,omitempty"`
}

type MenuRes struct {
	Name        string  `json:"name"`
	Price       int64   `json:"price"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Rating      float64 `json:"rating"`
	Stock       int     `json:"stock"`
}
