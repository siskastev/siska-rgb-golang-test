package models

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	ADMIN_ROLE Role = "admin"
	USER_ROLE  Role = "user"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;column:id;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"`
	Email     string    `gorm:"column:email;type:varchar(100);unique;not null"`
	Role      Role      `gorm:"column:role;not null;default:'team'"`
	Point     int       `gorm:"column:point;not null;default:0"`
	Password  string    `gorm:"column:password;type:varchar(100);not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      Role       `json:"role"`
	Point     int        `json:"point"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type UserRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6,max=10"`
}

type UserResponseWithToken struct {
	UserResponse
	Token string `json:"token"`
}
