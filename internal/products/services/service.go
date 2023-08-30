package services

import (
	"context"
	"siska-rgb-golang-test/internal/models"
)

type Service interface {
	GetProductCategories(ctx context.Context) ([]models.ProductCategoryResponse, error)
}
