package repositories

import (
	"context"
	"encoding/json"
	"siska-rgb-golang-test/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type productRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewProductRepository(db *gorm.DB, redisClient *redis.Client) Repository {
	return &productRepository{db: db, redisClient: redisClient}
}

func (p *productRepository) GetProductCategories(ctx context.Context) ([]models.ProductCategory, error) {
	productCategories := []models.ProductCategory{}
	cacheKey := "product_categories"
	cachedData, err := p.redisClient.Get(ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		return productCategories, err
	}

	if cachedData != "" {
		if err := json.Unmarshal([]byte(cachedData), &productCategories); err != nil {
			return nil, err
		}
		return productCategories, nil
	}

	if err := p.db.WithContext(ctx).Find(&productCategories).Error; err != nil {
		return nil, err
	}

	// Set cache
	jsonData, err := json.Marshal(productCategories)
	if err != nil {
		return nil, err
	}
	if err := p.redisClient.Set(ctx, cacheKey, jsonData, 1*time.Hour).Err(); err != nil {
		return nil, err
	}

	return productCategories, nil

}
