package services

import (
	"context"
	"siska-rgb-golang-test/internal/models"

	"github.com/google/uuid"
)

type Service interface {
	CreateRating(ctx context.Context, rating models.RatingReviewRequest, productID uuid.UUID, userID uuid.UUID) error
	IsExistUserRedemption(ctx context.Context, productID uuid.UUID, userID uuid.UUID) (bool, error)
	IsAlreadyRatingUser(ctx context.Context, productID uuid.UUID, userID uuid.UUID) (bool, error)
}
