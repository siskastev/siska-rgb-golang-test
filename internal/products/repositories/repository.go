package repositories

import (
	"context"
	"siska-rgb-golang-test/internal/models"

	"github.com/google/uuid"
)

type Repository interface {
	GetProductCategories(ctx context.Context) ([]models.ProductCategory, error)
	GetProductCategoryByID(ctx context.Context, id uint) (models.ProductCategory, error)
	CreateProductGift(ctx context.Context, product models.Product) (models.Product, error)
	GetProductGiftsByID(ctx context.Context, id uuid.UUID) (models.Product, error)
	UpdateProductGift(ctx context.Context, request models.Product) (models.Product, error)
	GetGiftsPagination(ctx context.Context, request models.GiftsFilter) ([]models.Product, error)
	DeleteGiftByID(ctx context.Context, id uuid.UUID) error
}
