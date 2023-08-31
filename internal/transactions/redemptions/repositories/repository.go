package repositories

import (
	"context"
	"siska-rgb-golang-test/internal/models"

	"github.com/google/uuid"
)

type Repository interface {
	GetUserRedemptionByIDProduct(ctx context.Context, productID uuid.UUID, userID uuid.UUID) (models.Redemption, error)
	CreateRedemption(ctx context.Context, redemption models.Redemption) error
}
