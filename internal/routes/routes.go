package routes

import (
	"os"
	"siska-rgb-golang-test/internal/middleware"
	product "siska-rgb-golang-test/internal/products/routes"
	redeem "siska-rgb-golang-test/internal/transactions/redemptions/routes"
	users "siska-rgb-golang-test/internal/users/routes"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	jwtMiddleware := middleware.NewAuthMiddleware(os.Getenv("JWT_PRIVATE_KEY"))
	//add route here
	users.RegisterRoutes(app, jwtMiddleware)
	product.RegisterRoutes(app, jwtMiddleware)
	redeem.RegisterRoutes(app, jwtMiddleware)
}
