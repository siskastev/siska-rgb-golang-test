package pokemon

import (
	"siska-rgb-golang-test/internal/database"
	"siska-rgb-golang-test/internal/middleware"
	"siska-rgb-golang-test/internal/users/handlers"
	"siska-rgb-golang-test/internal/users/repositories"
	"siska-rgb-golang-test/internal/users/services"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(route fiber.Router, jwtMiddleware *middleware.AuthMiddleware) {
	repo := repositories.NewUserRepository(database.DB)
	service := services.NewUserService(repo)
	handler := handlers.NewHandlerUser(service)

	route.Post("/register", handler.Register)
	route.Post("/login", handler.Login)

	route.Use(jwtMiddleware.AuthRequired())
	route.Get("/profile", handler.Profile)
	route.Patch("/profile", handler.Update)
}
