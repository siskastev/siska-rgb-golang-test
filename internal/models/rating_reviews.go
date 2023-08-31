package models

import (
	"time"

	"github.com/google/uuid"
)

type (
	RatingReview struct {
		ID        uint      `gorm:"column:id;primary_key;auto_increment"`
		Rating    float32   `gorm:"column:rating;not null;"`
		ProductID uuid.UUID `gorm:"column:product_id;not null"`
		UserID    uuid.UUID `gorm:"column:user_id;not null"`
		CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
		UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	}

	RatingReviewRequest struct {
		Rating float32 `form:"rating" validate:"required,min=1,max=5"`
	}
)
