package repositories

import (
	"context"
	"siska-rgb-golang-test/internal/models"

	"github.com/google/uuid"
)

type Repository interface {
	CreateRating(ctx context.Context, rating models.RatingReview) error
	GetRatingUserByProductID(ctx context.Context, productID uuid.UUID, userID uuid.UUID) (models.RatingReview, error)
}
