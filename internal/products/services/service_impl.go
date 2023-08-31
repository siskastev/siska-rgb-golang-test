package services

import (
	"context"
	"math"
	"mime/multipart"
	helpers "siska-rgb-golang-test/internal/helpers/image"
	"siska-rgb-golang-test/internal/models"
	"siska-rgb-golang-test/internal/products/repositories"
	"time"

	"github.com/google/uuid"
)

type productService struct {
	productRepository repositories.Repository
}

func NewProductService(productRepository repositories.Repository) Service {
	return &productService{productRepository: productRepository}
}

func (p *productService) GetProductCategories(ctx context.Context) ([]models.ProductCategoryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
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

func (p *productService) IsCategoryIDExist(ctx context.Context, id uint) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := p.productRepository.GetProductCategoryByID(ctx, id)
	if err != nil {
		return false, err
	}

	return result.ID != 0, nil
}

func (p *productService) CreateProductGift(ctx context.Context, request models.GiftRequest, file *multipart.FileHeader) (models.GiftsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	imageURL, err := helpers.UploadImageToCloudinary(ctx, file)
	if err != nil {
		return models.GiftsResponse{}, err
	}

	gift := models.Product{
		Name:         request.Name,
		CategoryID:   request.CategoryID,
		Descriptions: request.Descriptions,
		Qty:          request.Qty,
		Price:        request.Price,
		Point:        request.Point,
		Image:        imageURL,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result, err := p.productRepository.CreateProductGift(ctx, gift)
	if err != nil {
		return models.GiftsResponse{}, err
	}

	response := models.GiftsResponse{
		ID:           result.ID.String(),
		Name:         result.Name,
		CategoryID:   result.CategoryID,
		Descriptions: result.Descriptions,
		Qty:          result.Qty,
		Price:        result.Price,
		Point:        result.Point,
		Rating:       result.Rating,
		Image:        imageURL,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	return response, nil
}

func (p *productService) IsProductGiftIDExist(ctx context.Context, id uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := p.productRepository.GetProductGiftsByID(ctx, id)
	if err != nil {
		return false, err
	}

	return result.ID != uuid.Nil, nil
}

func (p *productService) GetProductGiftByID(ctx context.Context, id uuid.UUID) (models.GiftsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := p.productRepository.GetProductGiftsByID(ctx, id)
	if err != nil {
		return models.GiftsResponse{}, err
	}

	response := models.GiftsResponse{
		ID:           result.ID.String(),
		Name:         result.Name,
		CategoryID:   result.CategoryID,
		CategoryName: &result.ProductCategory.Name,
		Descriptions: result.Descriptions,
		Qty:          result.Qty,
		Price:        result.Price,
		Point:        result.Point,
		Rating:       float32(math.Round(float64(result.Rating))),
		Image:        result.Image,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	return response, nil
}

func (p *productService) UpdateProductGift(ctx context.Context, request models.GiftRequest, file *multipart.FileHeader, id uuid.UUID) (models.GiftsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	imageURL, err := helpers.UploadImageToCloudinary(ctx, file)
	if err != nil {
		return models.GiftsResponse{}, err
	}

	gift := models.Product{
		ID:           id,
		Name:         request.Name,
		CategoryID:   request.CategoryID,
		Descriptions: request.Descriptions,
		Qty:          request.Qty,
		Price:        request.Price,
		Point:        request.Point,
		Image:        imageURL,
		UpdatedAt:    time.Now(),
	}

	result, err := p.productRepository.UpdateProductGift(ctx, gift)
	if err != nil {
		return models.GiftsResponse{}, err
	}

	response := models.GiftsResponse{
		ID:           result.ID.String(),
		Name:         result.Name,
		CategoryID:   result.CategoryID,
		Descriptions: result.Descriptions,
		Qty:          result.Qty,
		Price:        result.Price,
		Point:        result.Point,
		Rating:       result.Rating,
		Image:        imageURL,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	return response, nil
}

func (p *productService) UpdateProductGiftStock(ctx context.Context, request models.GiftRequestStock, id uuid.UUID) (models.GiftsResponse, error) {
	gift := models.Product{
		ID:        id,
		Qty:       request.Qty,
		UpdatedAt: time.Now(),
	}

	result, err := p.productRepository.UpdateProductGift(ctx, gift)
	if err != nil {
		return models.GiftsResponse{}, err
	}

	response := models.GiftsResponse{
		ID:           result.ID.String(),
		Name:         result.Name,
		CategoryID:   result.CategoryID,
		Descriptions: result.Descriptions,
		Qty:          result.Qty,
		Price:        result.Price,
		Point:        result.Point,
		Rating:       result.Rating,
		Image:        result.Image,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.UpdatedAt,
	}

	return response, nil

}

func (p *productService) GetGiftsPagination(ctx context.Context, filter models.GiftsFilter) (models.GiftsResponsePagination, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := p.productRepository.GetGiftsPagination(ctx, filter)
	if err != nil {
		return models.GiftsResponsePagination{}, err
	}

	totalData := len(result)
	totalPage := (int(totalData) + filter.PageSize - 1) / filter.PageSize
	nextPage := filter.Page + 1
	if nextPage > totalPage {
		nextPage = 0
	}
	metaData := models.MetaPagination{
		TotalData:   totalData,
		TotalPage:   totalPage,
		CurrentPage: filter.Page,
		NextPage:    nextPage,
		PageSize:    filter.PageSize,
	}

	gifts := []models.GiftsResponse{}
	for _, gift := range result {
		gifts = append(gifts, models.GiftsResponse{
			ID:           gift.ID.String(),
			Name:         gift.Name,
			CategoryID:   gift.CategoryID,
			CategoryName: &gift.ProductCategory.Name,
			Descriptions: gift.Descriptions,
			Qty:          gift.Qty,
			Price:        gift.Price,
			Point:        gift.Point,
			Rating:       float32(math.Round(float64(gift.Rating))),
			Image:        gift.Image,
			CreatedAt:    gift.CreatedAt,
			UpdatedAt:    gift.UpdatedAt,
		})
	}

	response := models.GiftsResponsePagination{
		Data: gifts,
		Meta: metaData,
	}

	return response, nil
}

func (p *productService) DeleteProductGift(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := p.productRepository.DeleteGiftByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
