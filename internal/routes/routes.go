package routes

import (
	"os"
	"siska-rgb-golang-test/internal/middleware"
	users "siska-rgb-golang-test/internal/users/routes"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")
	jwtMiddleware := middleware.NewAuthMiddleware(os.Getenv("JWT_PRIVATE_KEY"))
	//add route here
	users.RegisterRoutes(api, jwtMiddleware)
}
