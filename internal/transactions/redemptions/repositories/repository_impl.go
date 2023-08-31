package repositories

import (
	"context"
	"siska-rgb-golang-test/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type redemptionRepository struct {
	db *gorm.DB
}

func NewRedemptionRepository(db *gorm.DB) Repository {
	return &redemptionRepository{db: db}
}

func (r *redemptionRepository) GetUserRedemptionByIDProduct(ctx context.Context, productID uuid.UUID, userID uuid.UUID) (models.Redemption, error) {
	redemption := models.Redemption{}
	if err := r.db.WithContext(ctx).Where("product_id = ? AND user_id = ?", productID, userID).First(&redemption).Error; err != nil {
		return redemption, err
	}
	return redemption, nil
}

func (r *redemptionRepository) CreateRedemption(ctx context.Context, redemption models.Redemption) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&redemption).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Product{}).Where("id = ?", redemption.ProductID).Update("qty", gorm.Expr("qty - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
