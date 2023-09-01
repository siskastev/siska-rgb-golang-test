package routes

import (
	"siska-rgb-golang-test/internal/database"
	"siska-rgb-golang-test/internal/middleware"
	"siska-rgb-golang-test/internal/models"
	"siska-rgb-golang-test/internal/products/handlers"
	"siska-rgb-golang-test/internal/products/repositories"
	"siska-rgb-golang-test/internal/products/services"
	"siska-rgb-golang-test/internal/redis"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(route fiber.Router, jwtMiddleware *middleware.AuthMiddleware) {
	repo := repositories.NewProductRepository(database.DB, redis.RedisClient)
	service := services.NewProductService(repo)
	handler := handlers.NewHandlerProduct(service)

	route.Get("/products/categories", handler.GetProductCategories)

	groupGift := route.Group("/gifts")
	groupGift.Get("/:id", handler.GetGiftsByID)
	groupGift.Get("", handler.GetGifts)

	groupGift.Use(jwtMiddleware.AuthRequired())
	groupGift.Post("", jwtMiddleware.HasRoles(string(models.ADMIN_ROLE)), handler.CreateGifts)
	groupGift.Put("/:id", jwtMiddleware.HasRoles(string(models.ADMIN_ROLE)), handler.UpdateProductGift)
	groupGift.Patch("/:id", jwtMiddleware.HasRoles(string(models.ADMIN_ROLE)), handler.UpdateProductGiftStock)
	groupGift.Delete("/:id", jwtMiddleware.HasRoles(string(models.ADMIN_ROLE)), handler.DeleteGiftsByID)
}
