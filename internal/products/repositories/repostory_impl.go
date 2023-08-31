package repositories

import (
	"context"
	"encoding/json"
	"fmt"
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
	if err := p.db.WithContext(ctx).Where("id = ? AND point > 0", id).Preload("ProductCategory").First(&product).Error; err != nil {
		return product, err
	}
	return product, nil
}

func (p *productRepository) UpdateProductGift(ctx context.Context, request models.Product) (models.Product, error) {
	product := models.Product{}
	if err := p.db.WithContext(ctx).Where("id = ?", request.ID).First(&product).Updates(request).Error; err != nil {
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

func (p *productRepository) GetGiftsPagination(ctx context.Context, request models.GiftsFilter) ([]models.Product, error) {
	cacheKey := fmt.Sprintf("gifts:page=%d:page_size:%d:filter=%s", request.Page, request.PageSize, map[string]interface{}{"is_stock": request.IsStock, "rating": request.Rating, "sort_by": request.SortBy})
	products := []models.Product{}

	cachedData, err := p.redisClient.Get(ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if cachedData != "" {
		var products []models.Product
		if err := json.Unmarshal([]byte(cachedData), &products); err != nil {
			return nil, err
		}
		return products, nil
	}
	offset := (request.Page - 1) * request.PageSize
	limit := request.PageSize

	query := p.db.WithContext(ctx).Where("point > 0").Preload("ProductCategory")

	if request.IsStock {
		query = query.Where("qty > 0")
	}

	if request.Rating != 0 {
		query = query.Where("rating > ?", request.Rating)
	}

	switch request.SortBy {
	case "created_at":
		query = query.Order("created_at ASC")
	default:
		query = query.Order("created_at DESC")
	}

	if condition := query.Order("rating DESC").Offset(offset).Limit(limit).Find(&products); condition.Error != nil {
		return products, condition.Error
	}

	// Set cache
	jsonData, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}

	if err := p.redisClient.Set(ctx, cacheKey, jsonData, 1*time.Hour).Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *productRepository) DeleteGiftByID(ctx context.Context, id uuid.UUID) error {
	product := models.Product{}
	if err := p.db.WithContext(ctx).Where("id = ?", id).Delete(&product).Error; err != nil {
		return err
	}

	// Clear cache for all battles
	keys, err := p.redisClient.Keys(ctx, "gifts:*").Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		result := p.redisClient.Del(ctx, keys...)
		if result.Err() != nil {
			return result.Err()
		}
	}

	return nil
}
