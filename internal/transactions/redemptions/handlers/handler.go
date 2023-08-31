package handlers

import (
	"fmt"
	uuid_helpers "siska-rgb-golang-test/internal/helpers/uuid"
	"siska-rgb-golang-test/internal/models"
	"siska-rgb-golang-test/internal/transactions/redemptions/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	redeemService services.Service
}

func NewHandler(redeemService services.Service) *Handler {
	return &Handler{
		redeemService: redeemService,
	}
}

func (h *Handler) CreateRedemption(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*logrus.Logger)

	userID := uuid_helpers.ParseUUID(c.Locals("user").(*models.UserResponse).ID)
	productID := uuid_helpers.ParseUUID(c.Params("id"))

	product, _ := h.redeemService.GetProductGiftByID(c.Context(), productID)
	if product.ID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "ProductID",
			"message": "Product ID not exists"},
		})
	}

	if product.Qty == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "ProductID",
			"message": "Product ID stock is not available"},
		})
	}

	isExistUserRedemption, _ := h.redeemService.IsExistUserRedemption(c.Context(), productID, userID)
	if isExistUserRedemption {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": fiber.Map{
			"field":   "ProductID",
			"message": "Product ID already redeemed"},
		})
	}

	if err := h.redeemService.CreateRedemption(c.Context(), productID, userID); err != nil {
		logger.WithFields(logrus.Fields{
			"method": c.Method(),
			"route":  c.Path(),
			"error":  err.Error(),
		}).Error("Failed to create redemption")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": fmt.Sprintf("Redemption for product ID %s successfully", productID)})
}
