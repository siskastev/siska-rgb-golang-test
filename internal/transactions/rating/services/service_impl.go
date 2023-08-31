package services

import (
	"context"
	"siska-rgb-golang-test/internal/models"
	"siska-rgb-golang-test/internal/transactions/rating/repositories"
	redeemRepository "siska-rgb-golang-test/internal/transactions/redemptions/repositories"
	"time"

	"github.com/google/uuid"
)

type ratingService struct {
	ratingRepo     repositories.Repository
	redemptionRepo redeemRepository.Repository
}

func NewRatingService(ratingRepo repositories.Repository, redemptionRepo redeemRepository.Repository) Service {
	return &ratingService{
		ratingRepo:     ratingRepo,
		redemptionRepo: redemptionRepo,
	}
}

func (r *ratingService) CreateRating(ctx context.Context, rating models.RatingReviewRequest, productID uuid.UUID, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ratingReview := models.RatingReview{
		ProductID: productID,
		UserID:    userID,
		Rating:    rating.Rating,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.ratingRepo.CreateRating(ctx, ratingReview); err != nil {
		return err
	}

	return nil
}

func (r *ratingService) IsExistUserRedemption(ctx context.Context, productID uuid.UUID, userID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	redemption, err := r.redemptionRepo.GetUserRedemptionByIDProduct(ctx, productID, userID)
	if err != nil {
		return false, err
	}

	return redemption.ID != uuid.Nil, nil
}

func (r *ratingService) IsAlreadyRatingUser(ctx context.Context, productID uuid.UUID, userID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rating, err := r.ratingRepo.GetRatingUserByProductID(ctx, productID, userID)
	if err != nil {
		return false, err
	}

	return rating.ID != 0, nil
}
