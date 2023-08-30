package routes

import (
	"siska-rgb-golang-test/internal/database"
	"siska-rgb-golang-test/internal/middleware"
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

	route.Use(jwtMiddleware.AuthRequired())
	route.Get("/products/categories", handler.GetProductCategories)

	groupGift := route.Group("/gifts")
	groupGift.Post("", jwtMiddleware.IsAdmin(), handler.CreateGifts)
	groupGift.Put("/:id", jwtMiddleware.IsAdmin(), handler.UpdateProductGift)
	groupGift.Patch("/:id", jwtMiddleware.IsAdmin(), handler.UpdateProductGiftDescriptions)
	groupGift.Get("/:id", handler.GetGiftsByID)
}
