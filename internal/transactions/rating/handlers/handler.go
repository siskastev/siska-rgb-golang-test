package handlers

import (
	"fmt"
	"siska-rgb-golang-test/internal/models"
	"siska-rgb-golang-test/internal/transactions/rating/services"

	uuid_helpers "siska-rgb-golang-test/internal/helpers/uuid"
	helpers "siska-rgb-golang-test/internal/helpers/validator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	ratingService services.Service
}

func NewHandler(ratingService services.Service) *Handler {
	return &Handler{
		ratingService: ratingService,
	}
}

func (h *Handler) CreateRating(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	var request models.RatingReviewRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request payload"})
	}

	userID := uuid_helpers.ParseUUID(c.Locals("user").(*models.UserResponse).ID)
	productID := uuid_helpers.ParseUUID(c.Params("id"))

	isAlreadyRedeem, _ := h.ratingService.IsExistUserRedemption(c.Context(), productID, userID)
	if !isAlreadyRedeem {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "Rating",
			"message": "Product ID not redeemed yet"},
		})
	}

	isAlreadyRating, _ := h.ratingService.IsAlreadyRatingUser(c.Context(), productID, userID)
	if isAlreadyRating {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "Rating",
			"message": "Product ID already rated"},
		})
	}

	if err := validator.New().Struct(request); err != nil {
		return helpers.ValidationErrorResponse(c, err.(validator.ValidationErrors))
	}

	ratingStr := fmt.Sprintf("%v", request.Rating)
	if !helpers.IsValidRating(ratingStr) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "Rating",
			"message": "Rating format is invalid"},
		})
	}

	if err := h.ratingService.CreateRating(c.Context(), request, productID, userID); err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to create rating")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": fmt.Sprintf("Rating for product ID %s successfully", productID)})
}
