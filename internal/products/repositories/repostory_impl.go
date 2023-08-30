package repositories

import (
	"context"
	"encoding/json"
	"siska-rgb-golang-test/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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

func (p *productRepository) GetProductCategoryByID(ctx context.Context, id uint) (category models.ProductCategory, err error) {
	if err := p.db.WithContext(ctx).Where(models.ProductCategory{ID: id}).First(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (p *productRepository) CreateProductGift(ctx context.Context, product models.Product) (models.Product, error) {
	if err := p.db.WithContext(ctx).Create(&product).Error; err != nil {
		return product, err
	}

	// Clear cache for all battles
	keys, err := p.redisClient.Keys(ctx, "gifts:*").Result()
	if err != nil {
		return product, err
	}

	if len(keys) > 0 {
		result := p.redisClient.Del(ctx, keys...)
		if result.Err() != nil {
			return product, result.Err()
		}
	}

	return product, nil
}

func (p *productRepository) GetProductGiftsByID(ctx context.Context, id uuid.UUID) (models.Product, error) {
	product := models.Product{}
	if err := p.db.WithContext(ctx).Where("id = ? AND point > 0", id).First(&product).Error; err != nil {
		return product, err
	}
	return product, nil
}
