package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	IsOnline  bool      `json:"is_online"`
	Phone     string    `json:"phone,omitempty"`
	Balance   int64     `json:"balance,omitempty"`
	Name      string    `json:"name"`
}
type UserResp struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	Phone     string    `json:"phone,omitempty"`
	Name      string    `json:"name"`
}
type RegisterReq struct {
	Username string `json:"username" validate:"required,min=3,max=24"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=driver merchant user"`
	Phone    string `json:"phone" validate:"required,e164"`
}
type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
