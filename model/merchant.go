package model

import "github.com/google/uuid"

type Merchant struct {
	MerchantID  uuid.UUID `json:"merchant_id,omitempty"`
	Name        string    `json:"name" validate:"required"`
	Rating      float64   `json:"rating,omitempty"`
	Address     string    `json:"address" validate:"required"`
	Category    string    `json:"category" validate:"required"`
	Description string    `json:"description" validate:"required"`
	UserID      uuid.UUID `json:"user_id,omitempty"`
	Owner       string    `json:"owner,omitempty"`
}
type MerchantRes struct {
	Name        string  `json:"name"`
	Rating      float64 `json:"rating"`
	Address     string  `json:"address"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Owner       string  `json:"owner"`
}
