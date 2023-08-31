package handlers

import (
	"siska-rgb-golang-test/internal/models"
	"siska-rgb-golang-test/internal/products/services"
	"strconv"

	"siska-rgb-golang-test/internal/helpers/uuid"
	helpers "siska-rgb-golang-test/internal/helpers/validator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type HandlerProduct struct {
	ProductService services.Service
}

func NewHandlerProduct(productService services.Service) *HandlerProduct {
	return &HandlerProduct{ProductService: productService}
}

func (h *HandlerProduct) GetProductCategories(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	productCategories, err := h.ProductService.GetProductCategories(c.Context())
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to get product categories")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": productCategories})
}

func (h *HandlerProduct) CreateGifts(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	var request models.GiftRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	if err := validator.New().Struct(request); err != nil {
		return helpers.ValidationErrorResponse(c, err.(validator.ValidationErrors))
	}

	isCategoryIDExists, _ := h.ProductService.IsCategoryIDExist(c.Context(), request.CategoryID)
	if !isCategoryIDExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "CategoryID",
			"message": "Category ID not exists"},
		})
	}

	imageFile, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "Image",
			"message": "Failed to get image"},
		})
	}

	if !helpers.IsValidImageType(imageFile) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "Image",
			"message": "Invalid image type. Only JPG, JPEG and PNG are allowed."},
		})
	}

	if !helpers.IsValidSizeImage(imageFile) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "Image",
			"message": "Image size exceeds the limit of 5MB"},
		})
	}

	giftResponse, err := h.ProductService.CreateProductGift(c.Context(), request, imageFile)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method":  c.Method(),
			"route":   c.Path(),
			"error":   err.Error(),
			"payload": request,
		}).Error("Failed to parse create gifts")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": giftResponse})
}

func (h *HandlerProduct) GetGiftsByID(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	id := uuid.ParseUUID(c.Params("id"))

	isGiftIDExists, err := h.ProductService.IsProductGiftIDExist(c.Context(), id)
	if !isGiftIDExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	giftResponse, err := h.ProductService.GetProductGiftByID(c.Context(), id)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to get gift by id")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": giftResponse})
}

func (h *HandlerProduct) UpdateProductGift(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	var request models.GiftRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	id := uuid.ParseUUID(c.Params("id"))

	isGiftIDExists, err := h.ProductService.IsProductGiftIDExist(c.Context(), id)
	if !isGiftIDExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	if err := validator.New().Struct(request); err != nil {
		return helpers.ValidationErrorResponse(c, err.(validator.ValidationErrors))
	}

	isCategoryIDExists, _ := h.ProductService.IsCategoryIDExist(c.Context(), request.CategoryID)
	if !isCategoryIDExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "CategoryID",
			"message": "Category ID not exists"},
		})
	}

	imageFile, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "Image",
			"message": "Failed to get image"},
		})
	}

	if !helpers.IsValidImageType(imageFile) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "Image",
			"message": "Invalid image type. Only JPG, JPEG and PNG are allowed."},
		})
	}

	if !helpers.IsValidSizeImage(imageFile) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "Image",
			"message": "Image size exceeds the limit of 5MB"},
		})
	}

	giftResponse, err := h.ProductService.UpdateProductGift(c.Context(), request, imageFile, id)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method":  c.Method(),
			"route":   c.Path(),
			"error":   err.Error(),
			"payload": request,
		}).Error("Failed to parse update gifts")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": giftResponse})
}

func (h *HandlerProduct) UpdateProductGiftStock(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	var request models.GiftRequestStock

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	id := uuid.ParseUUID(c.Params("id"))

	isGiftIDExists, err := h.ProductService.IsProductGiftIDExist(c.Context(), id)
	if !isGiftIDExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	if err := validator.New().Struct(request); err != nil {
		return helpers.ValidationErrorResponse(c, err.(validator.ValidationErrors))
	}

	giftResponse, err := h.ProductService.UpdateProductGiftStock(c.Context(), request, id)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method":  c.Method(),
			"route":   c.Path(),
			"error":   err.Error(),
			"payload": request,
		}).Error("Failed to parse update gifts")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": giftResponse})
}

func (h *HandlerProduct) GetGifts(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)
	var filter models.GiftsFilter
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}
	filter.Page = page
	filter.PageSize = pageSize
	filter.SortBy = c.Query("sort_by", "-created_at")
	filter.IsStock, _ = strconv.ParseBool(c.Query("is_stock", "false"))
	rating, _ := strconv.ParseFloat(c.Query("rating", "0"), 64)
	filter.Rating = float32(rating)

	if err := c.QueryParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request query"})
	}

	gifts, err := h.ProductService.GetGiftsPagination(c.Context(), filter)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to get gifts")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(gifts)
}

func (h *HandlerProduct) DeleteGiftsByID(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	id := uuid.ParseUUID(c.Params("id"))

	isGiftIDExists, err := h.ProductService.IsProductGiftIDExist(c.Context(), id)
	if !isGiftIDExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
	}

	err = h.ProductService.DeleteProductGift(c.Context(), id)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to delete gift")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Gift deleted successfully"})
}
