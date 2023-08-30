package services

import (
	"context"
	"mime/multipart"
	"siska-rgb-golang-test/internal/models"

	"github.com/google/uuid"
)

type Service interface {
	GetProductCategories(ctx context.Context) ([]models.ProductCategoryResponse, error)
	IsCategoryIDExist(ctx context.Context, categoryID uint) (bool, error)
	CreateProductGift(ctx context.Context, product models.GiftRequest, imageFile *multipart.FileHeader) (models.GiftsResponse, error)
	IsProductGiftIDExist(ctx context.Context, id uuid.UUID) (bool, error)
	GetProductGiftByID(ctx context.Context, id uuid.UUID) (models.GiftsResponse, error)
	UpdateProductGift(ctx context.Context, request models.GiftRequest, file *multipart.FileHeader, id uuid.UUID) (models.GiftsResponse, error)
	UpdateProductGiftDescriptions(ctx context.Context, request models.GiftRequestDescriptions, id uuid.UUID) (models.GiftsResponse, error)
}
