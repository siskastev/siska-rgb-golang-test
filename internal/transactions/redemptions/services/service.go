package services

import (
	"context"
	"siska-rgb-golang-test/internal/models"

	"github.com/google/uuid"
)

type Service interface {
	CreateRedemption(ctx context.Context, productID uuid.UUID, userID uuid.UUID) error
	IsExistUserRedemption(ctx context.Context, productID uuid.UUID, userID uuid.UUID) (bool, error)
	GetProductGiftByID(ctx context.Context, id uuid.UUID) (models.Product, error)
}
