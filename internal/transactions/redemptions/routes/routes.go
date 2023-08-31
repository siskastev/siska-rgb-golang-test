package routes

import (
	"siska-rgb-golang-test/internal/database"
	"siska-rgb-golang-test/internal/middleware"
	"siska-rgb-golang-test/internal/models"
	productRepositories "siska-rgb-golang-test/internal/products/repositories"
	"siska-rgb-golang-test/internal/redis"
	"siska-rgb-golang-test/internal/transactions/redemptions/handlers"
	"siska-rgb-golang-test/internal/transactions/redemptions/repositories"
	"siska-rgb-golang-test/internal/transactions/redemptions/services"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(route fiber.Router, jwtMiddleware *middleware.AuthMiddleware) {
	repo := repositories.NewRedemptionRepository(database.DB)
	productRepo := productRepositories.NewProductRepository(database.DB, redis.RedisClient)
	service := services.NewRedemptionService(repo, productRepo)
	handler := handlers.NewHandler(service)

	route.Use(jwtMiddleware.AuthRequired())
	route.Post("/gifts/:id/redeem", jwtMiddleware.HasRoles(string(models.USER_ROLE)), handler.CreateRedemption)
}
