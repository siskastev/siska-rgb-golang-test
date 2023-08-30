package services

import (
	"context"
	"siska-rgb-golang-test/internal/models"
	"siska-rgb-golang-test/internal/products/repositories"
)

type productService struct {
	productRepository repositories.Repository
}

func NewProductService(productRepository repositories.Repository) Service {
	return &productService{productRepository: productRepository}
}

func (p *productService) GetProductCategories(ctx context.Context) ([]models.ProductCategoryResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	productCategories, err := p.productRepository.GetProductCategories(ctx)
	if err != nil {
		return nil, err
	}

	productCategoryResponses := []models.ProductCategoryResponse{}
	for _, productCategory := range productCategories {
		productCategoryResponses = append(productCategoryResponses, models.ProductCategoryResponse{
			ID:   productCategory.ID,
			Name: productCategory.Name,
		})
	}

	return productCategoryResponses, nil
}
