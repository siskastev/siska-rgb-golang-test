package handlers

import (
	"siska-rgb-golang-test/internal/products/services"

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
