package repositories

import (
	"context"
	"siska-rgb-golang-test/internal/models"
)

type Repository interface {
	GetProductCategories(ctx context.Context) ([]models.ProductCategory, error)
}
