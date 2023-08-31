package routes

import (
	"siska-rgb-golang-test/internal/database"
	"siska-rgb-golang-test/internal/middleware"
	"siska-rgb-golang-test/internal/models"
	"siska-rgb-golang-test/internal/transactions/rating/handlers"
	"siska-rgb-golang-test/internal/transactions/rating/repositories"
	"siska-rgb-golang-test/internal/transactions/rating/services"
	redeemRepositories "siska-rgb-golang-test/internal/transactions/redemptions/repositories"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(route fiber.Router, jwtMiddleware *middleware.AuthMiddleware) {
	repo := repositories.NewRatingRepository(database.DB)
	redeemRepo := redeemRepositories.NewRedemptionRepository(database.DB)
	service := services.NewRatingService(repo, redeemRepo)
	handler := handlers.NewHandler(service)

	route.Use(jwtMiddleware.AuthRequired())
	route.Post("/gifts/:id/rating", jwtMiddleware.HasRoles(string(models.USER_ROLE)), handler.CreateRating)
}
