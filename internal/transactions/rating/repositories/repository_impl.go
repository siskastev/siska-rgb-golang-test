package repositories

import (
	"context"
	"siska-rgb-golang-test/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ratingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) Repository {
	return &ratingRepository{db: db}
}

func (r *ratingRepository) CreateRating(ctx context.Context, rating models.RatingReview) error {
	if err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&rating).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Product{}).Where("id = ?", rating.ProductID).Update("rating", gorm.Expr("rating + ?", rating.Rating)).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *ratingRepository) GetRatingUserByProductID(ctx context.Context, productID uuid.UUID, userID uuid.UUID) (models.RatingReview, error) {
	rating := models.RatingReview{}
	if err := r.db.WithContext(ctx).Where("product_id = ? AND user_id = ?", productID, userID).First(&rating).Error; err != nil {
		return rating, err
	}
	return rating, nil
}
