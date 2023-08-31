package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	Redemption struct {
		ID           uuid.UUID `gorm:"type:uuid;column:id;primaryKey;default:gen_random_uuid()"`
		ProductID    uuid.UUID `gorm:"column:product_id;not null"`
		UserID       uuid.UUID `gorm:"column:user_id;not null"`
		ProductName  string    `gorm:"column:product_name;type:varchar(100);not null"`
		CategoryName string    `gorm:"column:category_name;type:varchar(100);not null"`
		Point        int       `gorm:"column:point;not null;default:0"`
		Descriptions string    `gorm:"column:descriptions;not null"`
		Image        string    `gorm:"column:image;not null"`
		CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
		UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
	}

	RedemptionResponse struct {
		ID           string    `json:"id"`
		ProductID    string    `json:"product_id"`
		UserID       string    `json:"user_id"`
		ProductName  string    `json:"product_name"`
		CategoryName string    `json:"category_name"`
		Point        int       `json:"point"`
		Descriptions string    `json:"descriptions"`
		Image        string    `json:"image"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}
)
