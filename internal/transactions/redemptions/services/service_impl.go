package services

import (
	"context"
	"siska-rgb-golang-test/internal/models"
	productRepo "siska-rgb-golang-test/internal/products/repositories"
	"siska-rgb-golang-test/internal/transactions/redemptions/repositories"
	"time"

	"github.com/google/uuid"
)

type redemptionService struct {
	redemptionRepo repositories.Repository
	productRepo    productRepo.Repository
}

func NewRedemptionService(redemptionRepo repositories.Repository, productRepo productRepo.Repository) Service {
	return &redemptionService{redemptionRepo: redemptionRepo, productRepo: productRepo}
}

func (s *redemptionService) CreateRedemption(ctx context.Context, productID uuid.UUID, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	product, err := s.productRepo.GetProductGiftsByID(ctx, productID)
	if err != nil {
		return err
	}

	redemption := models.Redemption{
		ProductID:    productID,
		UserID:       userID,
		ProductName:  product.Name,
		CategoryName: product.ProductCategory.Name,
		Point:        product.Point,
		Descriptions: product.Descriptions,
		Image:        product.Image,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.redemptionRepo.CreateRedemption(ctx, redemption); err != nil {
		return err
	}

	return nil
}

func (s *redemptionService) IsExistUserRedemption(ctx context.Context, productID uuid.UUID, userID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	redemption, err := s.redemptionRepo.GetUserRedemptionByIDProduct(ctx, productID, userID)
	if err != nil {
		return false, err
	}

	return redemption.ID != uuid.Nil, nil
}

func (s *redemptionService) GetProductGiftByID(ctx context.Context, id uuid.UUID) (models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := s.productRepo.GetProductGiftsByID(ctx, id)
	if err != nil {
		return models.Product{}, err
	}

	return result, nil
}
