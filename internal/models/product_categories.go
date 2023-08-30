package models

import "time"

type ProductCategory struct {
	ID        uint      `gorm:"column:id;primary_key;auto_increment"`
	Name      string    `gorm:"column:name;varchar(100);not null;"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type ProductCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
