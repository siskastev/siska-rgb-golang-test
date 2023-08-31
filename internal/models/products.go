package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Product struct {
		ID              uuid.UUID       `gorm:"type:uuid;column:id;primaryKey;default:gen_random_uuid()"`
		CategoryID      uint            `gorm:"column:category_id;not null"`
		Name            string          `gorm:"column:name;type:varchar(100);not null"`
		Descriptions    string          `gorm:"column:descriptions;not null"`
		Qty             uint16          `gorm:"column:qty;not null;default:0"`
		Price           float64         `gorm:"column:price;not null"`
		Point           int             `gorm:"column:point;not null;default:0"`
		Rating          float32         `gorm:"column:rating;not null"`
		Image           string          `gorm:"column:image;not null"`
		CreatedAt       time.Time       `gorm:"column:created_at;autoCreateTime"`
		UpdatedAt       time.Time       `gorm:"column:updated_at;autoUpdateTime"`
		DeletedAt       gorm.DeletedAt  `gorm:"column:deleted_at;"` // soft delete
		ProductCategory ProductCategory `gorm:"foreignKey:CategoryID"`
	}

	GiftRequest struct {
		Name         string  `form:"name" validate:"required,min=3,max=100"`
		CategoryID   uint    `form:"category_id" validate:"required"`
		Descriptions string  `form:"descriptions" validate:"required"`
		Qty          uint16  `form:"qty" validate:"required,min=1,max=1000"`
		Price        float64 `form:"price" validate:"required,min=1"`
		Point        int     `form:"point" validate:"required,min=1"`
		Image        string  `form:"image"`
	}

	GiftsResponse struct {
		ID           string    `json:"id"`
		Name         string    `json:"name"`
		CategoryID   uint      `json:"category_id"`
		CategoryName *string   `json:"category_name,omitempty"`
		Descriptions string    `json:"descriptions"`
		Qty          uint16    `json:"qty"`
		Price        float64   `json:"price"`
		Point        int       `json:"point"`
		Rating       float32   `json:"rating"`
		Image        string    `json:"image"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	GiftRequestStock struct {
		Qty uint16 `form:"qty" validate:"required,min=1,max=1000"`
	}

	GiftsResponsePagination struct {
		Data []GiftsResponse `json:"data"`
		Meta MetaPagination  `json:"meta"`
	}

	MetaPagination struct {
		TotalData   int `json:"total_data"`
		TotalPage   int `json:"total_page"`
		CurrentPage int `json:"current_page"`
		NextPage    int `json:"next_page"`
		PageSize    int `json:"page_size"`
	}

	GiftsFilter struct {
		Page     int     `form:"page"`
		PageSize int     `form:"page_size"`
		IsStock  bool    `form:"is_stock"`
		Rating   float32 `form:"rating"`
		SortBy   string  `form:"sort_by"`
	}
)
